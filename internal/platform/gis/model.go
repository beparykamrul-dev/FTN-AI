package gis

import "time"

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IPAsset struct {
	ID         string     `json:"id"`
	IP         string     `json:"ip"`
	CIDR       string     `json:"cidr,omitempty"`
	MAC        string     `json:"mac,omitempty"`
	Hostname   string     `json:"hostname,omitempty"`
	ServerID   string     `json:"server_id,omitempty"`
	Location   *Coordinate `json:"location,omitempty"`
	Provider   string     `json:"provider,omitempty"`
	ASN        uint32     `json:"asn,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type MapNode struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Kind     string      `json:"kind"`
	Location Coordinate  `json:"location"`
	IPAssets []string    `json:"ip_assets,omitempty"`
}

type MapEdge struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind"`
}

type MapSnapshot struct {
	Nodes []MapNode `json:"nodes"`
	Edges []MapEdge `json:"edges"`
}
