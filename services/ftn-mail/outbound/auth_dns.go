package outbound

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type MailAuthDNS struct {
	Domain string
	IPv4 []string
	IPv6 []string
	DKIMRecord string
	DMARCPolicy string
}

func (d MailAuthDNS) SPFRecord() (string,error) {
	if d.Domain=="" { return "",errors.New("empty mail domain") }
	parts:=[]string{"v=spf1"}
	for _,ip:=range d.IPv4 { p:=net.ParseIP(ip);if p==nil||p.To4()==nil{return "",fmt.Errorf("invalid IPv4: %s",ip)};parts=append(parts,"ip4:"+ip) }
	for _,ip:=range d.IPv6 { p:=net.ParseIP(ip);if p==nil||p.To4()!=nil{return "",fmt.Errorf("invalid IPv6: %s",ip)};parts=append(parts,"ip6:"+ip) }
	parts=append(parts,"-all");return strings.Join(parts," "),nil
}

func (d MailAuthDNS) DMARCRecord() (string,error) {
	if d.Domain=="" { return "",errors.New("empty mail domain") }
	policy:=d.DMARCPolicy;if policy==""{policy="none"};switch policy{case "none","quarantine","reject":default:return "",errors.New("invalid DMARC policy")}
	return fmt.Sprintf("v=DMARC1; p=%s; adkim=s; aspf=s; pct=100",policy),nil
}
