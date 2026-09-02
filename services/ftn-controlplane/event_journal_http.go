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

func (j *HTTPEventJournal) client() *http.Client { if j.Client != nil { return j.Client }; return &http.Client{Timeout:10*time.Second} }
func (j *HTTPEventJournal) endpoint(path string) string { return strings.TrimRight(j.BaseURL,"/")+path }
func (j *HTTPEventJournal) do(method,path string,body any,out any) error { var rd io.Reader; if body!=nil { b,e:=json.Marshal(body);if e!=nil{return e};rd=bytes.NewReader(b) }; req,e:=http.NewRequest(method,j.endpoint(path),rd);if e!=nil{return e};req.Header.Set("Content-Type","application/json");if strings.TrimSpace(j.Token)!=""{req.Header.Set("Authorization","Bearer "+j.Token)};resp,e:=j.client().Do(req);if e!=nil{return e};defer resp.Body.Close();if resp.StatusCode<200||resp.StatusCode>=300{b,_:=io.ReadAll(io.LimitReader(resp.Body,4096));return fmt.Errorf("event api: %s: %s",resp.Status,strings.TrimSpace(string(b)))};if out==nil{return nil};return json.NewDecoder(resp.Body).Decode(out) }
func (j *HTTPEventJournal) Append(event JournalEvent)(JournalEvent,error){if strings.TrimSpace(event.TenantID)==""||strings.TrimSpace(event.Type)==""{return JournalEvent{},errors.New("event requires tenant_id and type")};var out JournalEvent;e:=j.do(http.MethodPost,"/api/v1/events/append",event,&out);return out,e}
func (j *HTTPEventJournal) ReadAfter(tenantID string,sequence uint64,limit int)([]JournalEvent,error){if limit<=0{return []JournalEvent{},nil};if limit>500{limit=500};var out struct{Events []JournalEvent `json:"events"`};path:=fmt.Sprintf("/api/v1/events?after=%d&limit=%d&tenant_id=%s",sequence,limit,url.QueryEscape(tenantID));e:=j.do(http.MethodGet,path,nil,&out);return out.Events,e}
func (j *HTTPEventJournal) CommitOffset(consumerID,tenantID string,sequence uint64)error{return j.do(http.MethodPost,"/api/v1/events/offset/commit",map[string]any{"consumer_id":consumerID,"tenant_id":tenantID,"sequence":sequence},nil)}
func (j *HTTPEventJournal) Offset(consumerID,tenantID string)uint64{var out eventOffsetResponse;path:=fmt.Sprintf("/api/v1/events/offset?consumer_id=%s",url.QueryEscape(consumerID));if e:=j.do(http.MethodGet,path,nil,&out);e!=nil{return 0};return out.Sequence}
var _ EventJournal = (*HTTPEventJournal)(nil)
