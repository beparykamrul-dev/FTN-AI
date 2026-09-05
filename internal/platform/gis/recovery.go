package gis

import("math";"strings";"time")
type RecoveryAction struct{ID string `json:"id"`;AssetID string `json:"asset_id"`;Kind string `json:"kind"`;Reason string `json:"reason"`;Confidence float64 `json:"confidence"`;RequiresApproval bool `json:"requires_approval"`;CreatedAt time.Time `json:"created_at"`}
type RecoveryEngine struct{}
func NewRecoveryEngine()*RecoveryEngine{return &RecoveryEngine{}}
func(e *RecoveryEngine)Propose(asset FiberAsset,reason string,confidence float64)RecoveryAction{if math.IsNaN(confidence)||math.IsInf(confidence,0){confidence=0}else if confidence<0{confidence=0}else if confidence>1{confidence=1};return RecoveryAction{AssetID:strings.TrimSpace(asset.ID),Kind:strings.TrimSpace(string(asset.Type)),Reason:strings.TrimSpace(reason),Confidence:confidence,RequiresApproval:true,CreatedAt:time.Now().UTC()}}
