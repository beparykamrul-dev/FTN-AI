package edge

// FastlyInspiredCapabilities captures edge-delivery capabilities FTN can
// implement independently: shielding, fanout, segmented cache, service
// discovery, telemetry and local emulation.
type FastlyInspiredCapabilities struct {
    OriginShielding   bool `json:"origin_shielding"`
    Fanout            bool `json:"fanout"`
    WebSocket         bool `json:"websocket"`
    SSE               bool `json:"sse"`
    LongPolling       bool `json:"long_polling"`
    SegmentedCaching  bool `json:"segmented_caching"`
    ServiceDiscovery  bool `json:"service_discovery"`
    EdgeTelemetry     bool `json:"edge_telemetry"`
    TracePropagation  bool `json:"trace_propagation"`
    EdgeLogAnalytics  bool `json:"edge_log_analytics"`
    MediaTelemetry    bool `json:"media_telemetry"`
    LocalEmulation    bool `json:"local_emulation"`
}

func ProductionCapabilities() FastlyInspiredCapabilities {
    return FastlyInspiredCapabilities{
        OriginShielding:true, Fanout:true, WebSocket:true, SSE:true,
        LongPolling:true, SegmentedCaching:true, ServiceDiscovery:true,
        EdgeTelemetry:true, TracePropagation:true, EdgeLogAnalytics:true,
        MediaTelemetry:true, LocalEmulation:true,
    }
}
