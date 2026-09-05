package agent

import (
	"context"
	"fmt"
	"strings"
)

type LayerHealth struct { LayerID string; Healthy bool; Reason string }
type HealthChecker interface { Check(context.Context) error }
func CheckLayer(ctx context.Context, layer Layer, checker HealthChecker) LayerHealth {
	layerID:=strings.TrimSpace(layer.ID)
	if ctx==nil{return LayerHealth{LayerID:layerID,Reason:"context is required"}}
	if err:=ctx.Err();err!=nil{return LayerHealth{LayerID:layerID,Reason:err.Error()}}
	if checker==nil{return LayerHealth{LayerID:layerID,Healthy:false,Reason:"health checker unavailable"}}
	if err:=checker.Check(ctx);err!=nil{return LayerHealth{LayerID:layerID,Healthy:false,Reason:err.Error()}}
	return LayerHealth{LayerID:layerID,Healthy:true,Reason:fmt.Sprintf("layer %s healthy",layerID)}
}
