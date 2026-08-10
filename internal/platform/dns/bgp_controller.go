package dns

import (
    "fmt"
    "sync"
)

type BGPController struct {
    mu sync.RWMutex
    advertisements map[string]BGPAdvertisement
    enabled bool
}

func NewBGPController() *BGPController {
    return &BGPController{advertisements: make(map[string]BGPAdvertisement)}
}

func (c *BGPController) SetEnabled(enabled bool) {
    c.mu.Lock(); defer c.mu.Unlock(); c.enabled = enabled
}

func (c *BGPController) Upsert(a BGPAdvertisement) error {
    if err := ValidateBGPAdvertisement(a); err != nil { return err }
    c.mu.Lock(); defer c.mu.Unlock()
    c.advertisements[a.NodeID+":"+a.Prefix] = a
    return nil
}

func (c *BGPController) Withdraw(nodeID, prefix string) {
    c.mu.Lock(); defer c.mu.Unlock()
    delete(c.advertisements, nodeID+":"+prefix)
}

func (c *BGPController) Announcements() ([]BGPAdvertisement, error) {
    c.mu.RLock(); defer c.mu.RUnlock()
    if !c.enabled { return nil, fmt.Errorf("BGP controller is disabled") }
    items := make([]BGPAdvertisement, 0, len(c.advertisements))
    for _, a := range c.advertisements { items = append(items, a) }
    return SelectAdvertisements(items), nil
}
