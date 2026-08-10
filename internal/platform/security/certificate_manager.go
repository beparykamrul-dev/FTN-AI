package security

import (
	"crypto/x509"
	"sync"
	"time"
)

type CertificateState string

const (
	CertificateActive CertificateState = "active"
	CertificateExpiring CertificateState = "expiring"
	CertificateExpired CertificateState = "expired"
)

type CertificateRecord struct {
	ID string `json:"id"`
	Subject string `json:"subject"`
	SANs []string `json:"sans,omitempty"`
	Issuer string `json:"issuer,omitempty"`
	NotBefore time.Time `json:"not_before"`
	NotAfter time.Time `json:"not_after"`
	State CertificateState `json:"state"`
}

type CertificateManager struct {
	mu sync.RWMutex
	records map[string]CertificateRecord
}

func NewCertificateManager() *CertificateManager {
	return &CertificateManager{records: make(map[string]CertificateRecord)}
}

func (m *CertificateManager) Put(r CertificateRecord) {
	m.mu.Lock(); defer m.mu.Unlock()
	m.records[r.ID] = r
}

func StateAt(r CertificateRecord, now time.Time, warning time.Duration) CertificateState {
	if now.After(r.NotAfter) { return CertificateExpired }
	if warning > 0 && !now.Add(warning).Before(r.NotAfter) { return CertificateExpiring }
	return CertificateActive
}

func ParseCertificate(der []byte) (*x509.Certificate, error) { return x509.ParseCertificate(der) }
