package adapters

import("context";"fmt";"net";"strings";"time";ftndns "github.com/beparykamrul-dev/FTN-AI/services/ftn-dns")
type GoDNSAdapter struct{Endpoint string;Timeout time.Duration}
func(g GoDNSAdapter)Name()string{return "godns"}
func(g GoDNSAdapter)Health(ctx context.Context)(ftndns.Health,error){if ctx==nil{return ftndns.Health{},fmt.Errorf("context is required")};endpoint:=strings.TrimSpace(g.Endpoint);if endpoint==""{return ftndns.Health{},fmt.Errorf("godns endpoint is required")};timeout:=g.Timeout;if timeout<=0{timeout=2*time.Second};start:=time.Now();dialer:=net.Dialer{Timeout:timeout};conn,err:=dialer.DialContext(ctx,"tcp",endpoint);latency:=time.Since(start);if err!=nil{return ftndns.Health{LatencyMS:latency.Milliseconds()},fmt.Errorf("godns endpoint unreachable: %w",err)};_ = conn.Close();return ftndns.Health{Reachable:true,LatencyMS:latency.Milliseconds(),LossRatio:0},nil}
func(g GoDNSAdapter)Query(ctx context.Context,name,recordType string)(ftndns.Response,error){if ctx==nil{return ftndns.Response{},fmt.Errorf("context is required")};if strings.TrimSpace(name)==""||strings.TrimSpace(recordType)==""{return ftndns.Response{},fmt.Errorf("DNS name and record type are required")};select{case<-ctx.Done():return ftndns.Response{},ctx.Err();default:};return ftndns.Response{},fmt.Errorf("godns query adapter requires the configured DNS transport")}
