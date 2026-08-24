package reseller

import "fmt"

type ServiceKind string
const (
    ServiceNetwork ServiceKind = "network"
    ServiceRouter ServiceKind = "router"
    ServiceOLT ServiceKind = "olt"
    ServiceFiberMap ServiceKind = "fiber-map"
    ServiceDNS ServiceKind = "dns"
    ServiceHosting ServiceKind = "hosting"
    ServiceEdge ServiceKind = "edge"
    ServiceCache ServiceKind = "cache"
    ServiceVPN ServiceKind = "vpn"
)

type ServiceOffer struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Kind ServiceKind `json:"kind"`
    Standalone bool `json:"standalone"`
    Bundleable bool `json:"bundleable"`
    ResellerEnabled bool `json:"reseller_enabled"`
    RequiresOLT bool `json:"requires_olt"`
    RequiresFiberMap bool `json:"requires_fiber_map"`
}

func ValidateOffer(o ServiceOffer) error {
    if o.ID == "" || o.Name == "" || o.Kind == "" { return fmt.Errorf("service offer identity is required") }
    if !o.Standalone && !o.Bundleable { return fmt.Errorf("offer must support standalone or bundle mode") }
    if o.RequiresFiberMap && o.Kind != ServiceFiberMap && !o.RequiresOLT { return fmt.Errorf("fiber-map dependent offer requires a valid access dependency") }
    return nil
}

var DefaultOffers = []ServiceOffer{
    {ID:"router-only", Name:"FTN Router Service", Kind:ServiceRouter, Standalone:true, Bundleable:true, ResellerEnabled:true},
    {ID:"olt-only", Name:"FTN OLT Access Service", Kind:ServiceOLT, Standalone:true, Bundleable:true, ResellerEnabled:true, RequiresOLT:true},
    {ID:"fiber-map", Name:"FTN Auto Fiber Map", Kind:ServiceFiberMap, Standalone:true, Bundleable:true, ResellerEnabled:true, RequiresFiberMap:true},
    {ID:"network-bundle", Name:"FTN Network Bundle", Kind:ServiceNetwork, Standalone:false, Bundleable:true, ResellerEnabled:true},
    {ID:"dns", Name:"FTN DNS", Kind:ServiceDNS, Standalone:true, Bundleable:true, ResellerEnabled:true},
    {ID:"hosting", Name:"FTN Hosting", Kind:ServiceHosting, Standalone:true, Bundleable:true, ResellerEnabled:true},
    {ID:"edge", Name:"FTN Edge", Kind:ServiceEdge, Standalone:true, Bundleable:true, ResellerEnabled:true},
    {ID:"cache", Name:"FTN Cache", Kind:ServiceCache, Standalone:true, Bundleable:true, ResellerEnabled:true},
    {ID:"vpn", Name:"FTN VPN", Kind:ServiceVPN, Standalone:true, Bundleable:true, ResellerEnabled:true},
}
