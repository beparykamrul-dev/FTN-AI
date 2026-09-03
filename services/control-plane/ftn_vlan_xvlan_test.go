package main

import "testing"

func TestNormalizeFTNVLAN(t *testing.T) {
	v, err := NormalizeFTNVLAN(FTNVLAN{ID:100, Name:" access-users ", Mode:"ACCESS", Enabled:true})
	if err != nil { t.Fatal(err) }
	if v.ID != 100 || v.Mode != FTNVLANAccess || v.Name != "access-users" { t.Fatalf("vlan=%+v", v) }
}

func TestNormalizeFTNQinQ(t *testing.T) {
	v, err := NormalizeFTNVLAN(FTNVLAN{ID:200, Mode:FTNVLANQinQ, OuterID:200, InnerID:100})
	if err != nil { t.Fatal(err) }
	if v.OuterID != 200 || v.InnerID != 100 { t.Fatalf("vlan=%+v", v) }
}

func TestNormalizeFTNVLANRejectsInvalidID(t *testing.T) {
	if _, err := NormalizeFTNVLAN(FTNVLAN{ID:4095}); err == nil { t.Fatal("expected invalid VLAN id") }
}

func TestNormalizeFTNVLANBinding(t *testing.T) {
	b, err := NormalizeFTNVLANBinding(FTNVLANBinding{VLANID:100, DeviceID:"olt-1", Interface:"ether1", Service:"internet", IPPrefix:"192.0.2.7/24"})
	if err != nil { t.Fatal(err) }
	if b.IPPrefix != "192.0.2.0/24" { t.Fatalf("binding=%+v", b) }
}
