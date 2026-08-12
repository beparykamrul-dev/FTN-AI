package model

// Dependency represents an authorized relationship in the FTN service graph.
type Dependency struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Kind      string `json:"kind"`
	Critical  bool   `json:"critical"`
	Healthy   bool   `json:"healthy"`
}
