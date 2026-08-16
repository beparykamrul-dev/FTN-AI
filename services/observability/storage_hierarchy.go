package observability

// StorageClass represents a local telemetry buffer tier.
type StorageClass string

const (
	StorageRAM StorageClass = "ram"
	StorageNVMe StorageClass = "nvme"
	StorageSSD StorageClass = "ssd"
	StorageHDD StorageClass = "hdd"
)

// StorageNode describes a candidate tier in the local storage hierarchy.
type StorageNode struct {
	Class StorageClass
	Healthy bool
	FreeGB float64
	LatencyMS float64
	Priority uint32
}

func (n StorageNode) Eligible() bool {
	return n.Healthy && n.FreeGB > 0
}
