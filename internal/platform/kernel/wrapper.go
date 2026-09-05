package kernel

import (
	"context"
	"errors"
	"fmt"
)
type ToolRequest struct { Tool string; Operation string; Target string; Parameters map[string]string }
type ToolResult struct { Status string; Output string }
type Tool interface { Name() string; Operations() []string; Execute(ctx context.Context, request ToolRequest)(ToolResult,error) }
type Registry struct { tools map[string]Tool }
func NewRegistry(tools ...Tool)(*Registry,error){r:=&Registry{tools:make(map[string]Tool)};for _,tool:=range tools{if tool==nil||tool.Name()==""{return nil,errors.New("kernel tool must have a name")};if _,exists:=r.tools[tool.Name()];exists{return nil,fmt.Errorf("duplicate kernel tool: %s",tool.Name())};r.tools[tool.Name()]=tool};return r,nil}
func(r *Registry)Execute(ctx context.Context,request ToolRequest)(ToolResult,error){if r==nil{return ToolResult{},errors.New("kernel registry is required")};if ctx==nil{return ToolResult{},errors.New("context is required")};select{case<-ctx.Done():return ToolResult{},ctx.Err();default:};if request.Tool==""||request.Operation==""{return ToolResult{},errors.New("tool and operation are required")};tool,ok:=r.tools[request.Tool];if !ok{return ToolResult{},fmt.Errorf("kernel tool not registered: %s",request.Tool)};params:=make(map[string]string,len(request.Parameters));for k,v:=range request.Parameters{params[k]=v};request.Parameters=params;return tool.Execute(ctx,request)}
func(r *Registry)DoExecute(ctx context.Context,request ToolRequest)(ToolResult,error){return r.Execute(ctx,request)}
