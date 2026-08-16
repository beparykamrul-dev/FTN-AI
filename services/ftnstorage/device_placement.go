package ftnstorage

// PlacementForDevice maps data temperature to the preferred physical tier.
func PlacementForDevice(hot bool, warm bool) DeviceClass {
	if hot {
		return NVMe
	}
	if warm {
		return SSD
	}
	return HDD
}

// RAM is reserved as an optional cache tier; persistent data should not rely
// on RAM as the only copy.
func CacheTier() DeviceClass { return RAM }
