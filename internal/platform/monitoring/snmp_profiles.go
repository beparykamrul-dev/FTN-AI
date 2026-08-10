package monitoring

import "time"

type SNMPProfile struct {
	Name string `json:"name"`
	Version string `json:"version"`
	OIDs []string `json:"oids"`
	Timeout time.Duration `json:"timeout"`
	Retries int `json:"retries"`
}

var DefaultSNMPProfiles = map[string]SNMPProfile{
	"generic-interface": {
		Name: "generic-interface",
		Version: "2c",
		OIDs: []string{
			"1.3.6.1.2.1.2.2.1.2",
			"1.3.6.1.2.1.2.2.1.5",
			"1.3.6.1.2.1.2.2.1.8",
			"1.3.6.1.2.1.2.2.1.10",
			"1.3.6.1.2.1.2.2.1.16",
		},
		Timeout: 3 * time.Second,
		Retries: 2,
	},
}

func GetSNMPProfile(name string) (SNMPProfile, bool) {
	p, ok := DefaultSNMPProfiles[name]
	return p, ok
}
