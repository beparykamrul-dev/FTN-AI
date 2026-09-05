package contracts

import "strings"

// LocalMediaItem is a provider-neutral metadata contract for authorized FTN-local media.
type LocalMediaItem struct{ID string `json:"id"`;Title string `json:"title"`;Type string `json:"type"`;Region string `json:"region,omitempty"`;Source string `json:"source,omitempty"`;OriginNode string `json:"origin_node,omitempty"`;CachePolicy string `json:"cache_policy,omitempty"`;DRM string `json:"drm,omitempty"`;Status string `json:"status"`}
type LocalMediaProfile struct{Live bool `json:"live"`;Playback bool `json:"playback"`;Adaptive bool `json:"adaptive"`;LocalCache bool `json:"local_cache"`;Search bool `json:"search"`}
func(m LocalMediaItem)Valid()bool{return strings.TrimSpace(m.ID)!=""&&len(strings.TrimSpace(m.ID))<=256&&strings.TrimSpace(m.Title)!=""&&len(strings.TrimSpace(m.Title))<=512&&strings.TrimSpace(m.Type)!=""&&len(strings.TrimSpace(m.Type))<=64&&len(m.Region)<=128&&len(m.Source)<=2048&&len(m.OriginNode)<=256&&len(m.CachePolicy)<=128&&len(m.DRM)<=128&&len(m.Status)<=64}
