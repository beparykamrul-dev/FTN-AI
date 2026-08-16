package ftnstorage

// StorageInventory is the normalized physical capacity view used by FTN Drive.
type StorageInventory struct {
	NodeID string
	Devices []DeviceTier
}

// PlacementDecision describes where a new object should be placed.
type PlacementDecision struct {
	NodeID string
	Device DeviceTier
	Reason string
}

// ChoosePlacement prefers the configured device tier and requires a writable,
// healthy device with sufficient capacity.
func ChoosePlacement(inv StorageInventory, desired DeviceClass, bytes uint64) (PlacementDecision, bool) {
	for _, d := range inv.Devices {
		if d.Class == desired && d.Healthy && d.Writable && d.FreeBytes >= bytes {
			return PlacementDecision{NodeID: inv.NodeID, Device: d, Reason: "healthy-tier-capacity"}, true
		}
	}
	return PlacementDecision{}, false
}
