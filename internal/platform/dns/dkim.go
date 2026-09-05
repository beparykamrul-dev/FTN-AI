package dns

import (
	"fmt"
	"strings"
)

type DKIMRecord struct { Domain string `json:"domain"`; Selector string `json:"selector"`; PublicKey string `json:"public_key"`; KeyType string `json:"key_type,omitempty"`; Hash string `json:"hash,omitempty"` }
func (r DKIMRecord) Validate() error {
	if strings.TrimSpace(r.Domain)==""{return fmt.Errorf("domain is required")};if strings.TrimSpace(r.Selector)==""{return fmt.Errorf("selector is required")};if strings.TrimSpace(r.PublicKey)==""{return fmt.Errorf("public key is required")}
	keyType:=strings.ToLower(strings.TrimSpace(r.KeyType));if keyType!=""&&keyType!="rsa"&&keyType!="ed25519"{return fmt.Errorf("unsupported DKIM key type")};hash:=strings.ToLower(strings.TrimSpace(r.Hash));if hash!=""&&hash!="sha256"{return fmt.Errorf("unsupported DKIM hash")};return nil
}
func (r DKIMRecord) TXTValue()(string,error){if err:=r.Validate();err!=nil{return "",err};keyType:=strings.ToLower(strings.TrimSpace(r.KeyType));if keyType==""{keyType="rsa"};parts:=[]string{"v=DKIM1","k="+keyType,"p="+strings.TrimSpace(r.PublicKey)};if hash:=strings.ToLower(strings.TrimSpace(r.Hash));hash!=""{parts=append(parts,"h="+hash)};return strings.Join(parts,"; "),nil}
func (r DKIMRecord) FQDN()string{return strings.TrimSuffix(strings.TrimSpace(r.Selector),".")+"._domainkey."+strings.TrimSuffix(strings.TrimSpace(r.Domain),".")}
