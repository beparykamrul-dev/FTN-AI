package contracts

// CameraServiceMode describes how FTN integrates an authorized camera.
type CameraServiceMode string

const (
    CameraIP     CameraServiceMode = "IP_CAMERA"
    CameraNonIP   CameraServiceMode = "NON_IP_CAMERA"
)

// CameraEndpoint keeps normalized camera metadata independent of vendor.
type CameraEndpoint struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Mode        CameraServiceMode `json:"mode"`
    SiteID      string            `json:"site_id"`
    ZoneID      string            `json:"zone_id,omitempty"`
    Address     string            `json:"address,omitempty"`
    Protocol    string            `json:"protocol,omitempty"`
    StreamURL   string            `json:"stream_url,omitempty"`
    GatewayID   string            `json:"gateway_id,omitempty"`
    RecorderID  string            `json:"recorder_id,omitempty"`
    Status      string            `json:"status"`
    AIEnabled   bool              `json:"ai_enabled"`
}

// CameraCapabilities describes normalized capabilities exposed to the control plane.
type CameraCapabilities struct {
    LiveView       bool `json:"live_view"`
    Recording      bool `json:"recording"`
    Playback       bool `json:"playback"`
    PTZ            bool `json:"ptz"`
    MotionEvents   bool `json:"motion_events"`
    AIAnalysis     bool `json:"ai_analysis"`
    LocalStorage   bool `json:"local_storage"`
    CloudStorage   bool `json:"cloud_storage"`
}

// Non-IP cameras are represented through an authorized encoder/DVR bridge;
// they are never treated as native IP endpoints.
type CameraBridge struct {
    ID          string `json:"id"`
    CameraID    string `json:"camera_id"`
    BridgeType  string `json:"bridge_type"`
    Input       string `json:"input"`
    Output      string `json:"output"`
    RecorderID  string `json:"recorder_id,omitempty"`
    Status      string `json:"status"`
}
