package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PowerDNSHTTPClient struct {
	Endpoint string
	ServerID string
	APIKey string
	Client *http.Client
}

type powerDNSRRSet struct {
	Name string `json:"name"`
	Type string `json:"type"`
	TTL uint32 `json:"ttl"`
	Changetype string `json:"changetype"`
	Records []struct { Content string `json:"content"`; Disabled bool `json:"disabled"` } `json:"records"`
}

type powerDNSZonePatch struct { RRsets []powerDNSRRSet `json:"rrsets"` }

func NewPowerDNSHTTPClient(endpoint, serverID, apiKey string) *PowerDNSHTTPClient {
	return &PowerDNSHTTPClient{Endpoint: NormalizeProviderEndpoint(endpoint), ServerID: strings.TrimSpace(serverID), APIKey: apiKey, Client: &http.Client{Timeout: 15 * time.Second}}
}

// ApplyZone sends an explicit RRset patch to PowerDNS Authoritative's HTTP API.
// Callers should provision credentials through a secret manager rather than source control.
func (c *PowerDNSHTTPClient) ApplyZone(ctx context.Context, zone Zone) error {
	if c.Endpoint == "" || c.ServerID == "" { return fmt.Errorf("PowerDNS endpoint and server ID are required") }
	if strings.TrimSpace(c.APIKey) == "" { return fmt.Errorf("PowerDNS API key is required") }
	if zone.Name == "" { return fmt.Errorf("zone name is required") }
	patch := powerDNSZonePatch{RRsets: make([]powerDNSRRSet, 0, len(zone.Records))}
	for _, r := range zone.Records {
		name := strings.TrimSuffix(strings.TrimSpace(r.Name), ".") + "."
		rr := powerDNSRRSet{Name: name, Type: strings.ToUpper(strings.TrimSpace(r.Type)), TTL: r.TTL, Changetype: "REPLACE"}
		rr.Records = append(rr.Records, struct { Content string `json:"content"`; Disabled bool `json:"disabled"` }{Content: r.Value, Disabled: false})
		patch.RRsets = append(patch.RRsets, rr)
	}
	body, err := json.Marshal(patch); if err != nil { return err }
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.Endpoint, c.ServerID, strings.TrimSuffix(zone.Name, ".")+".")
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body)); if err != nil { return err }
	req.Header.Set("X-API-Key", c.APIKey); req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req); if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("PowerDNS API returned HTTP %d", resp.StatusCode) }
	return nil
}
