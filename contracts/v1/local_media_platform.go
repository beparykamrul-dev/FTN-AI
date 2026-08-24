package contracts

// LocalMediaPlatform defines a production-oriented local-first media plane.
type LocalMediaPlatform struct {
    OriginNodes       []string `json:"origin_nodes"`
    EdgeNodes         []string `json:"edge_nodes"`
    CacheNodes        []string `json:"cache_nodes"`
    StorageTargets    []string `json:"storage_targets"`
    PlaybackProtocols []string `json:"playback_protocols"`
    IngestProtocols   []string `json:"ingest_protocols"`
    ABRProfiles       []string `json:"abr_profiles"`
    HealthChecks      []string `json:"health_checks"`
    FailoverEnabled   bool     `json:"failover_enabled"`
    LocalFirst        bool     `json:"local_first"`
}

// Recommended production profiles are capability profiles, not promises about
// source quality. A source can only be delivered at qualities it actually has.
var DefaultMediaProfiles = []string{
    "source",
    "360p",
    "480p",
    "720p",
    "1080p",
    "1440p",
    "2160p",
    "4320p",
}
