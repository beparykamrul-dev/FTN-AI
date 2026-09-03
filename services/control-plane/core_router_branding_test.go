package main

import "testing"

func TestDefaultFTNCoreRouterBranding(t *testing.T) {
	b := DefaultFTNCoreRouterBranding()
	if !ValidateFTNCoreRouterBranding(b) { t.Fatal("default FTN Core Router branding must validate") }
	if b.Product != "FTN Core Router" || b.Company != "Family Time Network" { t.Fatalf("unexpected identity: %+v", b) }
	if b.LogoKey != "ftn" || b.Background != "ftn-network" || b.Theme != "dark-network" { t.Fatalf("unexpected branding: %+v", b) }
}

func TestFTNCoreRouterBrandingRejectsWrongProduct(t *testing.T) {
	b := DefaultFTNCoreRouterBranding()
	b.Product = "MikroTik"
	if ValidateFTNCoreRouterBranding(b) { t.Fatal("non-FTN product must not validate as FTN Core Router") }
}
