// Package api Syncra DMS API
//
// Syncra DMS backend API.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package api

func swaggerOperations() {
	// swagger:operation GET /healthz system getHealth
	//
	// Health check.
	//
	// ---
	// responses:
	//   "200":
	//     description: API process is running.

	// swagger:operation GET /readyz system getReadiness
	//
	// Readiness check.
	//
	// ---
	// responses:
	//   "200":
	//     description: API dependencies are ready.
	//   "503":
	//     description: API dependencies are not ready.

	// swagger:operation GET /version system getVersion
	//
	// Version metadata.
	//
	// ---
	// responses:
	//   "200":
	//     description: API version metadata.

	// swagger:operation POST /api/auth/sign-up/email auth signUpEmail
	//
	// Sign up with email and password.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: User account was created or already exists and email verification is required.
	//   "400":
	//     description: Invalid sign-up request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/sign-in/email auth signInEmail
	//
	// Sign in with email and password.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: Invalid email or password, or trusted internal request required.
	//   "403":
	//     description: Email is not verified.

	// swagger:operation GET /api/auth/get-session auth getSession
	//
	// Load the current session from the auth.session_token cookie.
	//
	// Trusted SvelteKit server hook endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Current authenticated session, or null when no valid session exists.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/sign-out auth signOut
	//
	// Sign out the current session.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Session was deleted when present.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/email-otp/send-verification-otp auth sendVerificationOTP
	//
	// Send or rotate a six-digit email verification OTP.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: OTP was created when the account exists and is unverified.
	//   "400":
	//     description: Invalid OTP request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/email-otp/verify-email auth verifyEmailOTP
	//
	// Confirm an email address with a six-digit OTP.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Email was verified.
	//   "400":
	//     description: Invalid or expired verification code.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/password-reset/request auth requestPasswordReset
	//
	// Request a password reset email.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Reset request accepted.
	//   "400":
	//     description: Invalid password reset request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/password-reset/confirm auth confirmPasswordReset
	//
	// Confirm password reset with emailed token.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Password was reset and existing sessions were revoked.
	//   "400":
	//     description: Invalid or expired password reset token.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/oauth/google/start auth startGoogleOAuth
	//
	// Start Google OAuth.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Google authorization URL and state.
	//   "401":
	//     description: Trusted internal request required.
	//   "503":
	//     description: Google OAuth is not configured.

	// swagger:operation POST /api/auth/oauth/google/callback auth signInGoogleOAuth
	//
	// Complete Google OAuth sign-in.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: OAuth sign-in failed or trusted internal request required.

	// swagger:operation POST /api/auth/oauth/github/start auth startGitHubOAuth
	//
	// Start GitHub OAuth.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: GitHub authorization URL and state.
	//   "401":
	//     description: Trusted internal request required.
	//   "503":
	//     description: GitHub OAuth is not configured.

	// swagger:operation POST /api/auth/oauth/github/callback auth signInGitHubOAuth
	//
	// Complete GitHub OAuth sign-in.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: OAuth sign-in failed or trusted internal request required.
}
