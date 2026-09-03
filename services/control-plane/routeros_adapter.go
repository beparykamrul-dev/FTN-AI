package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RouterOSCredentialRef string

type RouterOSSecret struct { Username string; Password string }

type RouterOSAdapter struct {
	client *http.Client
	baseURL string
	credentialRef RouterOSCredentialRef
	serverName string
	resolveSecret func(context.Context, RouterOSCredentialRef) (RouterOSSecret, error)
}

func NewRouterOSAdapter(baseURL string, credentialRef RouterOSCredentialRef, serverName string, resolveSecret func(context.Context, RouterOSCredentialRef) (RouterOSSecret, error), timeout time.Duration) (*RouterOSAdapter, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" { return nil, errors.New("routeros base url is required") }
	if credentialRef == "" { return nil, errors.New("routeros credential reference is required") }
	if resolveSecret == nil { return nil, errors.New("routeros secret resolver is required") }
	if timeout <= 0 { timeout = 10 * time.Second }
	return &RouterOSAdapter{baseURL: baseURL, credentialRef: credentialRef, serverName: serverName, resolveSecret: resolveSecret, client: &http.Client{Timeout: timeout, Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, ServerName: serverName}}}}, nil
}
func (a *RouterOSAdapter) Protocol() string { return "routeros-api" }

func (a *RouterOSAdapter) doJSON(ctx context.Context, path string) ([]map[string]any, error) {
	secret, err := a.resolveSecret(ctx, a.credentialRef); if err != nil { return nil, fmt.Errorf("routeros secret lookup: %w", err) }
	if secret.Username == "" || secret.Password == "" { return nil, errors.New("routeros secret is incomplete") }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+"/rest/"+strings.TrimLeft(path, "/"), nil); if err != nil { return nil, err }
	req.SetBasicAuth(secret.Username, secret.Password); req.Header.Set("Accept", "application/json")
	resp, err := a.client.Do(req); if err != nil { return nil, err }; defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)); return nil, fmt.Errorf("routeros REST %s: status=%d body=%q", path, resp.StatusCode, strings.TrimSpace(string(body))) }
	var rows []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil { return nil, fmt.Errorf("routeros REST %s: decode: %w", path, err) }
	return rows, nil
}

func (a *RouterOSAdapter) Capabilities(ctx context.Context, device NetworkDevice) ([]string, error) {
	if err := validateRouterOSDevice(device); err != nil { return nil, err }
	if _, err := a.doJSON(ctx, "system/resource"); err != nil { return nil, err }
	return []string{"health.read", "interface.read", "routing.read", "qos.snapshot"}, nil
}
func (a *RouterOSAdapter) CollectInterfaceState(ctx context.Context, device NetworkDevice) ([]InterfaceState, error) {
	if err := validateRouterOSDevice(device); err != nil { return nil, err }; rows, err := a.doJSON(ctx, "interface"); if err != nil { return nil, err }
	out := make([]InterfaceState, 0, len(rows)); for _, row := range rows { name := stringValue(row["name"]); if name == "" { continue }; out = append(out, InterfaceState{DeviceID: device.ID, Name: name, AdminUp: !boolValue(row["disabled"]), OperUp: boolValue(row["running"]), RxBps: uint64Value(row["rx-byte"]), TxBps: uint64Value(row["tx-byte"]), RxErrors: uint64Value(row["rx-error"]), TxErrors: uint64Value(row["tx-error"])}) }; return out, nil
}
func (a *RouterOSAdapter) CollectRoutingState(ctx context.Context, device NetworkDevice) ([]RoutingState, error) {
	if err := validateRouterOSDevice(device); err != nil { return nil, err }; rows, err := a.doJSON(ctx, "ip/route"); if err != nil { return nil, err }
	out := make([]RoutingState, 0, len(rows)); for _, row := range rows { prefix := stringValue(row["dst-address"]); if prefix == "" { continue }; out = append(out, RoutingState{DeviceID: device.ID, Protocol: stringValue(row["protocol"]), VRF: stringValue(row["routing-table"]), Prefix: prefix, NextHop: stringValue(row["gateway"]), Metric: uint32Value(row["distance"]), Active: boolValue(row["active"])}) }; return out, nil
}

