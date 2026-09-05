package callcenter

import("strings";"time")
type Channel string
const(ChannelVoice Channel="voice";ChannelWeb Channel="web";ChannelAndroid Channel="android")
type CallSession struct{ID string `json:"id"`;CustomerID string `json:"customer_id"`;Channel Channel `json:"channel"`;Language string `json:"language,omitempty"`;Status string `json:"status"`;StartedAt time.Time `json:"started_at"`;EndedAt *time.Time `json:"ended_at,omitempty"`;AssignedAgent string `json:"assigned_agent,omitempty"`}
type CallEvent struct{SessionID string `json:"session_id"`;Type string `json:"type"`;Timestamp time.Time `json:"timestamp"`;Payload map[string]any `json:"payload,omitempty"`}
type AIActionPolicy struct{CanReadCustomerProfile bool `json:"can_read_customer_profile"`;CanCreateTicket bool `json:"can_create_ticket"`;CanChangeService bool `json:"can_change_service"`;CanExecuteNetworkAction bool `json:"can_execute_network_action"`;RequiresApproval bool `json:"requires_approval"`}
type CallCenterConfig struct{Policy AIActionPolicy `json:"policy"`;MaxSessionMinutes int `json:"max_session_minutes"`}
func(s CallSession)Valid()bool{return strings.TrimSpace(s.ID)!=""&&strings.TrimSpace(s.CustomerID)!=""&&(s.Channel==ChannelVoice||s.Channel==ChannelWeb||s.Channel==ChannelAndroid)&&!s.StartedAt.IsZero()}
