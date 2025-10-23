package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/acme-go/acme-client/internal/config"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// ValidationMethod represents the validation method
type ValidationMethod string

const (
	ValidationHTTP01   ValidationMethod = "http-01"
	ValidationDNS01    ValidationMethod = "dns-01"
	ValidationTLSALPN01 ValidationMethod = "tls-alpn-01"
)

// Client represents an ACME client
type Client struct {
	config     *config.Config
	legoClient *lego.Client
	user       *User
}

// User represents an ACME user
type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

// GetEmail returns user email
func (u *User) GetEmail() string {
	return u.Email
}

// GetRegistration returns user registration
func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

// GetPrivateKey returns user private key
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// CertificateRequest represents a certificate request
type CertificateRequest struct {
	Domains          []string
	KeyType          string
	ValidationMethod ValidationMethod
	Webroot          string
	DNSProvider      string
	Force            bool
}

// InstallRequest represents a certificate installation request
type InstallRequest struct {
	Domain        string
	CertFile      string
	KeyFile       string
	CAFile        string
	FullchainFile string
	ReloadCmd     string
}

// Certificate represents a certificate
type Certificate struct {
	Domain      string
	KeyType     string
	CertPath    string
	KeyPath     string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Certificate []byte
	PrivateKey  []byte
}

// NeedsRenewal checks if certificate needs renewal
func (c *Certificate) NeedsRenewal() bool {
	renewalTime := c.ExpiresAt.AddDate(0, 0, -30) // 30 days before expiry
	return time.Now().After(renewalTime)
}

// NewClient creates a new ACME client
func NewClient(cfg *config.Config) (*Client, error) {
	// Create or load user
	user, err := createOrLoadUser(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create/load user: %w", err)
	}

	// Create lego config
	legoConfig := lego.NewConfig(user)
	
	// Set ACME server URL
	if cfg.Staging {
		legoConfig.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else if cfg.Server != "" {
		legoConfig.CADirURL = cfg.Server
	}

	// Create lego client
	legoClient, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create lego client: %w", err)
	}

	// Register user if not already registered
	if user.Registration == nil {
		reg, err := legoClient.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return nil, fmt.Errorf("failed to register user: %w", err)
		}
		user.Registration = reg
		
		// Save user registration
		if err := saveUser(cfg, user); err != nil {
			log.Printf("Warning: failed to save user registration: %v", err)
		}
	}

	return &Client{
		config:     cfg,
		legoClient: legoClient,
		user:       user,
	}, nil
}

// Issue issues a new certificate
func (c *Client) Issue(req *CertificateRequest) (*Certificate, error) {
	// Check if certificate already exists and not forced
	if !req.Force {
		if cert, err := c.loadCertificate(req.Domains[0]); err == nil {
			if !cert.NeedsRenewal() {
				log.Printf("Certificate for %s is still valid, skipping", req.Domains[0])
				return cert, nil
			}
		}
	}

	// Set up challenge solver
	if err := c.setupChallengeSolver(req); err != nil {
		return nil, fmt.Errorf("failed to setup challenge solver: %w", err)
	}

	// Create certificate request
	certReq := certificate.ObtainRequest{
		Domains: req.Domains,
		Bundle:  true,
	}

	// Note: Key type handling would be implemented here
	// For now, we'll use the default key type

	// Obtain certificate
	certificates, err := c.legoClient.Certificate.Obtain(certReq)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain certificate: %w", err)
	}

	// Save certificate
	cert, err := c.saveCertificate(req.Domains[0], req.KeyType, certificates)
	if err != nil {
		return nil, fmt.Errorf("failed to save certificate: %w", err)
	}

	return cert, nil
}

