import { json } from "@sveltejs/kit";
import { Buffer } from "node:buffer";
import type Stripe from "stripe";

import {
	claimBillingInvoiceEmailDelivery,
	fetchBillingInvoicePDF,
	isBillingApiError,
	markBillingInvoiceEmailSent,
	markBillingOrderFailed,
	markBillingOrderPaid,
	type BillingInvoiceResponse,
	type BillingOrderResponse
} from "$lib/server/billing";
import { rootLogger, safeError, type Logger } from "$lib/server/logging";
import { sendInvoicePaidEmail } from "$lib/server/mail";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { isStripeCheckoutError, stripe, stripeWebhookSecret } from "$lib/server/stripe";
import type { RequestHandler } from "./$types";

/**
 * Handles error reporting for webhook processing failures.
 * Logs the error and generates standard JSON error response.
 * 
 * @param error - Webhook processing error object
 */
function webhookErrorResponse(error: unknown, logger: Logger) {
	logger.error("billing.stripe_webhook_error", { error: safeError(error) });
	if (isBillingApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	if (isStripeCheckoutError(error)) {
		return jsonPublicErrorResponse(500, error.message, "Stripe webhook failed");
	}

	throw error;
}

/**
 * Validates and extracts Stripe Checkout Session object from Event payload.
 * 
 * @param event - The validated Stripe Event
 */
function checkoutSession(event: Stripe.Event) {
	const object = event.data.object;
	if ("object" in object && object.object === "checkout.session") {
		return object as Stripe.Checkout.Session;
	}
	return null;
}

/**
 * Retrieves our internal billing order ID from Checkout Session metadata.
 */
function billingOrderID(session: Stripe.Checkout.Session) {
	return session.metadata?.billing_order_id?.trim() || "";
}

/**
 * Resolves the Stripe Payment Intent ID associated with checkout session.
 */
function paymentIntentID(session: Stripe.Checkout.Session) {
	const paymentIntent = session.payment_intent;
	if (typeof paymentIntent === "string") return paymentIntent;
	if (paymentIntent && typeof paymentIntent === "object" && "id" in paymentIntent) {
		return paymentIntent.id;
	}
	return undefined;
}

/**
 * Signals backend to mark the corresponding billing order as paid.
 */
async function markSessionPaid(
	fetchFn: typeof fetch,
	event: Stripe.Event,
	session: Stripe.Checkout.Session
): Promise<BillingOrderResponse | Response> {
	const orderId = billingOrderID(session);
	if (!orderId) {
		return json({ error: "missing billing order metadata" }, { status: 400 });
	}

	return await markBillingOrderPaid(fetchFn, {
		orderId,
		checkoutSessionId: session.id,
		paymentIntentId: paymentIntentID(session),
		paidAt: new Date(event.created * 1000).toISOString()
	});
}

/**
 * Signals backend to mark the billing order as failed.
 */
async function markSessionFailed(
	fetchFn: typeof fetch,
	session: Stripe.Checkout.Session,
	logger: Logger
) {
	const orderId = billingOrderID(session);
	if (!orderId) {
		return json({ error: "missing billing order metadata" }, { status: 400 });
	}

	try {
		await markBillingOrderFailed(fetchFn, orderId);
	} catch (error) {
		if (isBillingApiError(error) && error.status === 404) {
			logger.warn("billing.stripe_webhook_stale_failed_session_ignored", {
				checkout_session_id: session.id,
				order_id: orderId
			});
			return null;
		}
		throw error;
	}
	return null;
}

/**
 * Delivers a generated billing invoice PDF after payment processing succeeds.
 * Delivery failures are logged and do not affect Stripe webhook acknowledgement.
 */
async function deliverInvoiceEmailForPaidOrder(
	fetchFn: typeof fetch,
	order: BillingOrderResponse,
	logger: Logger
) {
	const invoice = order.invoice;
	if (!invoice?.id) return;

	try {
		const claim = await claimBillingInvoiceEmailDelivery(fetchFn, {
			invoiceId: invoice.id,
			userId: order.user_id
		});
		if (claim.status !== "claimed" || !claim.invoice) return;

		const pdf = await fetchBillingInvoicePDF(fetchFn, order.user_id, invoice.id, {
			download: true
		});
		const pdfContent = Buffer.from(await new Response(pdf.body).arrayBuffer());
		if (pdfContent.byteLength === 0) {
			throw new Error("invoice PDF is empty");
		}

		await sendInvoicePaidEmail({
			to: claim.invoice.billing_email,
			billingName: claim.invoice.billing_name,
			invoiceLabel: invoiceLabel(claim.invoice),
			invoiceDate: claim.invoice.invoice_date,
			totalAmount: claim.invoice.total_amount,
			currency: order.currency,
			credits: order.credits,
			orderId: order.id,
			pdf: {
				filename: invoicePDFFilename(pdf.headers, claim.invoice),
				content: pdfContent,
				contentType: pdf.headers.get("content-type") || "application/pdf"
			}
		});

		await markBillingInvoiceEmailSent(fetchFn, {
			invoiceId: invoice.id,
			userId: order.user_id
		});
	} catch (error) {
		logger.error("billing.invoice_email_delivery_failed", {
			order_id: order.id,
			invoice_id: invoice.id,
			error: safeError(error)
		});
	}
}

function invoiceLabel(invoice: BillingInvoiceResponse) {
	return `${invoice.invoice_serie}-${String(invoice.invoice_number).padStart(5, "0")}`;
}

function invoicePDFFilename(headers: Headers, invoice: BillingInvoiceResponse) {
	return (
		safeFilename(contentDispositionFilename(headers.get("content-disposition"))) ||
		invoicePDFFilenameFallback(invoice)
	);
}

function contentDispositionFilename(value: string | null) {
	if (!value) return "";

	const utf8Match = value.match(/filename\*=UTF-8''([^;]+)/i);
	if (utf8Match?.[1]) {
		try {
			return decodeURIComponent(utf8Match[1].trim());
		} catch {
			return utf8Match[1].trim();
		}
	}

	const quotedMatch = value.match(/filename="([^"]+)"/i);
	if (quotedMatch?.[1]) return quotedMatch[1].trim();

	const plainMatch = value.match(/filename=([^;]+)/i);
	return plainMatch?.[1]?.trim() ?? "";
}

