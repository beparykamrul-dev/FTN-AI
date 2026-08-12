package outbound

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// DKIMDNSRecord returns the TXT value published at <selector>._domainkey.<domain>.
// The public key is derived locally from the FTN-owned private key.
func DKIMDNSRecord(domain, selector string, key *rsa.PrivateKey) (string, error) {
	if domain == "" || selector == "" || key == nil { return "", errors.New("invalid DKIM DNS configuration") }
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil { return "", err }
	pub := base64.StdEncoding.EncodeToString(der)
	return fmt.Sprintf("v=DKIM1; k=rsa; p=%s", pub), nil
}

func DKIMPublicKeyPEM(key *rsa.PrivateKey) ([]byte, error) {
	if key == nil { return nil, errors.New("nil DKIM key") }
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil { return nil, err }
	return pem.EncodeToMemory(&pem.Block{Type:"PUBLIC KEY", Bytes:der}), nil
}
