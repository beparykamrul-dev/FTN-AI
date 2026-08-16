package ftnstorage

// DeviceHealth is the normalized health state for RAM/SSD/NVMe/HDD resources.
type DeviceHealth struct {
	Profile      DeviceProfile
	ReadErrors   uint64
	WriteErrors  uint64
	ChecksumBad  uint64
	TemperatureC int32
}

func (h DeviceHealth) HealthyForPlacement() bool {
	return h.Profile.Valid() && h.ReadErrors == 0 && h.WriteErrors == 0 && h.ChecksumBad == 0
}
