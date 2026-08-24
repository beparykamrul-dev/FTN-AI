package cctv

// RecoveryPlan describes a safe recovery action without directly changing
// physical or privileged camera infrastructure.
type RecoveryPlan struct {
    CameraID string `json:"camera_id"`
    Mode Mode `json:"mode"`
    Actions []string `json:"actions"`
    RequiresApproval bool `json:"requires_approval"`
}

func BuildRecoveryPlan(c Camera) RecoveryPlan {
    actions := []string{"refresh health", "reconnect authorized endpoint"}
    if c.Mode == NonIP { actions = append(actions, "validate authorized encoder/DVR bridge") }
    if c.RecorderID != "" { actions = append(actions, "validate recorder path") }
    return RecoveryPlan{CameraID:c.ID, Mode:c.Mode, Actions:actions, RequiresApproval:true}
}
