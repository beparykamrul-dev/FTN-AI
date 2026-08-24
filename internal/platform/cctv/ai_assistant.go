package cctv

// AssistantIntent is the normalized CCTV assistant command contract.
type AssistantIntent struct { UserID string `json:"user_id"`; CameraID string `json:"camera_id"`; Action string `json:"action"`; Start string `json:"start,omitempty"`; End string `json:"end,omitempty"` }

type AlertRule struct { CameraID string `json:"camera_id"`; Event string `json:"event"`; Enabled bool `json:"enabled"`; Severity string `json:"severity"` }

func AllowedActions() []string { return []string{"live_view","playback","health","ack_alert","search_event"} }
