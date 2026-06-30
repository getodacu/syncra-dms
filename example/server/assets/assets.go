package assets

import "embed"

const InvoiceTemplatePath = "invoice_template.html"

// FS contains runtime assets bundled into the syncra-server binary.
//
//go:embed invoice_template.html
var FS embed.FS