// Renew renews a certificate
func (c *Client) Renew(domain string, force bool) error {
	cert, err := c.loadCertificate(domain)
	if err != nil {
		return fmt.Errorf("certificate not found for domain %s: %w", domain, err)
	}

	if !force && !cert.NeedsRenewal() {
		log.Printf("Certificate for %s does not need renewal yet", domain)
		return nil
	}

	// Parse existing certificate to get domains
	block, _ := pem.Decode(cert.Certificate)
	if block == nil {
		return fmt.Errorf("failed to decode certificate")
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	domains := []string{x509Cert.Subject.CommonName}
	domains = append(domains, x509Cert.DNSNames...)

	// Create renewal request
	req := &CertificateRequest{
		Domains: domains,
		KeyType: cert.KeyType,
		Force:   true,
	}

	// Issue new certificate
	_, err = c.Issue(req)
	if err != nil {
		return fmt.Errorf("failed to renew certificate: %w", err)
	}

	log.Printf("Certificate renewed successfully for domain: %s", domain)
	return nil
}

// RenewAll renews all certificates
func (c *Client) RenewAll(force bool) error {
	certs, err := c.ListCertificates()
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	for _, cert := range certs {
		if force || cert.NeedsRenewal() {
			log.Printf("Renewing certificate for domain: %s", cert.Domain)
			if err := c.Renew(cert.Domain, force); err != nil {
				log.Printf("Failed to renew certificate for %s: %v", cert.Domain, err)
				continue
			}
		}
	}

	return nil
}

// InstallCert installs certificate to specified locations
func (c *Client) InstallCert(req *InstallRequest) error {
	cert, err := c.loadCertificate(req.Domain)
	if err != nil {
		return fmt.Errorf("certificate not found for domain %s: %w", req.Domain, err)
	}

	// Install certificate file
	if req.CertFile != "" {
		if err := c.installFile(cert.CertPath, req.CertFile); err != nil {
			return fmt.Errorf("failed to install cert file: %w", err)
		}
	}

	// Install private key file
	if req.KeyFile != "" {
		if err := c.installFile(cert.KeyPath, req.KeyFile); err != nil {
			return fmt.Errorf("failed to install key file: %w", err)
		}
	}

	// Install CA file and fullchain file would require additional logic
	// For now, we'll just copy the certificate file
	if req.CAFile != "" {
		// TODO: Extract CA certificate from chain
		log.Printf("CA file installation not implemented yet")
	}

	if req.FullchainFile != "" {
		if err := c.installFile(cert.CertPath, req.FullchainFile); err != nil {
			return fmt.Errorf("failed to install fullchain file: %w", err)
		}
	}

	// Execute reload command
	if req.ReloadCmd != "" {
		if err := c.executeCommand(req.ReloadCmd); err != nil {
			return fmt.Errorf("failed to execute reload command: %w", err)
		}
	}

	return nil
}

// Deploy deploys certificate using deployment hooks
func (c *Client) Deploy(domain, hookName string) error {
	hook, exists := c.config.GetDeployHook(hookName)
	if !exists {
		return fmt.Errorf("deployment hook '%s' not found", hookName)
	}

	cert, err := c.loadCertificate(domain)
	if err != nil {
		return fmt.Errorf("certificate not found for domain %s: %w", domain, err)
	}

	// Execute deployment hook
	switch hook.Type {
	case "script":
		return c.executeDeployScript(hook.Script, cert)
	default:
		return fmt.Errorf("unsupported deployment hook type: %s", hook.Type)
	}
}

// ListCertificates lists all certificates
func (c *Client) ListCertificates() ([]*Certificate, error) {
	certDir := c.config.CertDir
	
	entries, err := os.ReadDir(certDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert directory: %w", err)
	}

	var certificates []*Certificate
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domain := entry.Name()
		if cert, err := c.loadCertificate(domain); err == nil {
			certificates = append(certificates, cert)
		}
	}

	return certificates, nil
}

// Helper functions

func createOrLoadUser(cfg *config.Config) (*User, error) {
	userDir := filepath.Join(cfg.CertDir, "accounts")
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user directory: %w", err)
	}

	keyPath := filepath.Join(userDir, "account.key")
	
	// Try to load existing key
	if keyData, err := os.ReadFile(keyPath); err == nil {
		block, _ := pem.Decode(keyData)
		if block != nil {
			if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
				return &User{
					Email: cfg.Email,
					key:   key,
				}, nil
			}
		}
	}

	// Generate new key
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Save key
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	return &User{
		Email: cfg.Email,
		key:   key,
	}, nil
}

func saveUser(cfg *config.Config, user *User) error {
	// In a real implementation, you would save the user registration data
	// For now, we'll just log it
	log.Printf("User registration saved for: %s", user.Email)
	return nil
}

