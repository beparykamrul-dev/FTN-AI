package controlplane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RemoteEventJournal connects ftn-controlplane to the durable event journal
// exposed by the FTN control-plane. It keeps the service package independent
// of the PostgreSQL driver while retaining the same EventJournal contract.
type RemoteEventJournal struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewRemoteEventJournal(baseURL, token string) *RemoteEventJournal {
	return &RemoteEventJournal{BaseURL: strings.TrimRight(baseURL, "/"), Token: token, Client: &http.Client{Timeout: 5 * time.Second}}
}

func (j *RemoteEventJournal) client() *http.Client {
	if j.Client != nil { return j.Client }
	return &http.Client{Timeout: 5 * time.Second}
}

func (j *RemoteEventJournal) request(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
	if j.BaseURL == "" { return nil, 0, fmt.Errorf("remote event journal base URL is required") }
	req, err := http.NewRequestWithContext(ctx, method, j.BaseURL+path, bytes.NewReader(body))
	if err != nil { return nil, 0, err }
	if j.Token != "" { req.Header.Set("Authorization", "Bearer "+j.Token) }
	req.Header.Set("Content-Type", "application/json")
	resp, err := j.client().Do(req)
	if err != nil { return nil, 0, err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil { return nil, resp.StatusCode, err }
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return data, resp.StatusCode, fmt.Errorf("control-plane returned HTTP %d", resp.StatusCode) }
	return data, resp.StatusCode, nil
}

func (j *RemoteEventJournal) Append(event JournalEvent) (JournalEvent, error) {
	payload := json.RawMessage(event.Payload)
	if payload == nil { payload = json.RawMessage(`{}`) }
	body, _ := json.Marshal(struct {
		TenantID string `json:"tenant_id"`
		Type string `json:"type"`
		CorrelationID string `json:"correlation_id"`
		CausationID string `json:"causation_id,omitempty"`
		AggregateID string `json:"aggregate_id,omitempty"`
		Payload json.RawMessage `json:"payload"`
	}{event.TenantID,event.Type,event.CorrelationID,event.CausationID,event.AggregateID,payload})
	data, _, err := j.request(context.Background(), http.MethodPost, "/api/v1/events/append", body)
	if err != nil { return JournalEvent{}, err }
	var remote struct { ID string `json:"id"`; TenantID string `json:"tenant_id"`; Type string `json:"type"`; Sequence int64 `json:"sequence"`; CorrelationID string `json:"correlation_id"`; CausationID string `json:"causation_id"`; AggregateID string `json:"aggregate_id"`; Payload json.RawMessage `json:"payload"`; CreatedAt string `json:"created_at"` }
	if err := json.Unmarshal(data, &remote); err != nil { return JournalEvent{}, err }
	created, err := time.Parse(time.RFC3339Nano, remote.CreatedAt); if err != nil { created = time.Now().UTC() }
	return JournalEvent{ID:remote.ID,TenantID:remote.TenantID,Type:remote.Type,Sequence:uint64(remote.Sequence),CorrelationID:remote.CorrelationID,CausationID:remote.CausationID,AggregateID:remote.AggregateID,Payload:append([]byte(nil),remote.Payload...),CreatedAt:created}, nil
}

func (j *RemoteEventJournal) ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error) {
	path := "/api/v1/events?after=" + strconv.FormatUint(sequence,10)
	if limit > 0 { path += "&limit=" + strconv.Itoa(limit) }
	data, _, err := j.request(context.Background(), http.MethodGet, path, nil); if err != nil { return nil, err }
	var response struct { Events []struct { ID string `json:"id"`; TenantID string `json:"tenant_id"`; Type string `json:"type"`; Sequence int64 `json:"sequence"`; CorrelationID string `json:"correlation_id"`; CausationID string `json:"causation_id"`; AggregateID string `json:"aggregate_id"`; Payload json.RawMessage `json:"payload"`; CreatedAt string `json:"created_at"` } `json:"events"` }
	if err := json.Unmarshal(data,&response); err != nil { return nil, err }
	out:=make([]JournalEvent,0,len(response.Events));for _,e:=range response.Events{created,_:=time.Parse(time.RFC3339Nano,e.CreatedAt);out=append(out,JournalEvent{ID:e.ID,TenantID:e.TenantID,Type:e.Type,Sequence:uint64(e.Sequence),CorrelationID:e.CorrelationID,CausationID:e.CausationID,AggregateID:e.AggregateID,Payload:append([]byte(nil),e.Payload...),CreatedAt:created})};return out,nil
}

func (j *RemoteEventJournal) CommitOffset(consumerID, tenantID string, sequence uint64) error {
	body,_:=json.Marshal(EventOffset{ConsumerID:consumerID,TenantID:tenantID,Sequence:sequence});_,_,err:=j.request(context.Background(),http.MethodPost,"/api/v1/events/offset/commit",body);return err
}

func (j *RemoteEventJournal) Offset(consumerID, tenantID string) uint64 {
	data,_,err:=j.request(context.Background(),http.MethodGet,"/api/v1/events/offset?consumer_id="+consumerID,nil);if err!=nil{return 0};var o EventOffset;if json.Unmarshal(data,&o)!=nil{return 0};return o.Sequence
}
