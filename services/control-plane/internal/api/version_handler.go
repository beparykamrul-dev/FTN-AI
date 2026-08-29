package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type versionInfo struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func Version(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("FTN_CONTROL_PLANE_VERSION")
	if version == "" {
		version = "dev"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(versionInfo{Service: "control-plane", Version: version}))
}
