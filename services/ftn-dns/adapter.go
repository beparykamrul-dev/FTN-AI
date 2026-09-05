package dns

import (
	"context"
	"strings"
)

// Adapter is the common contract for FTN-owned and external DNS backends.
// Implementations must keep credentials outside this interface and return only
// normalized operational data to the FTN DNS control plane.
type Adapter interface {
	Name() string
	Health(ctx context.Context) (Health, error)
	Query(ctx context.Context, name, recordType string) (Response, error)
}

type Health struct {
	Reachable bool
	Secure    bool
	LatencyMS int64
	LossRatio float64
}

type Response struct {
	Name       string
	RecordType string
	Values     []string
	Secure     bool
}

func (h Health) Valid() bool {
	return h.LatencyMS >= 0 && h.LossRatio >= 0 && h.LossRatio <= 1
}

func (r Response) Normalized() Response {
	r.Name = strings.TrimSpace(r.Name)
	r.RecordType = strings.ToUpper(strings.TrimSpace(r.RecordType))
	values := make([]string, 0, len(r.Values))
	for _, value := range r.Values {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	r.Values = values
	return r
}
