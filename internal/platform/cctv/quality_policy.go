package cctv

// QualityPolicy separates source capability from delivery quality. Upscaling
// never claims to create missing source detail; adaptive transcoding may expose
// multiple authorized playback profiles when the source supports them.
type QualityPolicy struct {
    LiveProfiles []string `json:"live_profiles"`
    PlaybackProfiles []string `json:"playback_profiles"`
    AllowAdaptiveTranscode bool `json:"allow_adaptive_transcode"`
    PreferLocalCache bool `json:"prefer_local_cache"`
}

func DefaultQualityPolicy() QualityPolicy {
    return QualityPolicy{LiveProfiles: []string{"low","720p","1080p","4k","8k"}, PlaybackProfiles: []string{"low","720p","1080p","4k","8k"}, AllowAdaptiveTranscode:true, PreferLocalCache:true}
}
