package db

import "context"

type BillingStore interface { CountInvoices(context.Context) (int64, error) }
type NOCStore interface { CountNodes(context.Context) (int64, error) }
type ServiceRequestStore interface { SaveServiceRequest(context.Context, ServiceRequest) error }
