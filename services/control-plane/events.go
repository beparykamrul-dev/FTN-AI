package main

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "strings"

    "github.com/jackc/pgx/v5"
)

type DurableEvent struct { ID string `json:"id,omitempty"`; TenantID string `json:"tenant_id,omitempty"`; Type string `json:"type"`; Sequence int64 `json:"sequence,omitempty"`; CorrelationID string `json:"correlation_id"`; CausationID string `json:"causation_id,omitempty"`; AggregateID string `json:"aggregate_id,omitempty"`; Payload json.RawMessage `json:"payload"`; CreatedAt string `json:"created_at,omitempty"` }
type EventOffset struct { ConsumerID string `json:"consumer_id"`; TenantID string `json:"tenant_id"`; Sequence int64 `json:"sequence"` }

// appendEventTx is the transaction-safe primitive used by state-changing operations.
func appendEventTx(tx pgx.Tx, ctx context.Context, tenantID, eventType, correlationID, causationID, aggregateID string, payload json.RawMessage) (DurableEvent, error) {
    if strings.TrimSpace(eventType)=="" { return DurableEvent{}, errors.New("event_type_required") }
    if payload==nil { payload=json.RawMessage(`{}`) }
    if !json.Valid(payload) { return DurableEvent{}, errors.New("invalid_event_payload") }
    if _,err:=tx.Exec(ctx,`select pg_advisory_xact_lock(hashtext($1::text))`,tenantID);err!=nil{return DurableEvent{},err}
    var e DurableEvent
    err:=tx.QueryRow(ctx,`insert into event_journal(tenant_id,event_type,sequence,correlation_id,causation_id,aggregate_id,payload) values(nullif($1,'')::uuid,$2,coalesce((select max(sequence)+1 from event_journal where tenant_id=nullif($1,'')::uuid),1),$3,$4,$5,$6::jsonb) returning id::text,coalesce(tenant_id::text,''),event_type,sequence,correlation_id,causation_id,aggregate_id,payload::text,created_at::text`,tenantID,eventType,correlationID,causationID,aggregateID,string(payload)).Scan(&e.ID,&e.TenantID,&e.Type,&e.Sequence,&e.CorrelationID,&e.CausationID,&e.AggregateID,&e.Payload,&e.CreatedAt)
    return e,err
}

func (a *App) appendEvent(w http.ResponseWriter,r *http.Request){
    if !method(w,r,http.MethodPost)||!requirePermission(a,"event.append",w,r){return};var e DurableEvent
    if err:=json.NewDecoder(r.Body).Decode(&e);err!=nil||strings.TrimSpace(e.Type)==""{jsonResponse(w,400,map[string]string{"error":"invalid_event"});return}
    rc:=requestInfo(r);if e.CorrelationID==""{e.CorrelationID=rc.CorrelationID};if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return}
    tx,err:=a.db.Begin(r.Context());if err!=nil{jsonResponse(w,500,map[string]string{"error":"event_begin_failed"});return};defer tx.Rollback(r.Context())
    out,err:=appendEventTx(tx,r.Context(),rc.TenantID,e.Type,e.CorrelationID,e.CausationID,e.AggregateID,e.Payload);if err!=nil{jsonResponse(w,500,map[string]string{"error":"event_append_failed"});return}
    if err=tx.Commit(r.Context());err!=nil{jsonResponse(w,500,map[string]string{"error":"event_commit_failed"});return};a.audit(r,"event.append",out.Type,"accepted",out);jsonResponse(w,202,out)
}

func (a *App) readEvents(w http.ResponseWriter,r *http.Request){
    if !method(w,r,http.MethodGet)||!requirePermission(a,"event.read",w,r){return};if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return};rc:=requestInfo(r);after,_:=strconv.ParseInt(r.URL.Query().Get("after"),10,64);limit,_:=strconv.Atoi(r.URL.Query().Get("limit"));if limit<=0||limit>500{limit=100}
    rows,err:=a.db.Query(r.Context(),`select id::text,coalesce(tenant_id::text,''),event_type,sequence,correlation_id,causation_id,aggregate_id,payload::text,created_at::text from event_journal where tenant_id=nullif($1,'')::uuid and sequence>$2 order by sequence asc limit $3`,rc.TenantID,after,limit);if err!=nil{jsonResponse(w,500,map[string]string{"error":"event_read_failed"});return};defer rows.Close();out:=make([]DurableEvent,0,limit)
    for rows.Next(){var e DurableEvent;if err:=rows.Scan(&e.ID,&e.TenantID,&e.Type,&e.Sequence,&e.CorrelationID,&e.CausationID,&e.AggregateID,&e.Payload,&e.CreatedAt);err!=nil{jsonResponse(w,500,map[string]string{"error":"event_decode_failed"});return};out=append(out,e)};if err:=rows.Err();err!=nil{jsonResponse(w,500,map[string]string{"error":"event_read_failed"});return};jsonResponse(w,200,map[string]any{"events":out})
}

func (a *App) commitEventOffset(w http.ResponseWriter,r *http.Request){
    if !method(w,r,http.MethodPost)||!requirePermission(a,"event.commit",w,r){return};var o EventOffset;if err:=json.NewDecoder(r.Body).Decode(&o);err!=nil||strings.TrimSpace(o.ConsumerID)==""||o.Sequence<0{jsonResponse(w,400,map[string]string{"error":"invalid_offset"});return};rc:=requestInfo(r);if o.TenantID==""{o.TenantID=rc.TenantID};if o.TenantID==""{jsonResponse(w,400,map[string]string{"error":"tenant_required"});return};if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return}
    _,err:=a.db.Exec(r.Context(),`insert into event_consumer_offsets(consumer_id,tenant_id,sequence) values($1,$2,$3) on conflict(consumer_id,tenant_id) do update set sequence=greatest(event_consumer_offsets.sequence,excluded.sequence),updated_at=now()`,o.ConsumerID,o.TenantID,o.Sequence);if err!=nil{jsonResponse(w,500,map[string]string{"error":"offset_commit_failed"});return};a.audit(r,"event.offset.commit",o.ConsumerID,"accepted",o);jsonResponse(w,202,o)
}

func (a *App) eventOffset(w http.ResponseWriter,r *http.Request){
    if !method(w,r,http.MethodGet)||!requirePermission(a,"event.read",w,r){return};if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return};rc:=requestInfo(r);consumer:=strings.TrimSpace(r.URL.Query().Get("consumer_id"));if consumer==""||rc.TenantID==""{jsonResponse(w,400,map[string]string{"error":"consumer_id_required"});return};var seq int64;err:=a.db.QueryRow(r.Context(),`select sequence from event_consumer_offsets where consumer_id=$1 and tenant_id=$2`,consumer,rc.TenantID).Scan(&seq);if err!=nil{if errors.Is(err,pgx.ErrNoRows){seq=0}else{jsonResponse(w,500,map[string]string{"error":"offset_read_failed"});return}};jsonResponse(w,200,EventOffset{ConsumerID:consumer,TenantID:rc.TenantID,Sequence:seq})
}
