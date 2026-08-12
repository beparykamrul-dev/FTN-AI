package domaindns

import "context"

type Domain struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	DNSSEC     bool   `json:"dnssec"`
	Nameservers []string `json:"nameservers,omitempty"`
}

type Record struct {
	ID       string `json:"id"`
	Zone     string `json:"zone"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      uint32 `json:"ttl"`
	Priority uint16 `json:"priority,omitempty"`
}

type DomainStore interface {
	ListDomains(context.Context) ([]Domain, error)
	GetDomain(context.Context, string) (Domain, error)
	ListRecords(context.Context, string) ([]Record, error)
	CreateRecord(context.Context, Record) error
	DeleteRecord(context.Context, string) error
}
