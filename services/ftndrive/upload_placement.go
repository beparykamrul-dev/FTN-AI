package ftndrive

// UploadObject describes a customer upload before placement.
type UploadObject struct {
	ObjectID string
	Bytes    uint64
	Hot     bool
}

// PlacementAnalyzer converts upload characteristics into a storage class.
type PlacementAnalyzer struct{}

func (PlacementAnalyzer) Classify(o UploadObject) string {
	if o.Hot {
		return "nvme"
	}
	return "auto"
}
