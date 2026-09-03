package adapters

type DNSSECObservation struct {
	Validated     bool   `json:"validated"`
	Authenticated bool   `json:"authenticated"`
	Error         string `json:"error,omitempty"`
}
