package main

import "strings"

// FTNCoreRouterBranding is the canonical presentation contract for the native
// FTN Core Router. It is deliberately independent from MikroTik/RouterOS.
type FTNCoreRouterBranding struct {
	Product    string `json:"product"`
	Company    string `json:"company"`
	LogoKey    string `json:"logo_key"`
	Background string `json:"background"`
	Theme      string `json:"theme"`
}

func DefaultFTNCoreRouterBranding() FTNCoreRouterBranding {
	return FTNCoreRouterBranding{
		Product: "FTN Core Router",
		Company: "Family Time Network",
		LogoKey: "ftn",
		Background: "ftn-network",
		Theme: "dark-network",
	}
}

func ValidateFTNCoreRouterBranding(b FTNCoreRouterBranding) bool {
	return strings.TrimSpace(b.Product) == "FTN Core Router" &&
		strings.TrimSpace(b.Company) == "Family Time Network" &&
		strings.TrimSpace(b.LogoKey) != "" &&
		strings.TrimSpace(b.Background) != "" &&
		strings.TrimSpace(b.Theme) != ""
}
