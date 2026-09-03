package dns

// Record is the normalized provider-independent DNS resource record used by
// consistency, drift detection and reconciliation planning.
type Record struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}
