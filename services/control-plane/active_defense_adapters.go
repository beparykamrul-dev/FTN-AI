package main

import "strings"
type ActiveDefenseAdapter string
const(AdapterNFTables ActiveDefenseAdapter="nftables";AdapterXDP ActiveDefenseAdapter="xdp";AdapterMikroTik ActiveDefenseAdapter="mikrotik")
type ActiveDefenseActionPlan struct{Adapter ActiveDefenseAdapter `json:"adapter"`;Operation string `json:"operation"`;TargetAsset string `json:"target_asset"`;DurationSeconds int `json:"duration_seconds"`;SnapshotRequired bool `json:"snapshot_required"`;VerificationNeeded bool `json:"verification_needed"`;Automatic bool `json:"automatic"`;RequiresApproval bool `json:"requires_approval"`}
func BuildActiveDefenseActionPlan(intent ActiveDefenseExecutionIntent)(ActiveDefenseActionPlan,bool){asset:=strings.TrimSpace(intent.TargetAsset);if intent.TargetScope!="ftn-owned-asset"||asset==""||len(asset)>256{return ActiveDefenseActionPlan{},false};if intent.DurationSeconds<=0||intent.DurationSeconds>3600{return ActiveDefenseActionPlan{},false};plan:=ActiveDefenseActionPlan{TargetAsset:asset,DurationSeconds:intent.DurationSeconds,SnapshotRequired:true,VerificationNeeded:true,Automatic:intent.Automatic,RequiresApproval:true};switch intent.Action{case WazuhHealthRecover:plan.Adapter,plan.Operation=AdapterNFTables,"temporary-containment";case WazuhAlertAction:plan.Adapter,plan.Operation=AdapterNFTables,"rate-limit";default:return ActiveDefenseActionPlan{},false};return plan,true}
