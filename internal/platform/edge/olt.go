package edge

import("context";"math";"net";"strings")
type OLTState struct{ID string `json:"id"`;Vendor string `json:"vendor"`;Model string `json:"model,omitempty"`;ManagementIP string `json:"management_ip"`;Status string `json:"status"`}
type ONUState struct{ID string `json:"id"`;Serial string `json:"serial"`;PON string `json:"pon"`;Status string `json:"status"`;RXPowerDBM float64 `json:"rx_power_dbm,omitempty"`;TXPowerDBM float64 `json:"tx_power_dbm,omitempty"`}
type OLTDriver interface{Identity(context.Context)(OLTState,error);DiscoverONUs(context.Context)([]ONUState,error);Apply(context.Context,ChangeRequest)error}
func(s OLTState)Valid()bool{id:=strings.TrimSpace(s.ID);vendor:=strings.TrimSpace(s.Vendor);model:=strings.TrimSpace(s.Model);status:=strings.TrimSpace(s.Status);return id!=""&&vendor!=""&&len(id)<=256&&len(vendor)<=128&&len(model)<=256&&len(status)<=64&&net.ParseIP(strings.TrimSpace(s.ManagementIP))!=nil}
func(s ONUState)Valid()bool{id:=strings.TrimSpace(s.ID);serial:=strings.TrimSpace(s.Serial);pon:=strings.TrimSpace(s.PON);status:=strings.TrimSpace(s.Status);return id!=""&&serial!=""&&pon!=""&&len(id)<=256&&len(serial)<=256&&len(pon)<=128&&len(status)<=64&&finitePower(s.RXPowerDBM)&&finitePower(s.TXPowerDBM)}
func finitePower(v float64)bool{return !math.IsNaN(v)&&!math.IsInf(v,0)}
