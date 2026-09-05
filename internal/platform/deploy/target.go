package deploy

import (
	"errors"
	"net"
	"strings"
)

type Target struct {
	ID string `json:"id"`; Name string `json:"name"`; IP string `json:"ip"`; MAC string `json:"mac,omitempty"`; Serial string `json:"serial,omitempty"`; OS string `json:"os,omitempty"`; AgentID string `json:"agent_id,omitempty"`; Labels map[string]string `json:"labels,omitempty"`; Online bool `json:"online"`
}

func (t Target) Validate() error {
	t = t.Normalize()
	if t.ID == "" || t.Name == "" { return errors.New("target id and name are required") }
	if len(t.ID) > 256 || len(t.Name) > 256 || len(t.Serial) > 256 || len(t.AgentID) > 256 || len(t.OS) > 128 { return errors.New("target field is too large") }
	if t.IP != "" && net.ParseIP(t.IP) == nil { return errors.New("invalid target IP") }
	if t.MAC != "" { if _, err := net.ParseMAC(t.MAC); err != nil { return errors.New("invalid target MAC") } }
	if t.Serial == "" && t.AgentID == "" { return errors.New("serial or authenticated agent id is required") }
	if len(t.Labels) > 128 { return errors.New("too many target labels") }
	for k, v := range t.Labels { if strings.TrimSpace(k) == "" || len(k) > 128 || len(v) > 1024 { return errors.New("target label is invalid or too large") } }
	return nil
}

func (t Target) Normalize() Target {
	t.ID = strings.TrimSpace(t.ID); t.Name = strings.TrimSpace(t.Name); t.IP = strings.TrimSpace(t.IP); t.MAC = strings.ToLower(strings.TrimSpace(t.MAC)); t.Serial = strings.TrimSpace(t.Serial); t.OS = strings.TrimSpace(t.OS); t.AgentID = strings.TrimSpace(t.AgentID)
	if t.Labels != nil { labels := make(map[string]string, len(t.Labels)); for k, v := range t.Labels { k = strings.TrimSpace(k); if k == "" || len(k) > 128 { continue }; labels[k] = strings.TrimSpace(v) }; t.Labels = labels }
	return t
}
