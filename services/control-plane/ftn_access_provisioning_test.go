package main

import "testing"

func TestValidateFTNSubscriberAccessIntent(t *testing.T) {
	in := FTNSubscriberAccessIntent{SubscriberID:"sub-1", OLTID:"olt-1", ONUID:"onu-1", MikroTikID:"mt-1", CustomerVLAN:100, OuterVLAN:200, ServiceVLAN:300, VRF:"internet", Service:"internet", Approved:true}
	if err := ValidateFTNSubscriberAccessIntent(in); err != nil { t.Fatal(err) }
}

func TestValidateFTNSubscriberAccessRequiresApproval(t *testing.T) {
	in := FTNSubscriberAccessIntent{SubscriberID:"sub-1", OLTID:"olt-1", ONUID:"onu-1", MikroTikID:"mt-1", CustomerVLAN:100, VRF:"internet", Service:"internet"}
	if err := ValidateFTNSubscriberAccessIntent(in); err == nil { t.Fatal("expected approval requirement") }
}
