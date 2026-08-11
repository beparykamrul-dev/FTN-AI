package controlplane

// Feature describes an FTN capability exposed to the Control Plane UI/API.
// Implementation is provided by modules/adapters; the catalog is the stable
// discovery contract and does not claim a feature is enabled by itself.
type Feature struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Capabilities []string `json:"capabilities"`
}

var EnterpriseFeatureCatalog = []Feature{
	{ID:"ai-dashboard", Name:"AI Dashboard Builder", Category:"analytics", Capabilities:[]string{"pdf","excel","csv","json","xml"}},
	{ID:"ai-sql", Name:"AI SQL Assistant", Category:"analytics", Capabilities:[]string{"query-generation","optimization","database-chat"}},
	{ID:"gis-builder", Name:"GIS Map Builder", Category:"gis", Capabilities:[]string{"layers","geojson","kml","shapefile","geofencing","routing","heatmap","offline-maps"}},
	{ID:"fiber-gis", Name:"Fiber Network Mapping", Category:"isp-gis", Capabilities:[]string{"route","joint","splice","odf","splitter","olt","onu","pop","customer","cable"}},
	{ID:"network-topology", Name:"Network Topology Map", Category:"network", Capabilities:[]string{"live-status","backbone","bgp","device-links"}},
	{ID:"fleet", Name:"Infrastructure Fleet", Category:"control-plane", Capabilities:[]string{"server","router","olt","onu","pc","android","tv","virtual"}},
	{ID:"builder", Name:"Application Builder", Category:"builder", Capabilities:[]string{"php","html","react","vue","angular","flutter"}},
	{ID:"deployment", Name:"Deployment Control", Category:"devops", Capabilities:[]string{"docker","kubernetes","on-premise","backup-restore","versioning"}},
	{ID:"observability", Name:"Unified Observability", Category:"monitoring", Capabilities:[]string{"metrics","logs","traces","flow","ebpf","ndpi"}},
	{ID:"ai-ops", Name:"AI Network Operations", Category:"ai", Capabilities:[]string{"anomaly","forecast","fault-localization","capacity","recovery-plan"}},
}
