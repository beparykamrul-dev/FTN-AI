package domaindns

import (
	"context"
	"errors"
	"strings"
)

type Domain struct { ID string `json:"id"`; Name string `json:"name"`; Status string `json:"status"`; DNSSEC bool `json:"dnssec"`; Nameservers []string `json:"nameservers,omitempty"` }
type Record struct { ID string `json:"id"`; Zone string `json:"zone"`; Name string `json:"name"`; Type string `json:"type"`; Value string `json:"value"`; TTL uint32 `json:"ttl"`; Priority uint16 `json:"priority,omitempty"` }
type DomainStore interface { ListDomains(context.Context)([]Domain,error); GetDomain(context.Context,string)(Domain,error); ListRecords(context.Context,string)([]Record,error); CreateRecord(context.Context,Record)error; DeleteRecord(context.Context,string)error }

func ValidateDomain(d Domain) error {
	if strings.TrimSpace(d.ID)==""||strings.TrimSpace(d.Name)==""{return errors.New("domain id and name are required")}
	name:=strings.TrimSpace(d.Name)
	if len(name)>253||strings.ContainsAny(name," \t\r\n/"){return errors.New("invalid domain name")}
	if strings.HasPrefix(name,".")||strings.HasSuffix(name,"."){return errors.New("invalid domain name")}
	return nil
}
func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.ID)==""||strings.TrimSpace(r.Zone)==""||strings.TrimSpace(r.Name)==""||strings.TrimSpace(r.Type)==""||strings.TrimSpace(r.Value)==""{return errors.New("record id, zone, name, type and value are required")}
	if r.TTL==0{return errors.New("record TTL must be greater than zero")}
	if r.TTL>604800{return errors.New("record TTL exceeds seven days")}
	if len(strings.TrimSpace(r.Value))>1<<20{return errors.New("record value is too large")}
	return nil
}
