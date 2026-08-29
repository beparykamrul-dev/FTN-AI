package codec

// Capability describes an implementation-neutral processing capability.
// Implementations are selected by the FTN control plane; the core dataplane
// does not depend on a particular encoder or codec implementation.
type Capability struct {
	ID             string   `json:"id"`
	Class          string   `json:"class"`
	Modes          []string `json:"modes,omitempty"`
	Hardware       bool     `json:"hardware,omitempty"`
	Lossless       bool     `json:"lossless,omitempty"`
	WorkerIsolated bool     `json:"worker_isolated,omitempty"`
}

// Job is the common contract for codec/media workers.
type Job struct {
	CapabilityID string            `json:"capability_id"`
	InputURI     string            `json:"input_uri"`
	OutputURI    string            `json:"output_uri,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
	Preserve     bool              `json:"preserve_original"`
}

// Result is returned only after the worker has verified its output.
type Result struct {
	JobID        string `json:"job_id"`
	Status       string `json:"status"`
	OutputURI    string `json:"output_uri,omitempty"`
	InputSHA256  string `json:"input_sha256,omitempty"`
	OutputSHA256 string `json:"output_sha256,omitempty"`
	BytesIn      int64  `json:"bytes_in,omitempty"`
	BytesOut     int64  `json:"bytes_out,omitempty"`
}

// DefaultCapabilities is intentionally small and provider-neutral. External
// projects are registered as replaceable workers rather than copied into the
// control plane.
func DefaultCapabilities() []Capability {
	return []Capability{
		{ID: "binary-framing", Class: "transport", Modes: []string{"stream", "datagram"}},
		{ID: "compression", Class: "transfer", Modes: []string{"lossless"}, Lossless: true},
		{ID: "chunking", Class: "transfer", Modes: []string{"resumable", "parallel"}, Lossless: true},
		{ID: "deduplication", Class: "transfer", Modes: []string{"content-addressed"}, Lossless: true},
		{ID: "hardware-video-encode", Class: "media", Modes: []string{"h265"}, Hardware: true, WorkerIsolated: true},
		{ID: "media-cut", Class: "media", Modes: []string{"video-processing"}, WorkerIsolated: true},
	}
}
