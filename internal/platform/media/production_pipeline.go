package media

import (
 "fmt"
 "sort"
)

type Variant struct { Name string `json:"name"`; Width int `json:"width"`; Height int `json:"height"`; BitrateKbps int `json:"bitrate_kbps"`; Codec string `json:"codec"` }
type Pipeline struct { ItemID string `json:"item_id"`; Origin string `json:"origin"`; Variants []Variant `json:"variants"`; Live bool `json:"live"`; Playback bool `json:"playback"`; CacheEnabled bool `json:"cache_enabled"`; Health string `json:"health"` }

// BuildPipeline creates a provider-neutral media delivery plan. It does not
// fabricate source quality: variants above the source capability are omitted.
func BuildPipeline(item Item, sourceWidth, sourceHeight int) (Pipeline,error) {
 if item.ID=="" { return Pipeline{},fmt.Errorf("media item id is required") }
 if sourceWidth<=0 || sourceHeight<=0 { return Pipeline{},fmt.Errorf("source dimensions are required") }
 candidates:=[]Variant{{"360p",640,360,800,"H264"},{"720p",1280,720,2500,"H264"},{"1080p",1920,1080,5000,"H264"},{"4K",3840,2160,12000,"H265"},{"8K",7680,4320,30000,"H265"}}
 out:=Pipeline{ItemID:item.ID,Origin:item.OriginNode,Live:true,Playback:true,CacheEnabled:true,Health:"unknown"}
 for _,v:=range candidates { if v.Width<=sourceWidth && v.Height<=sourceHeight { out.Variants=append(out.Variants,v) } }
 sort.Slice(out.Variants,func(i,j int)bool{return out.Variants[i].Height<out.Variants[j].Height})
 return out,nil
}
