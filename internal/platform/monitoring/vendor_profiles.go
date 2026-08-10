package monitoring

import "time"

type VendorProfile struct {
	Name string `json:"name"`
	Family string `json:"family"`
	SNMPVersion string `json:"snmp_version"`
	Profiles []string `json:"profiles"`
}

var DefaultVendorProfiles = []VendorProfile{
	{Name: "mikrotik", Family: "router", SNMPVersion: "2c", Profiles: []string{"generic-interface"}},
	{Name: "generic-olt", Family: "olt", SNMPVersion: "2c", Profiles: []string{"generic-interface"}},
}

func VendorProfileFor(name string) (VendorProfile, bool) {
	for _, p := range DefaultVendorProfiles {
		if p.Name == name { return p, true }
	}
	return VendorProfile{}, false
}

type PollSchedule struct { Interval time.Duration `json:"interval"` }
