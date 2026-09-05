package callcenter

import("errors";"strings")
type CustomerContext struct{CustomerID string `json:"customer_id"`;Name string `json:"name,omitempty"`;ServiceID string `json:"service_id,omitempty"`;PackageName string `json:"package_name,omitempty"`;PPPoEUser string `json:"pppoe_user,omitempty"`;ONUID string `json:"onu_id,omitempty"`;RouterID string `json:"router_id,omitempty"`;IP string `json:"ip,omitempty"`;ServiceStatus string `json:"service_status,omitempty"`;ActiveIncident string `json:"active_incident,omitempty"`}
type ContextProvider interface{Load(customerID string)(CustomerContext,error)}
func(c CustomerContext)Valid()bool{return strings.TrimSpace(c.CustomerID)!=""}
func BuildSessionContext(provider ContextProvider,customerID string)(CustomerContext,error){customerID=strings.TrimSpace(customerID);if customerID==""{return CustomerContext{},errors.New("customer_id is required")};if provider==nil{return CustomerContext{CustomerID:customerID},nil};ctx,err:=provider.Load(customerID);if err!=nil{return CustomerContext{},err};ctx.CustomerID=strings.TrimSpace(ctx.CustomerID);if ctx.CustomerID==""{ctx.CustomerID=customerID};return ctx,nil}
