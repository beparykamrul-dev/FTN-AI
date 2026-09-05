package cctv

import "strings"
type QualityPolicy struct{LiveProfiles []string `json:"live_profiles"`;PlaybackProfiles []string `json:"playback_profiles"`;AllowAdaptiveTranscode bool `json:"allow_adaptive_transcode"`;PreferLocalCache bool `json:"prefer_local_cache"`}
func(q QualityPolicy)Valid()bool{return len(q.LiveProfiles)>0&&len(q.PlaybackProfiles)>0}
func DefaultQualityPolicy()QualityPolicy{return QualityPolicy{LiveProfiles:[]string{"low","720p","1080p","4k","8k"},PlaybackProfiles:[]string{"low","720p","1080p","4k","8k"},AllowAdaptiveTranscode:true,PreferLocalCache:true}}
func NormalizeProfiles(in []string)[]string{out:=make([]string,0,len(in));seen:=map[string]struct{}{};for _,p:=range in{p=strings.ToLower(strings.TrimSpace(p));if p==""{continue};if _,ok:=seen[p];ok{continue};seen[p]=struct{}{};out=append(out,p)};return out}
