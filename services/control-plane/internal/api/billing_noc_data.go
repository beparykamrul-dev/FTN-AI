package api

type BillingSummary struct {
	InvoiceCount int64 `json:"invoice_count"`
}

type NOCSummary struct {
	NodeCount int64 `json:"node_count"`
}