// ReadQoSSnapshot is strictly read-only. It requires an FTN ownership record,
// passes the common network execution guard, and reads only FTN-managed queue
// comments. No RouterOS mutation endpoint is invoked.
func (a *RouterOSAdapter) ReadQoSSnapshot(ctx context.Context, device NetworkDevice, owned bool) (RouterOSSnapshot, error) {
	if err := ValidateFTNOwnership(device, owned); err != nil { return RouterOSSnapshot{}, err }
	intent := NetworkExecutionIntent{Device: device, Action: NetworkRead}
	if err := ValidateNetworkExecutionIntent(intent); err != nil { return RouterOSSnapshot{}, err }
	decision := EvaluateNetworkExecution(intent); if !decision.Allowed { return RouterOSSnapshot{}, errors.New(decision.Reason) }
	rows, err := a.doJSON(ctx, "queue/simple"); if err != nil { return RouterOSSnapshot{}, err }
	rules, err := parseRouterOSQoSComments(rows); if err != nil { return RouterOSSnapshot{}, err }
	return NormalizeRouterOSSnapshot(RouterOSSnapshot{DeviceID: device.ID, Rules: rules, CapturedAt: time.Now().UTC()})
}

var routerOSQoSCommentRE = regexp.MustCompile(`^FTN-QOS service=([^ ]+) class=([^ ]+) path=([^ ]+) dscp=([0-9]+) priority=([0-9]+)$`)
func parseRouterOSQoSComments(rows []map[string]any) ([]RouterOSQoSState, error) {
	out := make([]RouterOSQoSState, 0, len(rows))
	for _, row := range rows { comment := strings.TrimSpace(stringValue(row["comment"])); if !strings.HasPrefix(comment, "FTN-QOS ") { continue }; m := routerOSQoSCommentRE.FindStringSubmatch(comment); if len(m) != 6 { return nil, errors.New("routeros_qos_malformed_ftn_comment") }; dscp, err := strconv.ParseUint(m[4], 10, 8); if err != nil || dscp > 63 { return nil, errors.New("routeros_qos_invalid_dscp") }; priority, err := strconv.ParseUint(m[5], 10, 8); if err != nil { return nil, errors.New("routeros_qos_invalid_priority") }; out = append(out, RouterOSQoSState{ServiceID:m[1], Class:TrafficClass(m[2]), PathID:m[3], DSCP:uint8(dscp), Priority:uint8(priority)}) }
	return normalizeRouterOSQoSState(out)
}

func validateRouterOSDevice(d NetworkDevice) error { if strings.TrimSpace(d.ID)=="" || strings.TrimSpace(d.Address)=="" { return errors.New("routeros device id and address are required") }; if strings.TrimSpace(d.Protocol)=="" { return errors.New("routeros protocol is required") }; if normalizeProtocol(d.Protocol)!="routeros-api" { return errors.New("routeros protocol mismatch") }; if !isFTNDeviceKind(d.Kind) { return errors.New("routeros target ownership is not verified") }; return nil }
func stringValue(v any) string { s, _ := v.(string); return strings.TrimSpace(s) }
func boolValue(v any) bool { b, _ := v.(bool); return b }
func uint64Value(v any) uint64 { switch n := v.(type) { case float64: if n>0 { return uint64(n) }; case json.Number: if value,err:=n.Int64(); err==nil && value>0 { return uint64(value) }; case string: var value uint64; if _,err:=fmt.Sscan(n,&value); err==nil { return value } }; return 0 }
func uint32Value(v any) uint32 { value:=uint64Value(v); if value>^uint32(0) { return ^uint32(0) }; return uint32(value) }