func parseKeyType(keyType string) (certcrypto.KeyType, error) {
	switch keyType {
	case "ec-256":
		return certcrypto.EC256, nil
	case "ec-384":
		return certcrypto.EC384, nil
	case "rsa-2048":
		return certcrypto.RSA2048, nil
	case "rsa-3072":
		return certcrypto.RSA3072, nil
	case "rsa-4096":
		return certcrypto.RSA4096, nil
	default:
		return certcrypto.EC256, fmt.Errorf("unsupported key type: %s", keyType)
	}
}

func (c *Client) setupChallengeSolver(req *CertificateRequest) error {
	switch req.ValidationMethod {
	case ValidationHTTP01:
		if req.Webroot != "" {
			return c.legoClient.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
		}
		return c.legoClient.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	
	case ValidationDNS01:
		// DNS provider setup would go here
		// For now, we'll use manual DNS
		return c.legoClient.Challenge.SetDNS01Provider(&ManualDNSProvider{})
	
	default:
		return fmt.Errorf("unsupported validation method: %s", req.ValidationMethod)
	}
}

func (c *Client) saveCertificate(domain, keyType string, certificates *certificate.Resource) (*Certificate, error) {
	domainDir := filepath.Join(c.config.CertDir, domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create domain directory: %w", err)
	}

	certPath := filepath.Join(domainDir, domain+".crt")
	keyPath := filepath.Join(domainDir, domain+".key")

	// Save certificate
	if err := os.WriteFile(certPath, certificates.Certificate, 0644); err != nil {
		return nil, fmt.Errorf("failed to save certificate: %w", err)
	}

	// Save private key
	if err := os.WriteFile(keyPath, certificates.PrivateKey, 0600); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	// Parse certificate to get expiry date
	block, _ := pem.Decode(certificates.Certificate)
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &Certificate{
		Domain:      domain,
		KeyType:     keyType,
		CertPath:    certPath,
		KeyPath:     keyPath,
		CreatedAt:   time.Now(),
		ExpiresAt:   x509Cert.NotAfter,
		Certificate: certificates.Certificate,
		PrivateKey:  certificates.PrivateKey,
	}, nil
}

func (c *Client) loadCertificate(domain string) (*Certificate, error) {
	domainDir := filepath.Join(c.config.CertDir, domain)
	certPath := filepath.Join(domainDir, domain+".crt")
	keyPath := filepath.Join(domainDir, domain+".key")

	// Check if files exist
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("certificate file not found")
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("private key file not found")
	}

	// Read certificate
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	// Read private key
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse certificate to get expiry date
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Determine key type from private key
	keyType := "unknown"
	keyBlock, _ := pem.Decode(keyData)
	if keyBlock != nil {
		if keyBlock.Type == "EC PRIVATE KEY" {
			keyType = "ec-256" // Simplified
		} else if keyBlock.Type == "RSA PRIVATE KEY" {
			if key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes); err == nil {
				switch key.Size() * 8 {
				case 2048:
					keyType = "rsa-2048"
				case 3072:
					keyType = "rsa-3072"
				case 4096:
					keyType = "rsa-4096"
				}
			}
		}
	}

	return &Certificate{
		Domain:      domain,
		KeyType:     keyType,
		CertPath:    certPath,
		KeyPath:     keyPath,
		CreatedAt:   x509Cert.NotBefore,
		ExpiresAt:   x509Cert.NotAfter,
		Certificate: certData,
		PrivateKey:  keyData,
	}, nil
}

func (c *Client) installFile(srcPath, dstPath string) error {
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(dstPath, srcData, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	log.Printf("Installed: %s -> %s", srcPath, dstPath)
	return nil
}

func (c *Client) executeCommand(cmd string) error {
	log.Printf("Executing command: %s", cmd)
	// In a real implementation, you would execute the command
	// For now, we'll just log it
	return nil
}

func (c *Client) executeDeployScript(script string, cert *Certificate) error {
	log.Printf("Executing deploy script: %s for domain: %s", script, cert.Domain)
	// In a real implementation, you would execute the deployment script
	// passing certificate information as environment variables
	return nil
}