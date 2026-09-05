package agent

// Category identifies a specialized FTN AI responsibility.
type Category string

const (
	CategoryStudio     Category = "studio"
	CategoryCallCenter Category = "call_center"
	CategoryBilling    Category = "billing_gateway"
	CategoryNetwork    Category = "network_service"
	CategoryDeveloper  Category = "developer"
	CategoryCustomer   Category = "customer"
	CategoryExecutive  Category = "executive_summary"
)

type CategoryConfig struct {
	ID             Category
	SummaryEnabled bool
	DetailsEnabled bool
	ToolScope      string
}

func DefaultCategories() []CategoryConfig {
	return []CategoryConfig{
		{ID: CategoryStudio, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "studio"},
		{ID: CategoryCallCenter, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "support"},
		{ID: CategoryBilling, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "billing"},
		{ID: CategoryNetwork, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "network"},
		{ID: CategoryDeveloper, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "developer"},
		{ID: CategoryCustomer, SummaryEnabled: true, DetailsEnabled: false, ToolScope: "customer"},
		{ID: CategoryExecutive, SummaryEnabled: true, DetailsEnabled: true, ToolScope: "reporting"},
	}
}
