package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RouterOSCredentials contains references to credentials held by the FTN
// secret store. Raw passwords must never be persisted in device models.
type RouterOSCredentials struct {
	Username    string
	Password    string
	ServerName  string
	InsecureTLS bool
}

type RouterOSAdapter struct {
	client  *http.Client
	baseURL string
	creds   RouterOSCredentials
}

func NewRouterOSAdapter(baseURL string, creds RouterOSCredentials, timeout time.Duration) (*RouterOSAdapter, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, errors.New("routeros base url is required")
	}
	if creds.Username == "" {
		return nil, errors.New("routeros username is required")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &RouterOSAdapter{
		baseURL: baseURL,
		creds:   creds,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					ServerName: creds.ServerName,
					InsecureSkipVerify: creds.InsecureTLS, //nolint:gosec // explicitly controlled by FTN device policy
				},
			},
		},
	}, nil
}

func (a *RouterOSAdapter) Protocol() string { return "routeros-api" }

func (a *RouterOSAdapter) doJSON(ctx context.Context, path string) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+"/rest/"+strings.TrimLeft(path, "/"), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(a.creds.Username, a.creds.Password)
	req.Header.Set("Accept", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("routeros REST %s: status=%d body=%q", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var rows []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("routeros REST %s: decode: %w", path, err)
	}
	return rows, nil
}

func (a *RouterOSAdapter) Capabilities(ctx context.Context, device NetworkDevice) ([]string, error) {
	if err := a.probe(ctx); err != nil {
		return nil, err
	}
	return []string{"health.read", "interface.read", "routing.read"}, nil
}

func (a *RouterOSAdapter) probe(ctx context.Context) error {
	_, err := a.doJSON(ctx, "system/resource")
	return err
}

func (a *RouterOSAdapter) CollectInterfaceState(ctx context.Context, device NetworkDevice) ([]InterfaceState, error) {
	rows, err := a.doJSON(ctx, "interface")
	if err != nil {
		return nil, err
	}
	out := make([]InterfaceState, 0, len(rows))
	for _, row := range rows {
		name := stringValue(row["name"])
		if name == "" {
			continue
		}
		out = append(out, InterfaceState{
			DeviceID: device.ID,
			Name: name,
			AdminUp: !boolValue(row["disabled"]),
			OperUp: boolValue(row["running"]),
			RxBps: uint64Value(row["rx-byte"]),
			TxBps: uint64Value(row["tx-byte"]),
		})
	}
	return out, nil
}

func (a *RouterOSAdapter) CollectRoutingState(ctx context.Context, device NetworkDevice) ([]RoutingState, error) {
	rows, err := a.doJSON(ctx, "ip/route")
	if err != nil {
		return nil, err
	}
	out := make([]RoutingState, 0, len(rows))
	for _, row := range rows {
		prefix := stringValue(row["dst-address"])
		if prefix == "" {
			continue
		}
		out = append(out, RoutingState{
			DeviceID: device.ID,
			Protocol: stringValue(row["protocol"]),
			VRF: stringValue(row["routing-table"]),
			Prefix: prefix,
			NextHop: stringValue(row["gateway"]),
			Metric: uint32Value(row["distance"]),
			Active: boolValue(row["active"]),
		})
	}
	return out, nil
}

func stringValue(v any) string {
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func boolValue(v any) bool {
	b, _ := v.(bool)
	return b
}

func uint64Value(v any) uint64 {
	switch n := v.(type) {
	case float64:
		if n > 0 { return uint64(n) }
	case json.Number:
		if value, err := n.Int64(); err == nil && value > 0 { return uint64(value) }
	case string:
		var value uint64
		if _, err := fmt.Sscan(n, &value); err == nil { return value }
	}
	return 0
}

func uint32Value(v any) uint32 {
	value := uint64Value(v)
	if value > ^uint32(0) { return ^uint32(0) }
	return uint32(value)
}
