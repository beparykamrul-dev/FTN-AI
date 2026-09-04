package controlplane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type remoteEventOffset struct { ConsumerID string `json:"consumer_id"`; TenantID string `json:"tenant_id"`; Sequence uint64 `json:"sequence"` }

// RemoteEventJournal connects ftn-controlplane to the durable event journal
// exposed by the FTN control-plane without coupling this service to PostgreSQL.
type RemoteEventJournal struct { BaseURL string; Token string; Client *http.Client }

func validRemoteJournalURL(raw string) bool {
	u, err := url.Parse(strings.TrimRight(strings.TrimSpace(raw), "/"))
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil
}

func NewRemoteEventJournal(baseURL, token string) *RemoteEventJournal { return &RemoteEventJournal{BaseURL:strings.TrimRight(strings.TrimSpace(baseURL),"/"),Token:token,Client:&http.Client{Timeout:5*time.Second}} }
func (j *RemoteEventJournal) client()*http.Client{if j.Client!=nil{return j.Client};return &http.Client{Timeout:5*time.Second}}
func(j *RemoteEventJournal)request(ctx context.Context,method,path string,body []byte)([]byte,int,error){if j==nil||!validRemoteJournalURL(j.BaseURL){return nil,0,fmt.Errorf("remote event journal base URL is invalid")};if ctx==nil{ctx=context.Background()};reqURL,err:=url.Parse(j.BaseURL+path);if err!=nil{return nil,0,err};if reqURL.User!=nil||reqURL.Scheme!="http"&&reqURL.Scheme!="https"||reqURL.Host==""{return nil,0,fmt.Errorf("remote event journal request URL is invalid")};req,err:=http.NewRequestWithContext(ctx,method,reqURL.String(),bytes.NewReader(body));if err!=nil{return nil,0,err};if j.Token!=""{req.Header.Set("Authorization","Bearer "+j.Token)};req.Header.Set("Content-Type","application/json");resp,err:=j.client().Do(req);if err!=nil{return nil,0,err};defer resp.Body.Close();data,err:=io.ReadAll(io.LimitReader(resp.Body,2<<20));if err!=nil{return nil,resp.StatusCode,err};if resp.StatusCode<200||resp.StatusCode>=300{return data,resp.StatusCode,fmt.Errorf("control-plane returned HTTP %d",resp.StatusCode)};return data,resp.StatusCode,nil}
func(j *RemoteEventJournal)Append(event JournalEvent)(JournalEvent,error){if strings.TrimSpace(event.TenantID)==""||strings.TrimSpace(event.Type)==""{return JournalEvent{},fmt.Errorf("event requires tenant_id and type")};payload:=json.RawMessage(event.Payload);if payload==nil{payload=json.RawMessage(`{}`)};if !json.Valid(payload){return JournalEvent{},fmt.Errorf("event payload must be valid JSON")};body,err:=json.Marshal(struct{TenantID string `json:"tenant_id"`;Type string `json:"type"`;CorrelationID string `json:"correlation_id"`;CausationID string `json:"causation_id,omitempty"`;AggregateID string `json:"aggregate_id,omitempty"`;Payload json.RawMessage `json:"payload"`}{event.TenantID,event.Type,event.CorrelationID,event.CausationID,event.AggregateID,payload});if err!=nil{return JournalEvent{},err};data,_,err:=j.request(context.Background(),http.MethodPost,"/api/v1/events/append",body);if err!=nil{return JournalEvent{},err};var e struct{ID string `json:"id"`;TenantID string `json:"tenant_id"`;Type string `json:"type"`;Sequence int64 `json:"sequence"`;CorrelationID string `json:"correlation_id"`;CausationID string `json:"causation_id"`;AggregateID string `json:"aggregate_id"`;Payload json.RawMessage `json:"payload"`;CreatedAt string `json:"created_at"`};if err=json.Unmarshal(data,&e);err!=nil{return JournalEvent{},err};if e.Sequence<0{return JournalEvent{},fmt.Errorf("control-plane returned negative event sequence")};created,err:=time.Parse(time.RFC3339Nano,e.CreatedAt);if err!=nil{return JournalEvent{},fmt.Errorf("invalid event created_at: %w",err)};if e.TenantID!=event.TenantID{return JournalEvent{},fmt.Errorf("control-plane returned mismatched event tenant")};return JournalEvent{ID:e.ID,TenantID:e.TenantID,Type:e.Type,Sequence:uint64(e.Sequence),CorrelationID:e.CorrelationID,CausationID:e.CausationID,AggregateID:e.AggregateID,Payload:append([]byte(nil),e.Payload...),CreatedAt:created},nil}
func(j *RemoteEventJournal)ReadAfter(tenantID string,sequence uint64,limit int)([]JournalEvent,error){if strings.TrimSpace(tenantID)==""{return nil,fmt.Errorf("tenant_id is required")};path:="/api/v1/events?after="+strconv.FormatUint(sequence,10)+"&tenant_id="+url.QueryEscape(tenantID);if limit>0{if limit>500{limit=500};path+="&limit="+strconv.Itoa(limit)};data,_,err:=j.request(context.Background(),http.MethodGet,path,nil);if err!=nil{return nil,err};var response struct{Events []struct{ID string `json:"id"`;TenantID string `json:"tenant_id"`;Type string `json:"type"`;Sequence int64 `json:"sequence"`;CorrelationID string `json:"correlation_id"`;CausationID string `json:"causation_id"`;AggregateID string `json:"aggregate_id"`;Payload json.RawMessage `json:"payload"`;CreatedAt string `json:"created_at"`} `json:"events"`};if err=json.Unmarshal(data,&response);err!=nil{return nil,err};out:=make([]JournalEvent,0,len(response.Events));for _,e:=range response.Events{if e.TenantID!=tenantID{return nil,fmt.Errorf("control-plane returned cross-tenant event")};if e.Sequence<0{return nil,fmt.Errorf("control-plane returned negative event sequence")};created,err:=time.Parse(time.RFC3339Nano,e.CreatedAt);if err!=nil{return nil,fmt.Errorf("invalid event created_at: %w",err)};out=append(out,JournalEvent{ID:e.ID,TenantID:e.TenantID,Type:e.Type,Sequence:uint64(e.Sequence),CorrelationID:e.CorrelationID,CausationID:e.CausationID,AggregateID:e.AggregateID,Payload:append([]byte(nil),e.Payload...),CreatedAt:created})};return out,nil}
func(j *RemoteEventJournal)CommitOffset(consumerID,tenantID string,sequence uint64)error{if strings.TrimSpace(consumerID)==""||strings.TrimSpace(tenantID)==""{return fmt.Errorf("consumer_id and tenant_id are required")};body,err:=json.Marshal(remoteEventOffset{ConsumerID:consumerID,TenantID:tenantID,Sequence:sequence});if err!=nil{return err};_,_,err=j.request(context.Background(),http.MethodPost,"/api/v1/events/offset/commit",body);return err}
func(j *RemoteEventJournal)Offset(consumerID,tenantID string)uint64{if strings.TrimSpace(consumerID)==""||strings.TrimSpace(tenantID)==""{return 0};data,_,err:=j.request(context.Background(),http.MethodGet,"/api/v1/events/offset?consumer_id="+url.QueryEscape(consumerID)+"&tenant_id="+url.QueryEscape(tenantID),nil);if err!=nil{return 0};var o remoteEventOffset;if json.Unmarshal(data,&o)!=nil||o.TenantID!=tenantID||o.ConsumerID!=consumerID{return 0};return o.Sequence}
