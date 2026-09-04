package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const maxClickHouseBatchSize = maxSiLKBatchSize

type FTNClickHouseStore struct { endpoint string; database string; user string; password string; client *http.Client }

type clickHouseFlowRow struct {
	ObservedAt time.Time `json:"observed_at"`; FirstSeen *time.Time `json:"first_seen,omitempty"`; LastSeen *time.Time `json:"last_seen,omitempty"`
	TenantID string `json:"tenant_id"`; ExporterID string `json:"exporter_id"`; Version uint16 `json:"version"`
	SourceIP string `json:"source_ip"`; DestinationIP string `json:"destination_ip"`; SourcePort uint16 `json:"source_port"`; DestinationPort uint16 `json:"destination_port"`; Protocol uint8 `json:"protocol"`
	Bytes uint64 `json:"bytes"`; Packets uint64 `json:"packets"`; SamplingRate uint32 `json:"sampling_rate"`; Fingerprint string `json:"fingerprint"`
	ServiceID string `json:"service_id"`; SubscriberID string `json:"subscriber_id"`; MainServerID string `json:"main_server_id"`
}

func NewFTNClickHouseStoreFromEnv() (*FTNClickHouseStore, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(os.Getenv("FTN_CLICKHOUSE_URL")), "/")
	if endpoint == "" { return nil, nil }
	u, err := url.Parse(endpoint); if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" { return nil, errors.New("invalid_clickhouse_url") }
	database := strings.TrimSpace(os.Getenv("FTN_CLICKHOUSE_DATABASE")); if database == "" { database = "ftn_analytics" }
	return &FTNClickHouseStore{endpoint:endpoint,database:database,user:os.Getenv("FTN_CLICKHOUSE_USER"),password:os.Getenv("FTN_CLICKHOUSE_PASSWORD"),client:&http.Client{Timeout:5*time.Second}},nil
}

func (s *FTNClickHouseStore) InsertBatchObservations(ctx context.Context, tenantID string, observations []FlowObservation, serviceID, subscriberID, mainServerID string) error {
	if s == nil || len(observations) == 0 { return nil }
	if len(observations) > maxClickHouseBatchSize { return fmt.Errorf("flow_batch_too_large: %d > %d",len(observations),maxClickHouseBatchSize) }
	if strings.TrimSpace(tenantID)=="" { return errors.New("tenant_required") }
	var body bytes.Buffer; enc:=json.NewEncoder(&body)
	for _, o := range observations {
		if err:=o.Validate(); err!=nil{return err}
		r:=o.FlowRecord; row:=clickHouseFlowRow{ObservedAt:o.ObservedAt.UTC(),TenantID:tenantID,ExporterID:r.ExporterID,Version:r.Version,SourceIP:r.SourceIP,DestinationIP:r.DestinationIP,SourcePort:r.SourcePort,DestinationPort:r.DestinationPort,Protocol:r.Protocol,Bytes:r.Bytes,Packets:r.Packets,SamplingRate:r.SamplingRate,Fingerprint:fmt.Sprintf("%016x",FlowRecordFingerprint(r)),ServiceID:serviceID,SubscriberID:subscriberID,MainServerID:mainServerID}
		if !o.FirstSeen.IsZero(){v:=o.FirstSeen.UTC();row.FirstSeen=&v};if !o.LastSeen.IsZero(){v:=o.LastSeen.UTC();row.LastSeen=&v}
		if err:=enc.Encode(row);err!=nil{return err}
	}
	query:="INSERT INTO "+quoteClickHouseIdent(s.database)+".flow_records FORMAT JSONEachRow"
	req,err:=http.NewRequestWithContext(ctx,http.MethodPost,s.endpoint+"/?query="+url.QueryEscape(query),&body);if err!=nil{return err};req.Header.Set("Content-Type","application/x-ndjson");if s.user!=""{req.SetBasicAuth(s.user,s.password)}
	resp,err:=s.client.Do(req);if err!=nil{return err};defer resp.Body.Close();if resp.StatusCode/100!=2{return fmt.Errorf("clickhouse_insert_http_%s",strconv.Itoa(resp.StatusCode))};return nil
}

func (s *FTNClickHouseStore) InsertBatch(ctx context.Context, observedAt time.Time, tenantID string, records []FlowRecord, serviceID, subscriberID, mainServerID string) error {
	if len(records)==0{return nil};observations:=make([]FlowObservation,0,len(records));for _,r:=range records{o,err:=normalizeFlowObservation(r,observedAt,time.Time{},time.Time{});if err!=nil{return err};observations=append(observations,o)};return s.InsertBatchObservations(ctx,tenantID,observations,serviceID,subscriberID,mainServerID)
}

func (s *FTNClickHouseStore) Ping(ctx context.Context) error { if s==nil{return errors.New("clickhouse_store_required")};req,err:=http.NewRequestWithContext(ctx,http.MethodPost,s.endpoint+"/?query=SELECT%201",nil);if err!=nil{return err};if s.user!=""{req.SetBasicAuth(s.user,s.password)};resp,err:=s.client.Do(req);if err!=nil{return err};defer resp.Body.Close();if resp.StatusCode/100!=2{return fmt.Errorf("clickhouse_ping_http_%s",strconv.Itoa(resp.StatusCode))};return nil }
func quoteClickHouseIdent(v string) string{return "`"+strings.ReplaceAll(v,"`","``")+"`"}
