package contracts

// LocalMediaItem is a provider-neutral metadata contract for authorized FTN-local media.
type LocalMediaItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	Region     string `json:"region,omitempty"`
	Source     string `json:"source,omitempty"`
	OriginNode string `json:"origin_node,omitempty"`
	CachePolicy string `json:"cache_policy,omitempty"`
	DRM        string `json:"drm,omitempty"`
	Status     string `json:"status"`
}

type LocalMediaProfile struct {
	Live       bool `json:"live"`
	Playback   bool `json:"playback"`
	Adaptive   bool `json:"adaptive"`
	LocalCache bool `json:"local_cache"`
	Search     bool `json:"search"`
}
