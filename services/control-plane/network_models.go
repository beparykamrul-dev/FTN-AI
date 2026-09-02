package main

import "strings"

type NetworkDevice struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Vendor string `json:"vendor,omitempty"`
	Model string `json:"model,omitempty"`
	Address string `json:"address,omitempty"`
	Protocol string `json:"protocol"`
	Region string `json:"region,omitempty"`
	Role string `json:"role,omitempty"`
	Healthy bool `json:"healthy"`
}

type InterfaceState struct {
	DeviceID string `json:"device_id"`
	Name string `json:"name"`
	AdminUp bool `json:"admin_up"`
	OperUp bool `json:"oper_up"`
	RxBps uint64 `json:"rx_bps"`
	TxBps uint64 `json:"tx_bps"`
	RxErrors uint64 `json:"rx_errors"`
	TxErrors uint64 `json:"tx_errors"`
}

type RoutingState struct {
	DeviceID string `json:"device_id"`
	Protocol string `json:"protocol"`
	VRF string `json:"vrf,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	NextHop string `json:"next_hop,omitempty"`
	Metric uint32 `json:"metric,omitempty"`
	Active bool `json:"active"`
}

type FlowRecord struct {
	ExporterID string `json:"exporter_id"`
	Version uint16 `json:"version"`
	SourceIP string `json:"source_ip,omitempty"`
	DestinationIP string `json:"destination_ip,omitempty"`
	SourcePort uint16 `json:"source_port,omitempty"`
	DestinationPort uint16 `json:"destination_port,omitempty"`
	Protocol uint8 `json:"protocol,omitempty"`
	Bytes uint64 `json:"bytes"`
	Packets uint64 `json:"packets"`
	SamplingRate uint32 `json:"sampling_rate,omitempty"`
}

func normalizeProtocol(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
