package gis

import "encoding/json"
import "net/http"

type API struct { IPAM *IPAM; Map *MapStore }

func (a *API) IPAssets(w http.ResponseWriter, r *http.Request) {
	if w==nil||r==nil{return}
	if r.Method!=http.MethodGet{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
	if a==nil||a.IPAM==nil{http.Error(w,"GIS IPAM not configured",http.StatusServiceUnavailable);return}
	w.Header().Set("Content-Type","application/json")
	if err:=json.NewEncoder(w).Encode(a.IPAM.List());err!=nil{return}
}

func (a *API) MapSnapshot(w http.ResponseWriter, r *http.Request) {
	if w==nil||r==nil{return}
	if r.Method!=http.MethodGet{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
	if a==nil||a.Map==nil{http.Error(w,"GIS map not configured",http.StatusServiceUnavailable);return}
	w.Header().Set("Content-Type","application/json")
	_ = json.NewEncoder(w).Encode(a.Map.Snapshot())
}
