package adapters

// DNSSECObservation is the adapter-local normalized result used by policy
// evaluation. Provider/global DNS packages may map their observations into it.
type DNSSECObservation struct {
	Domain string
	Validated bool
	Authenticated bool
	Error string
}
