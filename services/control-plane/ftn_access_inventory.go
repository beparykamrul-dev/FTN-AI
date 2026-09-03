package main

import (
	"fmt"
	"net/netip"
	"strings"
)

type FTNOLT struct { ID, Name, Vendor, Model, Address, POPID string; Healthy bool `json:"healthy"` }
type FTNONU struct { ID, Serial, OLTID, Name string; VLAN uint16; ManagementIP string; Online bool `json:"online"` }
type FTNMikroTik struct { ID, Name, Model, Address, POPID string; Healthy bool `json:"healthy"` }

func NormalizeFTNOLT(v FTNOLT) (FTNOLT,error) {
	v.ID=strings.TrimSpace(v.ID); v.Name=strings.TrimSpace(v.Name); v.Address=strings.TrimSpace(v.Address); v.POPID=strings.TrimSpace(v.POPID)
	if v.ID=="" || v.POPID=="" { return FTNOLT{},fmt.Errorf("olt id and pop are required") }
	if v.Address!="" { if _,err:=netip.ParseAddr(v.Address); err!=nil{return FTNOLT{},fmt.Errorf("invalid OLT address: %w",err)} }
	return v,nil
}
func NormalizeFTNONU(v FTNONU) (FTNONU,error) {
	v.ID=strings.TrimSpace(v.ID); v.Serial=strings.TrimSpace(v.Serial); v.OLTID=strings.TrimSpace(v.OLTID); v.ManagementIP=strings.TrimSpace(v.ManagementIP)
	if v.ID=="" || v.OLTID=="" || v.Serial=="" { return FTNONU{},fmt.Errorf("ONU id, serial and OLT are required") }
	if v.ManagementIP!="" { if _,err:=netip.ParseAddr(v.ManagementIP); err!=nil{return FTNONU{},fmt.Errorf("invalid ONU management address: %w",err)} }
	return v,nil
}
func NormalizeFTNMikroTik(v FTNMikroTik) (FTNMikroTik,error) {
	v.ID=strings.TrimSpace(v.ID); v.Address=strings.TrimSpace(v.Address); v.POPID=strings.TrimSpace(v.POPID)
	if v.ID=="" || v.POPID=="" { return FTNMikroTik{},fmt.Errorf("router id and pop are required") }
	if v.Address!="" { if _,err:=netip.ParseAddr(v.Address); err!=nil{return FTNMikroTik{},fmt.Errorf("invalid router address: %w",err)} }
	return v,nil
}
