package fiber

import "context"

type FieldRecord struct {
	ID string
	EntityID string
	TechnicianID string
	Action string
	Notes string
	Latitude float64
	Longitude float64
	CreatedAt string
}

// FieldRepository keeps field-work evidence tied to topology entities.
type FieldRepository interface {
	Record(context.Context, FieldRecord) error
	History(context.Context, string) ([]FieldRecord, error)
}
