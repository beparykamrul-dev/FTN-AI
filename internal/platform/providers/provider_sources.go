package providers

// Source describes an external open-source ecosystem that can contribute a
// provider adapter, deployment recipe, telemetry connector, or API client.
// Source code is not copied into the FTN core; adapters remain isolated.
type Source struct {
	Provider string
	Repository string
	Capabilities []string
	AdapterKind string
}

var Sources = []Source{
	{Provider:"Akamai", Repository:"akamai/cli", Capabilities:[]string{"edge","dns","api","deployment"}, AdapterKind:"api"},
	{Provider:"Akamai", Repository:"akamai/terraform-provider-akamai", Capabilities:[]string{"infrastructure","dns","edge"}, AdapterKind:"iac"},
	{Provider:"Akamai", Repository:"akamai/AkamaiOPEN-edgegrid-golang", Capabilities:[]string{"api","authentication"}, AdapterKind:"sdk"},
	{Provider:"Akamai Developers", Repository:"akamai-developers/rag-langgraph-k8s-quickstart", Capabilities:[]string{"ai","cloud","kubernetes"}, AdapterKind:"reference"},
	{Provider:"Akamai Compute Marketplace", Repository:"akamai-compute-marketplace/marketplace-apps", Capabilities:[]string{"one-click-deploy","ansible","stackscript","cloud"}, AdapterKind:"deployment"},
	{Provider:"Akamai Contrib", Repository:"akamai-contrib/nodejs-connector-of-datastream", Capabilities:[]string{"telemetry","datastream","prometheus","grafana"}, AdapterKind:"telemetry"},
}
