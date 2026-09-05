package gis

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
)

type FiberAssetType string
const ( FiberRoute FiberAssetType="fiber_route"; FiberJoint FiberAssetType="joint"; FiberSplitter FiberAssetType="splitter"; FiberCut FiberAssetType="cut"; FiberDamage FiberAssetType="damage"; FiberPPPoE FiberAssetType="pppoe"; FiberONU FiberAssetType="onu"; FiberRouter FiberAssetType="router"; FiberUser FiberAssetType="user" )
type FiberPoint struct { Lat float64 `json:"lat"`; Lng float64 `json:"lng"` }
type FiberAsset struct { ID string `json:"id"`; Type FiberAssetType `json:"type"`; Name string `json:"name"`; Points []FiberPoint `json:"points,omitempty"`; ParentID string `json:"parent_id,omitempty"`; Status string `json:"status,omitempty"`; DistanceM float64 `json:"distance_m,omitempty"`; DatabaseRef string `json:"database_ref,omitempty"` }
type FiberMap struct { mu sync.RWMutex; assets map[string]FiberAsset }
func NewFiberMap()*FiberMap{return &FiberMap{assets:make(map[string]FiberAsset)}}
func(m *FiberMap)Upsert(a FiberAsset)error{if m==nil{return fmt.Errorf("fiber map is nil")};a.ID=strings.TrimSpace(a.ID);a.Name=strings.TrimSpace(a.Name);if a.ID==""||a.Name==""{return fmt.Errorf("fiber asset id and name are required")};points:=append([]FiberPoint(nil),a.Points...);for _,p:=range points{if math.IsNaN(p.Lat)||math.IsNaN(p.Lng)||math.IsInf(p.Lat,0)||math.IsInf(p.Lng,0)||p.Lat < -90||p.Lat>90||p.Lng < -180||p.Lng>180{return fmt.Errorf("invalid fiber point")}};a.Points=points;if len(a.Points)>0{a.DistanceM=routeDistance(a.Points)};m.mu.Lock();m.assets[a.ID]=a;m.mu.Unlock();return nil}
func(m *FiberMap)List()[]FiberAsset{if m==nil{return []FiberAsset{}};m.mu.RLock();out:=make([]FiberAsset,0,len(m.assets));for _,a:=range m.assets{a.Points=append([]FiberPoint(nil),a.Points...);out=append(out,a)};m.mu.RUnlock();sort.SliceStable(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out}
func routeDistance(p []FiberPoint)float64{var d float64;for i:=1;i<len(p);i++{d+=haversine(p[i-1],p[i])};return d}
func haversine(a,b FiberPoint)float64{const R=6371000.0;lat1:=a.Lat*math.Pi/180;lat2:=b.Lat*math.Pi/180;dlat:=lat2-lat1;dlng:=(b.Lng-a.Lng)*math.Pi/180;h:=math.Sin(dlat/2)*math.Sin(dlat/2)+math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlng/2)*math.Sin(dlng/2);if h<0{h=0};if h>1{h=1};return 2*R*math.Asin(math.Sqrt(h))}
