package deploy

import (
	"errors"
	"net"
	"strings"
)

type Target struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	IP        string            `json:"ip"`
	MAC       string            `json:"mac,omitempty"`
	Serial    string            `json:"serial,omitempty"`
	OS        string            `json:"os,omitempty"`
	AgentID   string            `json:"agent_id,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Online    bool              `json:"online"`
}

func (t Target) Validate() error {
	if t.ID == "" || t.Name == "" { return errors.New("target id and name are required") }
	if t.IP != "" && net.ParseIP(t.IP) == nil { return errors.New("invalid target IP") }
	if t.MAC != "" {
		if _, err := net.ParseMAC(t.MAC); err != nil { return errors.New("invalid target MAC") }
	}
	if strings.TrimSpace(t.Serial) == "" && strings.TrimSpace(t.AgentID) == "" {
		return errors.New("serial or authenticated agent id is required")
	}
	return nil
}
