package edge

import("context";"math";"net/netip";"strings")
type RouterState struct { ID string `json:"id"`; Hostname string `json:"hostname"`; ManagementIP string `json:"management_ip"`; Status string `json:"status"`; Interfaces []InterfaceState `json:"interfaces,omitempty"` }
type InterfaceState struct { Name string `json:"name"`; Address string `json:"address,omitempty"`; State string `json:"state"`; RXBPS float64 `json:"rx_bps,omitempty"`; TXBPS float64 `json:"tx_bps,omitempty"` }
type CoreRouter interface { Identity(context.Context)(RouterState,error); Interfaces(context.Context)([]InterfaceState,error); Apply(context.Context,ChangeRequest)error }
type ChangeRequest struct { ID string `json:"id"`; Action string `json:"action"`; Target string `json:"target"`; Reason string `json:"reason"`; ApprovedBy string `json:"approved_by"` }
func (r RouterState) Valid() bool { return strings.TrimSpace(r.ID)!=""&&len(strings.TrimSpace(r.ID))<=256&&len(strings.TrimSpace(r.Hostname))<=256&&len(strings.TrimSpace(r.ManagementIP))<=64&&(r.ManagementIP==""||func()bool{_,e:=netip.ParseAddr(strings.TrimSpace(r.ManagementIP));return e==nil}())&&len(r.Interfaces)<=10000 }
func (i InterfaceState) Valid() bool { return strings.TrimSpace(i.Name)!=""&&len(strings.TrimSpace(i.Name))<=256&&len(i.State)<=64&&!math.IsNaN(i.RXBPS)&&!math.IsInf(i.RXBPS,0)&&i.RXBPS>=0&&!math.IsNaN(i.TXBPS)&&!math.IsInf(i.TXBPS,0)&&i.TXBPS>=0 }
func (c ChangeRequest) Valid() bool { return strings.TrimSpace(c.ID)!=""&&len(strings.TrimSpace(c.ID))<=256&&strings.TrimSpace(c.Action)!=""&&len(strings.TrimSpace(c.Action))<=64&&strings.TrimSpace(c.Target)!=""&&len(strings.TrimSpace(c.Target))<=512&&len(strings.TrimSpace(c.Reason))<=2048&&strings.TrimSpace(c.ApprovedBy)!=""&&len(strings.TrimSpace(c.ApprovedBy))<=256 }
