package reseller

import "testing"

func TestDefaultOffersSupportStandaloneAndBundles(t *testing.T) {
    for _, o := range DefaultOffers {
        if err := ValidateOffer(o); err != nil { t.Fatalf("%s: %v", o.ID, err) }
        if !o.ResellerEnabled { t.Fatalf("%s should be reseller enabled", o.ID) }
    }
}

func TestRouterOnlyOffer(t *testing.T) {
    var found bool
    for _, o := range DefaultOffers { if o.ID == "router-only" { found = o.Standalone && o.Bundleable; break } }
    if !found { t.Fatal("router-only offer missing or not modular") }
}
