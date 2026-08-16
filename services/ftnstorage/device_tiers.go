package ftnstorage

// DeviceClass represents a physical resource tier used by FTN storage.
type DeviceClass string

const (
	RAM   DeviceClass = "ram"
	SSD   DeviceClass = "ssd"
	NVMe  DeviceClass = "nvme"
	HDD   DeviceClass = "hdd"
)

// DeviceProfile describes capacity and performance characteristics without
// assuming a particular vendor or filesystem.
type DeviceProfile struct {
	NodeID      string
	Class       DeviceClass
	Capacity    uint64
	Used        uint64
	LatencyMS   uint32
	Healthy     bool
	Writable    bool
}

func (d DeviceProfile) Valid() bool {
	if d.NodeID == "" || d.Capacity == 0 || d.Used > d.Capacity || !d.Healthy {
		return false
	}
	switch d.Class {
	case RAM, SSD, NVMe, HDD:
		return true
	default:
		return false
	}
}
