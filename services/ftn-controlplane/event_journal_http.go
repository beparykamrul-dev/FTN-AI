package controlplane

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// HTTPEventJournal connects ftn-controlplane to the durable event APIs exposed by FTN control-plane.
// It deliberately uses only the standard library so this service keeps its existing module boundary.
type HTTPEventJournal struct {
    BaseURL string
    Token   string
    Client  *http.Client
}

func (j *HTTPEventJournal) client() *http.Client { if j.Client != nil { return j.Client }; return &http.Client{Timeout: 10 * time.Second} }
func (j *HTTPEventJournal) url(path string) string { return strings.TrimRight(j.BaseURL, "/") + path }

func (j *HTTPEventJournal) do(method, path string, body any, out any) error {
    var reader io.Reader
    if body != nil { b, err := json.Marshal(body); if err != nil { return err }; reader = bytes.NewReader(b) }
    req, err := http.NewRequest(method, j.url(path), reader); if err != nil { return err }
    req.Header.Set("Content-Type", "application/json")
    if strings.TrimSpace(j.Token) != "" { req.Header.Set("Authorization", "Bearer "+j.Token) }
    resp, err := j.client().Do(req); if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { b,_:=io.ReadAll(io.LimitReader(resp.Body, 4096)); return fmt.Errorf("event api: %s: %s", resp.Status, strings.TrimSpace(string(b))) }
    if out == nil { return nil }
    return json.NewDecoder(resp.Body).Decode(out)
}

func (j *HTTPEventJournal) Append(event JournalEvent) (JournalEvent, error) {
    if strings.TrimSpace(event.TenantID)=="" || strings.TrimSpace(event.Type)=="" { return JournalEvent{}, errors.New("event requires tenant_id and type") }
    var out JournalEvent
    err:=j.do(http.MethodPost,"/api/v1/events/append",event,&out)
    return out,err
}

func (j *HTTPEventJournal) ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error) {
    if limit<=0 { return []JournalEvent{},nil }; if limit>500 { limit=500 }
    var out struct{ Events []JournalEvent `json:"events"` }
    path:=fmt.Sprintf("/api/v1/events?after=%d&limit=%d&tenant_id=%s",sequence,limit,httpURLQuery(tenantID))
    err:=j.do(http.MethodGet,path,nil,&out); return out.Events,err
}

func (j *HTTPEventJournal) CommitOffset(consumerID, tenantID string, sequence uint64) error {
    return j.do(http.MethodPost,"/api/v1/events/offset/commit",map[string]any{"consumer_id":consumerID,"tenant_id":tenantID,"sequence":sequence},nil)
}

func (j *HTTPEventJournal) Offset(consumerID, tenantID string) uint64 {
    var out EventOffset
    path:=fmt.Sprintf("/api/v1/events/offset?consumer_id=%s",httpURLQuery(consumerID))
    if err:=j.do(http.MethodGet,path,nil,&out); err!=nil { return 0 }; return out.Sequence
}

func httpURLQuery(s string) string { return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s,"%","%25")," ","%20"),"/","%2F") }
var _ = strconv.IntSize
