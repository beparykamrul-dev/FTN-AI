package main

import (
	"errors"
	"math"
	"sort"
	"strings"
)

type FTNDNSDDNSRecord struct { Name string `json:"name"`; Type string `json:"type"`; Value string `json:"value"`; TTL uint32 `json:"ttl"`; Authenticated bool `json:"authenticated"` }
type FTNAnycastDNSNode struct { NodeID string `json:"node_id"`; Address string `json:"address"`; Healthy bool `json:"healthy"`; BFDUp bool `json:"bfd_up"`; RPKIValid bool `json:"rpki_valid"`; LatencyMS float64 `json:"latency_ms"`; Capacity float64 `json:"capacity"` }
func NormalizeFTNDNSDDNSRecord(r FTNDNSDDNSRecord)(FTNDNSDDNSRecord,error){r.Name=strings.ToLower(strings.TrimSpace(strings.TrimSuffix(r.Name,".")));r.Type=strings.ToUpper(strings.TrimSpace(r.Type));r.Value=strings.TrimSpace(r.Value);if r.Name==""||r.Type==""||r.Value==""{return FTNDNSDDNSRecord{},errors.New("ftn_dns_record_required")};if r.TTL==0{r.TTL=60};if r.TTL>86400{return FTNDNSDDNSRecord{},errors.New("ftn_dns_ttl_out_of_range")};if !r.Authenticated{return FTNDNSDDNSRecord{},errors.New("ftn_ddns_authentication_required")};return r,nil}
func SelectFTNAnycastDNSNodes(nodes []FTNAnycastDNSNode,max int)[]FTNAnycastDNSNode{if max<1{max=1};out:=make([]FTNAnycastDNSNode,0,len(nodes));for _,n:=range nodes{if n.Healthy&&n.BFDUp&&n.RPKIValid&&strings.TrimSpace(n.NodeID)!=""&&strings.TrimSpace(n.Address)!=""&&!math.IsNaN(n.LatencyMS)&&!math.IsInf(n.LatencyMS,0)&&!math.IsNaN(n.Capacity)&&!math.IsInf(n.Capacity,0)&&n.LatencyMS>=0&&n.Capacity>=0{out=append(out,n)}};sort.SliceStable(out,func(i,j int)bool{if out[i].LatencyMS!=out[j].LatencyMS{return out[i].LatencyMS<out[j].LatencyMS};if out[i].Capacity!=out[j].Capacity{return out[i].Capacity>out[j].Capacity};return out[i].NodeID<out[j].NodeID});if len(out)>max{out=out[:max]};return out}
