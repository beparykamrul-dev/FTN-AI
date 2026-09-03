package main

import (
	"context"
	"fmt"
	"strings"
)

type FTNSubscriberAccessIntent struct {
	SubscriberID string `json:"subscriber_id"`
	OLTID string `json:"olt_id"`
	ONUID string `json:"onu_id"`
	MikroTikID string `json:"mikrotik_id"`
	CustomerVLAN uint16 `json:"customer_vlan"`
	OuterVLAN uint16 `json:"outer_vlan,omitempty"`
	ServiceVLAN uint16 `json:"service_vlan,omitempty"`
	VRF string `json:"vrf"`
	Service string `json:"service"`
	Approved bool `json:"approved"`
}

func ValidateFTNSubscriberAccessIntent(in FTNSubscriberAccessIntent) error {
	for name, value := range map[string]string{"subscriber_id":in.SubscriberID,"olt_id":in.OLTID,"onu_id":in.ONUID,"mikrotik_id":in.MikroTikID,"service":in.Service} {
		if strings.TrimSpace(value)=="" { return fmt.Errorf("%s is required", name) }
	}
	if in.CustomerVLAN == 0 || in.CustomerVLAN > 4094 { return fmt.Errorf("invalid customer VLAN") }
	if in.OuterVLAN > 4094 || in.ServiceVLAN > 4094 { return fmt.Errorf("invalid outer or service VLAN") }
	if strings.TrimSpace(in.VRF)=="" { return fmt.Errorf("VRF is required") }
	if !in.Approved { return fmt.Errorf("access provisioning approval required") }
	return nil
}

// FTNAccessProvisioner is the device-neutral boundary for approved
// subscriber provisioning. Concrete OLT/ONU/MikroTik adapters implement it.
type FTNAccessProvisioner interface {
	Provision(context.Context, FTNSubscriberAccessIntent) error
	Verify(context.Context, FTNSubscriberAccessIntent) error
	Rollback(context.Context, FTNSubscriberAccessIntent) error
}
