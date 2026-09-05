package agent

import (
	"strings"
	"time"
)

type DeviceKind string
const ( DeviceServer DeviceKind="server"; DeviceRouter DeviceKind="router"; DeviceOLT DeviceKind="olt"; DeviceONU DeviceKind="onu"; DevicePC DeviceKind="pc"; DeviceAndroid DeviceKind="android"; DeviceTV DeviceKind="tv"; DeviceVirtual DeviceKind="virtual" )
type DeviceIdentity struct { ID string `json:"id"`; Kind DeviceKind `json:"kind"`; Name string `json:"name"`; Hostname string `json:"hostname,omitempty"`; IP string `json:"ip,omitempty"`; MAC string `json:"mac,omitempty"`; Serial string `json:"serial,omitempty"`; OS string `json:"os,omitempty"`; AgentID string `json:"agent_id,omitempty"`; Online bool `json:"online"`; LastSeen time.Time `json:"last_seen,omitempty"` }
func(d DeviceIdentity)HasStableIdentity()bool{return strings.TrimSpace(d.ID)!=""||strings.TrimSpace(d.AgentID)!=""||strings.TrimSpace(d.Serial)!=""||strings.TrimSpace(d.MAC)!=""}
func(d DeviceIdentity)Normalized()DeviceIdentity{d.ID=strings.TrimSpace(d.ID);d.Name=strings.TrimSpace(d.Name);d.Hostname=strings.TrimSpace(d.Hostname);d.IP=strings.TrimSpace(d.IP);d.MAC=strings.ToLower(strings.TrimSpace(d.MAC));d.Serial=strings.TrimSpace(d.Serial);d.OS=strings.TrimSpace(d.OS);d.AgentID=strings.TrimSpace(d.AgentID);if !d.LastSeen.IsZero(){d.LastSeen=d.LastSeen.UTC()};return d}
