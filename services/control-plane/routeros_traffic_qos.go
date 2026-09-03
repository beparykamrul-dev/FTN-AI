package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RouterOSTrafficQoSRule struct { ServiceID string `json:"service_id"`; Class TrafficClass `json:"class"`; DSCP uint8 `json:"dscp"`; Priority uint8 `json:"priority"`; PathID string `json:"path_id"` }
type RouterOSTrafficQoSPlan struct { DeviceID string `json:"device_id"`; Rules []RouterOSTrafficQoSRule `json:"rules"`; RequiresApproval bool `json:"requires_approval"` }

func BuildRouterOSTrafficQoSPlan(device NetworkDevice, decisions []TrafficDecision) (RouterOSTrafficQoSPlan,error) {
	if strings.TrimSpace(device.ID)=="" || !device.Healthy || !isFTNDeviceKind(device.Kind) { return RouterOSTrafficQoSPlan{},errors.New("invalid FTN RouterOS target") }
	out:=RouterOSTrafficQoSPlan{DeviceID:device.ID,RequiresApproval:true};seen:=make(map[string]struct{},len(decisions))
	for _,d:=range decisions { if strings.TrimSpace(d.ServiceID)==""||strings.TrimSpace(d.PathID)==""||d.DSCP>63{continue};if _,ok:=trafficPolicyByID(d.ServiceID);!ok{continue};key:=d.ServiceID+"\x00"+d.PathID;if _,ok:=seen[key];ok{continue};seen[key]=struct{}{};out.Rules=append(out.Rules,RouterOSTrafficQoSRule{ServiceID:d.ServiceID,Class:d.Class,DSCP:d.DSCP,Priority:d.Priority,PathID:d.PathID}) }
	if len(out.Rules)==0{return RouterOSTrafficQoSPlan{},errors.New("no valid traffic decisions")};return out,nil
}
func RouterOSTrafficQoSAction(device NetworkDevice) NetworkExecutionIntent{return NetworkExecutionIntent{Device:device,Action:"configure traffic-qos",Approved:false,Explicit:true,PrechangeSnapshot:true,VerificationRequired:true,Timeout:30*time.Second}}
func RenderRouterOSTrafficQoSComment(rule RouterOSTrafficQoSRule) string{return fmt.Sprintf("FTN-QOS service=%s class=%s path=%s dscp=%d priority=%d",rule.ServiceID,rule.Class,rule.PathID,rule.DSCP,rule.Priority)}
