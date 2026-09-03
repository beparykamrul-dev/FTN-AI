package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// CertificateSource identifies where an FTN certificate is issued from.
type CertificateSource string

const (
	CertificateSourceInternal CertificateSource = "internal-pki"
	CertificateSourceACME     CertificateSource = "acme"
	CertificateSourceAkamai   CertificateSource = "akamai"
)

// CertificateRecord is metadata only; private keys and certificate secrets are never stored here.
type CertificateRecord struct {
	ID            string
	Source        CertificateSource
	Subject       string
	DNSNames      []string
	NotBefore     time.Time
	NotAfter      time.Time
	Fingerprint   string
	Issuer        string
	Serial        string
	RenewalBefore time.Duration
}

func CertificateFingerprint(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}

func ValidateCertificate(cert *x509.Certificate, hostname string, roots *x509.CertPool, intermediates *x509.CertPool) error {
	if cert == nil {
		return errors.New("certificate is required")
	}
	if roots == nil {
		return errors.New("trust roots are required")
	}
	if cert.IsCA {
		return errors.New("leaf certificate must not be a CA")
	}
	if cert.NotAfter.Before(cert.NotBefore) {
		return errors.New("invalid certificate validity interval")
	}
	usageOK := false
	for _, usage := range cert.ExtKeyUsage {
		if usage == x509.ExtKeyUsageServerAuth || usage == x509.ExtKeyUsageClientAuth {
			usageOK = true
			break
		}
	}
	if !usageOK {
		return errors.New("certificate lacks TLS client/server authentication EKU")
	}
	opts := x509.VerifyOptions{Roots: roots, Intermediates: intermediates}
	if strings.TrimSpace(hostname) != "" {
		opts.DNSName = strings.TrimSpace(hostname)
	}
	_, err := cert.Verify(opts)
	return err
}
