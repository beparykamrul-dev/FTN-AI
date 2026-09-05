package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

type RouterOSCredentialRef string
type RouterOSSecret struct { Username string; Password string }
type RouterOSAdapter struct { client *http.Client; baseURL string; credentialRef RouterOSCredentialRef; serverName string; resolveSecret func(context.Context, RouterOSCredentialRef)(RouterOSSecret,error) }
func NewRouterOSAdapter(baseURL string, credentialRef RouterOSCredentialRef, serverName string, resolveSecret func(context.Context, RouterOSCredentialRef)(RouterOSSecret,error), timeout time.Duration)(*RouterOSAdapter,error){baseURL=strings.TrimRight(strings.TrimSpace(baseURL),"/");if baseURL==""{return nil,errors.New("routeros base url is required")};u,err:=url.Parse(baseURL);if err!=nil||u.Scheme==""||u.Host==""||u.User!=nil{return nil,errors.New("routeros base url must be an absolute URL without userinfo")};if u.Scheme!="http"&&u.Scheme!="https"{return nil,errors.New("routeros base url scheme must be http or https")};if u.Scheme=="http"&&!isLoopbackHost(u.Hostname()){return nil,errors.New("routeros remote transport must use https")};if credentialRef==""{return nil,errors.New("routeros credential reference is required")};if resolveSecret==nil{return nil,errors.New("routeros secret resolver is required")};if timeout<=0{timeout=10*time.Second};return &RouterOSAdapter{baseURL:baseURL,credentialRef:credentialRef,serverName:serverName,resolveSecret:resolveSecret,client:&http.Client{Timeout:timeout,Transport:&http.Transport{TLSClientConfig:&tls.Config{MinVersion:tls.VersionTLS12,ServerName:serverName}}}},nil}
func isLoopbackHost(host string)bool{ip,err:=netip.ParseAddr(host);return err==nil&&ip.IsLoopback()}
func(a *RouterOSAdapter)Protocol()string{return "routeros-api"}
func(a *RouterOSAdapter)doJSON(ctx context.Context,path string)([]map[string]any,error){if ctx==nil{return nil,errors.New("routeros context is required")};if a==nil||a.client==nil||a.resolveSecret==nil{return nil,errors.New("routeros adapter is not initialized")};secret,err:=a.resolveSecret(ctx,a.credentialRef);if err!=nil{return nil,fmt.Errorf("routeros secret lookup: %w",err)};if secret.Username==""||secret.Password==""{return nil,errors.New("routeros secret is incomplete")};req,err:=http.NewRequestWithContext(ctx,http.MethodGet,a.baseURL+"/rest/"+strings.TrimLeft(path,"/"),nil);if err!=nil{return nil,err};req.SetBasicAuth(secret.Username,secret.Password);req.Header.Set("Accept","application/json");resp,err:=a.client.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();if resp.StatusCode<200||resp.StatusCode>=300{body,_:=io.ReadAll(io.LimitReader(resp.Body,4096));return nil,fmt.Errorf("routeros REST %s: status=%d body=%q",path,resp.StatusCode,strings.TrimSpace(string(body)))};var rows []map[string]any;dec:=json.NewDecoder(io.LimitReader(resp.Body,8<<20));if err:=dec.Decode(&rows);err!=nil{return nil,fmt.Errorf("routeros REST %s: decode: %w",path,err)};var extra any;if err:=dec.Decode(&extra);err!=io.EOF{return nil,fmt.Errorf("routeros REST %s: multiple JSON values",path)};return rows,nil}
func(a *RouterOSAdapter)Capabilities(ctx context.Context,device NetworkDevice)([]string,error){if err:=validateRouterOSDevice(device);err!=nil{return nil,err};if _,err:=a.doJSON(ctx,"system/resource");err!=nil{return nil,err};return []string{"health.read","interface.read","routing.read"},nil}
func(a *RouterOSAdapter)CollectInterfaceState(ctx context.Context,device NetworkDevice)([]InterfaceState,error){if err:=validateRouterOSDevice(device);err!=nil{return nil,err};rows,err:=a.doJSON(ctx,"interface");if err!=nil{return nil,err};out:=make([]InterfaceState,0,len(rows));for _,row:=range rows{name:=stringValue(row["name"]);if name==""{continue};out=append(out,InterfaceState{DeviceID:device.ID,Name:name,AdminUp:!boolValue(row["disabled"]),OperUp:boolValue(row["running"]),RxBps:uint64Value(row["rx-byte"]),TxBps:uint64Value(row["tx-byte"]),RxErrors:uint64Value(row["rx-error"]),TxErrors:uint64Value(row["tx-error"])})};return out,nil}
func(a *RouterOSAdapter)CollectRoutingState(ctx context.Context,device NetworkDevice)([]RoutingState,error){if err:=validateRouterOSDevice(device);err!=nil{return nil,err};rows,err:=a.doJSON(ctx,"ip/route");if err!=nil{return nil,err};out:=make([]RoutingState,0,len(rows));for _,row:=range rows{prefix:=stringValue(row["dst-address"]);if prefix==""{continue};out=append(out,RoutingState{DeviceID:device.ID,Protocol:stringValue(row["protocol"]),VRF:stringValue(row["routing-table"]),Prefix:prefix,NextHop:stringValue(row["gateway"]),Metric:uint32Value(row["distance"]),Active:boolValue(row["active"])})};return out,nil}
func validateRouterOSDevice(d NetworkDevice)error{if strings.TrimSpace(d.ID)==""||strings.TrimSpace(d.Address)==""{return errors.New("routeros device id and address are required")};if !isFTNDeviceKind(d.Kind){return errors.New("routeros target ownership is not verified")};return nil}
func stringValue(v any)string{s,_:=v.(string);return strings.TrimSpace(s)}
func boolValue(v any)bool{b,_:=v.(bool);return b}
func uint64Value(v any)uint64{switch n:=v.(type){case float64:if n>=1&&n<=float64(^uint64(0))&&!math.IsNaN(n)&&!math.IsInf(n,0){return uint64(n)};case json.Number:if value,err:=n.Int64();err==nil&&value>0{return uint64(value)};case string:var value uint64;if _,err:=fmt.Sscan(strings.TrimSpace(n),&value);err==nil{return value}};return 0}
func uint32Value(v any)uint32{value:=uint64Value(v);if value>^uint32(0){return ^uint32(0)};return uint32(value)}
