// Package config provides TLS utilities including self-signed certificate generation
package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/gibram-io/gibram/pkg/logging"
)

// GenerateSelfSignedCert generates a self-signed TLS certificate
// Returns the certificate and key as PEM-encoded bytes
func GenerateSelfSignedCert(hosts []string, validFor time.Duration) (certPEM, keyPEM []byte, err error) {
	// Generate ECDSA private key (P-256 curve)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Generate serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Create certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"GibRAM Self-Signed"},
			CommonName:   "GibRAM Server",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Add hosts (IP addresses and DNS names)
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	// Always add localhost
	template.IPAddresses = append(template.IPAddresses, net.IPv4(127, 0, 0, 1), net.IPv6loopback)
	template.DNSNames = append(template.DNSNames, "localhost")

	// Create certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// Encode private key to PEM
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return certPEM, keyPEM, nil
}

// LoadOrGenerateTLSConfig loads TLS config from files or generates a self-signed certificate
// Returns the tls.Config and a boolean indicating if TLS should be enabled
func (cfg *TLSConfig) LoadOrGenerateTLSConfig(dataDir string) (*tls.Config, bool, error) {
	// First, check if cert/key files are provided
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, false, fmt.Errorf("failed to load TLS certificates: %w", err)
		}
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}, true, nil
	}

	// If auto_cert is enabled, generate or load cached self-signed cert
	if cfg.AutoCert {
		return cfg.loadOrGenerateAutoCert(dataDir)
	}

	// No TLS configured
	return nil, false, nil
}

// loadOrGenerateAutoCert handles auto-generated certificates with caching
func (cfg *TLSConfig) loadOrGenerateAutoCert(dataDir string) (*tls.Config, bool, error) {
	// Define paths for cached certificates
	certPath := filepath.Join(dataDir, "auto_cert.pem")
	keyPath := filepath.Join(dataDir, "auto_key.pem")

	// Try to load existing cached certificates
	if cert, err := tls.LoadX509KeyPair(certPath, keyPath); err == nil {
		// Verify the certificate is still valid
		if len(cert.Certificate) > 0 {
			x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
			if err == nil && time.Now().Before(x509Cert.NotAfter) {
				logging.Info("Using cached auto-generated TLS certificate (expires: %s)", x509Cert.NotAfter.Format(time.RFC3339))
				return &tls.Config{
					Certificates: []tls.Certificate{cert},
					MinVersion:   tls.VersionTLS12,
				}, true, nil
			}
		}
	}

	// Generate new self-signed certificate (valid for 1 year)
	logging.Info("Generating self-signed TLS certificate...")
	certPEM, keyPEM, err := GenerateSelfSignedCert([]string{"localhost"}, 365*24*time.Hour)
	if err != nil {
		return nil, false, fmt.Errorf("failed to generate self-signed cert: %w", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, false, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save certificates to disk for reuse
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		logging.Warn("Failed to cache certificate: %v", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		logging.Warn("Failed to cache key: %v", err)
	}

	// Parse the certificate for use
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse generated certificate: %w", err)
	}

	logging.Info("Self-signed TLS certificate generated (cached at %s)", dataDir)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, true, nil
}
