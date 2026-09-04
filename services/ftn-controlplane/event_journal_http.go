package controlplane

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

type HTTPEventJournal struct { BaseURL string; Token string; Client *http.Client }
type eventOffsetResponse struct { ConsumerID string `json:"consumer_id"`; TenantID string `json:"tenant_id"`; Sequence uint64 `json:"sequence"` }

func (j *HTTPEventJournal) client() *http.Client { if j != nil && j.Client != nil { return j.Client }; return &http.Client{Timeout:10*time.Second} }
func (j *HTTPEventJournal) endpoint(path string) string { if j == nil { return "" }; return strings.TrimRight(j.BaseURL,"/")+path }
func (j *HTTPEventJournal) do(method,path string,body any,out any) error { if j==nil||strings.TrimSpace(j.BaseURL)==""{return errors.New("event journal base URL is required")}; var rd io.Reader; if body!=nil { b,e:=json.Marshal(body);if e!=nil{return e};rd=bytes.NewReader(b) }; req,e:=http.NewRequest(method,j.endpoint(path),rd);if e!=nil{return e};req.Header.Set("Content-Type","application/json");if strings.TrimSpace(j.Token)!=""{req.Header.Set("Authorization","Bearer "+j.Token)};resp,e:=j.client().Do(req);if e!=nil{return e};defer resp.Body.Close();if resp.StatusCode<200||resp.StatusCode>=300{b,_:=io.ReadAll(io.LimitReader(resp.Body,4096));return fmt.Errorf("event api: %s: %s",resp.Status,strings.TrimSpace(string(b)))};if out==nil{return nil};return json.NewDecoder(io.LimitReader(resp.Body,2<<20)).Decode(out) }
func (j *HTTPEventJournal) Append(event JournalEvent)(JournalEvent,error){if strings.TrimSpace(event.TenantID)==""||strings.TrimSpace(event.Type)==""{return JournalEvent{},errors.New("event requires tenant_id and type")};if !json.Valid(event.Payload)&&len(event.Payload)>0{return JournalEvent{},errors.New("event payload must be valid JSON")};var out JournalEvent;e:=j.do(http.MethodPost,"/api/v1/events/append",event,&out);if e!=nil{return JournalEvent{},e};if out.TenantID!=event.TenantID{return JournalEvent{},errors.New("event journal returned mismatched tenant")};return out,nil}
func (j *HTTPEventJournal) ReadAfter(tenantID string,sequence uint64,limit int)([]JournalEvent,error){if strings.TrimSpace(tenantID)==""{return nil,errors.New("tenant_id is required")};if limit<=0{return []JournalEvent{},nil};if limit>500{limit=500};var out struct{Events []JournalEvent `json:"events"`};path:=fmt.Sprintf("/api/v1/events?after=%d&limit=%d&tenant_id=%s",sequence,limit,url.QueryEscape(tenantID));e:=j.do(http.MethodGet,path,nil,&out);if e!=nil{return nil,e};for _,event:=range out.Events{if event.TenantID!=tenantID{return nil,errors.New("event journal returned cross-tenant event")}};return out.Events,nil}
func (j *HTTPEventJournal) CommitOffset(consumerID,tenantID string,sequence uint64)error{if strings.TrimSpace(consumerID)==""||strings.TrimSpace(tenantID)==""{return errors.New("consumer_id and tenant_id are required")};return j.do(http.MethodPost,"/api/v1/events/offset/commit",map[string]any{"consumer_id":consumerID,"tenant_id":tenantID,"sequence":sequence},nil)}
func (j *HTTPEventJournal) Offset(consumerID,tenantID string)uint64{if strings.TrimSpace(consumerID)==""||strings.TrimSpace(tenantID)==""{return 0};var out eventOffsetResponse;path:=fmt.Sprintf("/api/v1/events/offset?consumer_id=%s&tenant_id=%s",url.QueryEscape(consumerID),url.QueryEscape(tenantID));if e:=j.do(http.MethodGet,path,nil,&out);e!=nil||out.TenantID!=tenantID||out.ConsumerID!=consumerID{return 0};return out.Sequence}
var _ EventJournal = (*HTTPEventJournal)(nil)
