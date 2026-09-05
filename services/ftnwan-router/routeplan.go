package router

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/netip"
	"strings"
)

type RoutePlanner struct{ Plane Dataplane }
type RoutePlan struct { ID string `json:"id"`; Route Route `json:"route"` }

func (p RoutePlanner) PlanRouteChange(ctx context.Context, r Route) (string, error) {
	if ctx==nil{return "",fmt.Errorf("context is required")};if p.Plane == nil { return "", fmt.Errorf("dataplane is required") }
	if err := validateRoute(r); err != nil { return "", err };if err:=ctx.Err();err!=nil{return "",err}
	if !p.Plane.Ready(ctx) { return "", fmt.Errorf("dataplane %q is not ready", p.Plane.Name()) }
	payload, err := json.Marshal(r); if err != nil { return "", fmt.Errorf("marshal route: %w", err) };h := sha256.Sum256(payload);plan := RoutePlan{ID: "route-" + hex.EncodeToString(h[:8]), Route: r};out, err := json.Marshal(plan); if err != nil { return "", err };return string(out), nil
}
func validateRoute(r Route) error {
	prefixText:=strings.TrimSpace(r.Prefix);if prefixText==""{return fmt.Errorf("route prefix is required")};if len(prefixText)>64{return fmt.Errorf("route prefix too long")};prefix,err:=netip.ParsePrefix(prefixText);if err!=nil{return fmt.Errorf("invalid route prefix: %w",err)}
	if r.NextHop!=""{if len(r.NextHop)>64{return fmt.Errorf("next hop too long")};hop,err:=netip.ParseAddr(strings.TrimSpace(r.NextHop));if err!=nil{return fmt.Errorf("invalid next hop: %w",err)};if hop.Is4()!=prefix.Addr().Is4(){return fmt.Errorf("next hop address family differs from prefix")}}
	if r.Metric>4294967294{return fmt.Errorf("route metric out of range")};return nil
}
func validatePlan(plan RoutePlan) error {if strings.TrimSpace(plan.ID)==""||len(plan.ID)>64{return fmt.Errorf("plan id is invalid")};if err:=validateRoute(plan.Route);err!=nil{return err};payload,err:=json.Marshal(plan.Route);if err!=nil{return fmt.Errorf("marshal route: %w",err)};h:=sha256.Sum256(payload);if plan.ID!="route-"+hex.EncodeToString(h[:8]){return fmt.Errorf("route plan integrity check failed")};return nil}
func decodePlan(planJSON string)(RoutePlan,error){if strings.TrimSpace(planJSON)==""{return RoutePlan{},fmt.Errorf("route plan is required")};if len(planJSON)>256<<10{return RoutePlan{},fmt.Errorf("route plan too large")};dec:=json.NewDecoder(strings.NewReader(planJSON));var plan RoutePlan;if err:=dec.Decode(&plan);err!=nil{return RoutePlan{},fmt.Errorf("invalid route plan: %w",err)};var extra any;if err:=dec.Decode(&extra);err!=io.EOF{if err==nil{return RoutePlan{},fmt.Errorf("multiple route plan values")};return RoutePlan{},fmt.Errorf("invalid trailing route plan data: %w",err)};if err:=validatePlan(plan);err!=nil{return RoutePlan{},err};return plan,nil}
func(p RoutePlanner)ApplyApprovedPlan(ctx context.Context,planJSON string)error{if ctx==nil{return fmt.Errorf("context is required")};if p.Plane==nil{return fmt.Errorf("dataplane is required")};plan,err:=decodePlan(planJSON);if err!=nil{return err};return p.Plane.ApplyRoute(ctx,plan.Route)}
func(p RoutePlanner)Rollback(ctx context.Context,planJSON string)error{if ctx==nil{return fmt.Errorf("context is required")};if p.Plane==nil{return fmt.Errorf("dataplane is required")};plan,err:=decodePlan(planJSON);if err!=nil{return err};return p.Plane.DeleteRoute(ctx,plan.Route)}
