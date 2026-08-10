package security

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"sync"
)

// CertificateGateway provides the local certificate/TLS boundary for FTN
// services. It deliberately does not contain private keys; key material is
// loaded from the node's protected keystore at runtime.
type CertificateGateway struct {
	mu sync.RWMutex
	certs map[string]tls.Certificate
	roots *x509.CertPool
}

func NewCertificateGateway(roots *x509.CertPool) *CertificateGateway {
	if roots == nil { roots = x509.NewCertPool() }
	return &CertificateGateway{certs: make(map[string]tls.Certificate), roots: roots}
}

func (g *CertificateGateway) Register(id string, cert tls.Certificate) error {
	if id == "" { return errors.New("certificate id is required") }
	if len(cert.Certificate) == 0 { return errors.New("certificate chain is empty") }
	g.mu.Lock(); defer g.mu.Unlock()
	g.certs[id] = cert
	return nil
}

func (g *CertificateGateway) TLSConfig(serverName string, clientAuth tls.ClientAuthType) *tls.Config {
	g.mu.RLock(); defer g.mu.RUnlock()
	cfg := &tls.Config{MinVersion: tls.VersionTLS13, ServerName: serverName, ClientAuth: clientAuth, ClientCAs: g.roots}
	for _, cert := range g.certs { cfg.Certificates = append(cfg.Certificates, cert) }
	return cfg
}
