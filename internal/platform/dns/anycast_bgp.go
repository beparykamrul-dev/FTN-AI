package dns

import (
    "fmt"
    "net"
    "sort"
    "strings"
)

type BGPAdvertisement struct {
    Prefix string `json:"prefix"`
    NodeID string `json:"node_id"`
    Enabled bool `json:"enabled"`
    Healthy bool `json:"healthy"`
    LocalPreference uint32 `json:"local_preference"`
    Community []string `json:"community,omitempty"`
}

func ValidateBGPAdvertisement(a BGPAdvertisement) error {
    if strings.TrimSpace(a.NodeID) == "" { return fmt.Errorf("node ID is required") }
    _, network, err := net.ParseCIDR(strings.TrimSpace(a.Prefix))
    if err != nil || network == nil { return fmt.Errorf("invalid BGP prefix: %s", a.Prefix) }
    return nil
}

func SelectAdvertisements(items []BGPAdvertisement) []BGPAdvertisement {
    out := make([]BGPAdvertisement, 0, len(items))
    for _, a := range items {
        if a.Enabled && a.Healthy && ValidateBGPAdvertisement(a) == nil { out = append(out, a) }
    }
    sort.SliceStable(out, func(i, j int) bool {
        if out[i].LocalPreference == out[j].LocalPreference { return out[i].NodeID < out[j].NodeID }
        return out[i].LocalPreference > out[j].LocalPreference
    })
    return out
}
