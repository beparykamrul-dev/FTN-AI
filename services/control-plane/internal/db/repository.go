package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRequest struct {
	ServiceID string
	DeviceBrand string
	Model string
	MAC string
	Serial string
	Scope string
}

type Repository struct { Pool *pgxpool.Pool }

func (r *Repository) SaveServiceRequest(ctx context.Context, q ServiceRequest) error {
	if r == nil || r.Pool == nil { return nil }
	_, err := r.Pool.Exec(ctx, `INSERT INTO service_requests(service_id,device_brand,model,mac,serial,scope) VALUES($1,$2,$3,$4,$5,$6)`, q.ServiceID,q.DeviceBrand,q.Model,q.MAC,q.Serial,q.Scope)
	return err
}
