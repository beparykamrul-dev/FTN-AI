package fiber

import("encoding/json";"net/http";"strings")
type GISService struct{RecordBuilder func(string)(GISRecord,error)}
func(s GISService)RecordHandler(w http.ResponseWriter,r *http.Request){if w==nil{return};if r==nil||r.Method!=http.MethodGet{w.WriteHeader(http.StatusMethodNotAllowed);return};if s.RecordBuilder==nil{http.Error(w,"GIS service not configured",http.StatusServiceUnavailable);return};id:=strings.TrimSpace(r.URL.Query().Get("node_id"));if id==""||len(id)>256{http.Error(w,"node_id is invalid",http.StatusBadRequest);return};record,err:=s.RecordBuilder(id);if err!=nil{http.Error(w,"failed to load GIS record",http.StatusInternalServerError);return};w.Header().Set("Content-Type","application/json");if err:=json.NewEncoder(w).Encode(record);err!=nil{return}}
