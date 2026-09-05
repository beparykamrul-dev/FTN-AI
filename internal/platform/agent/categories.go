package agent

import "strings"

type Category string
const(CategoryStudio Category="studio";CategoryCallCenter Category="call_center";CategoryBilling Category="billing_gateway";CategoryNetwork Category="network_service";CategoryDeveloper Category="developer";CategoryCustomer Category="customer";CategoryExecutive Category="executive_summary")
type CategoryConfig struct{ID Category;SummaryEnabled bool;DetailsEnabled bool;ToolScope string}
func(c CategoryConfig)Valid()bool{if strings.TrimSpace(string(c.ID))==""||strings.TrimSpace(c.ToolScope)==""{return false};switch c.ID{case CategoryStudio,CategoryCallCenter,CategoryBilling,CategoryNetwork,CategoryDeveloper,CategoryCustomer,CategoryExecutive:return true;default:return false}}
func(c CategoryConfig)Normalized()CategoryConfig{c.ID=Category(strings.ToLower(strings.TrimSpace(string(c.ID))));c.ToolScope=strings.ToLower(strings.TrimSpace(c.ToolScope));return c}
func DefaultCategories()[]CategoryConfig{return []CategoryConfig{{ID:CategoryStudio,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"studio"},{ID:CategoryCallCenter,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"support"},{ID:CategoryBilling,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"billing"},{ID:CategoryNetwork,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"network"},{ID:CategoryDeveloper,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"developer"},{ID:CategoryCustomer,SummaryEnabled:true,DetailsEnabled:false,ToolScope:"customer"},{ID:CategoryExecutive,SummaryEnabled:true,DetailsEnabled:true,ToolScope:"reporting"}}}