function safeFilename(value: string) {
	const filename = value.trim();
	if (!filename || filename.includes("/") || filename.includes("\\")) return "";
	return filename;
}

function invoicePDFFilenameFallback(invoice: BillingInvoiceResponse) {
	const compactDate = invoice.invoice_date.replaceAll("-", "").slice(2) || "invoice";
	return `${invoice.invoice_serie}_${String(invoice.invoice_number).padStart(5, "0")}_${compactDate}.pdf`;
}

/**
 * POST /api/billing/stripe/webhook
 * Handles incoming webhooks from Stripe (session completions, payments, failures).
 * Validates Stripe signature for authenticity.
 */
export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	const logger = (locals?.logger ?? rootLogger).child({
		component: "stripe_webhook",
		domain: "billing"
	});
	const signature = request.headers.get("stripe-signature");
	if (!signature) {
		return json({ error: "missing stripe signature" }, { status: 400 });
	}

	let secret: string;
	try {
		secret = stripeWebhookSecret();
	} catch (error) {
		return webhookErrorResponse(error, logger);
	}

	const body = await request.text();

	let event: Stripe.Event;
	try {
		event = stripe().webhooks.constructEvent(body, signature, secret);
	} catch (error) {
		logger.error("billing.stripe_webhook_signature_failed", { error: safeError(error) });
		return json({ error: "invalid stripe signature" }, { status: 400 });
	}

	try {
		if (event.type === "checkout.session.completed") {
			const session = checkoutSession(event);
			if (!session) return json({ error: "invalid checkout session" }, { status: 400 });

			if (session.payment_status === "paid") {
				const paidOrder = await markSessionPaid(fetch, event, session);
				if (paidOrder instanceof Response) return paidOrder;
				await deliverInvoiceEmailForPaidOrder(fetch, paidOrder, logger);
			}
		}

		if (event.type === "checkout.session.async_payment_succeeded") {
			const session = checkoutSession(event);
			if (!session) return json({ error: "invalid checkout session" }, { status: 400 });

			const paidOrder = await markSessionPaid(fetch, event, session);
			if (paidOrder instanceof Response) return paidOrder;
			await deliverInvoiceEmailForPaidOrder(fetch, paidOrder, logger);
		}

		if (
			event.type === "checkout.session.expired" ||
			event.type === "checkout.session.async_payment_failed"
		) {
			const session = checkoutSession(event);
			if (!session) return json({ error: "invalid checkout session" }, { status: 400 });

			const errorResponse = await markSessionFailed(fetch, session, logger);
			if (errorResponse) return errorResponse;
		}

		return json({ received: true });
	} catch (error) {
		return webhookErrorResponse(error, logger);
	}
};
