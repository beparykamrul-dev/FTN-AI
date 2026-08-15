package dynamic

import "errors"

// Namespace identifies the FTN DNS authority used for a record.
type Namespace string

const (
	FamilyTimeNet Namespace = "familytimenet.com"
	FTNDNSNet     Namespace = "ftndns.net"
)

var ErrInvalidNamespace = errors.New("invalid FTNDNS namespace")

func ValidateNamespace(ns Namespace) error {
	if ns != FamilyTimeNet && ns != FTNDNSNet {
		return ErrInvalidNamespace
	}
	return nil
}

// ResolveName returns the canonical FQDN for a service/node name.
// Dynamic infrastructure names belong to ftndns.net; public FTN service
// names belong to familytimenet.com.
func ResolveName(ns Namespace, name string) (string, error) {
	if err := ValidateNamespace(ns); err != nil || name == "" {
		return "", ErrInvalidNamespace
	}
	return name + "." + string(ns), nil
}
