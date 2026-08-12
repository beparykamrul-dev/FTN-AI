package domainregistry

import "errors"

// ConfidentialProviderBinding deliberately separates external-provider
// credentials/metadata from the normal domain resource. UI and ordinary staff
// roles should only receive the opaque binding ID and capability status.
type ConfidentialProviderBinding struct {
	BindingID        string   `json:"binding_id"`
	ProviderClass    string   `json:"provider_class"`
	ExternalObjectID string   `json:"external_object_id"`
	SecretRef        string   `json:"secret_ref"`
	AllowedActions   []string `json:"allowed_actions"`
	Active           bool     `json:"active"`
}

type DomainAccessView struct {
	DomainID          string `json:"domain_id"`
	DisplayName       string `json:"display_name"`
	ProviderBindingID string `json:"provider_binding_id,omitempty"`
	ProviderVisible   bool   `json:"provider_visible"`
}

// PublicDomainView never exposes provider credentials, external account IDs,
// secret references, or provider-specific operational metadata.
func PublicDomainView(domainID, displayName string) DomainAccessView {
	return DomainAccessView{DomainID: domainID, DisplayName: displayName, ProviderVisible: false}
}

func (b ConfidentialProviderBinding) Validate() error {
	if b.BindingID == "" || b.ProviderClass == "" || b.ExternalObjectID == "" || b.SecretRef == "" {
		return errors.New("invalid confidential provider binding")
	}
	if !b.Active {
		return errors.New("provider binding is inactive")
	}
	return nil
}
