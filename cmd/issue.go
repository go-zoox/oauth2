package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/acme-go/acme-client/internal/acme"
	"github.com/acme-go/acme-client/internal/config"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue a new SSL/TLS certificate",
	Long: `Issue a new SSL/TLS certificate for the specified domain(s).

Examples:
  # Issue certificate using webroot validation
  acme-go issue -d example.com -w /var/www/html

  # Issue certificate using DNS validation
  acme-go issue -d example.com --dns dns_cf

  # Issue certificate for multiple domains
  acme-go issue -d example.com -d www.example.com -w /var/www/html

  # Issue wildcard certificate using DNS
  acme-go issue -d "*.example.com" --dns dns_cf`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIssue(cmd, args)
	},
}

var (
	domains     []string
	webroot     string
	dnsProvider string
	keyType     string
	email       string
	server      string
	staging     bool
	force       bool
)

func init() {
	rootCmd.AddCommand(issueCmd)

	issueCmd.Flags().StringSliceVarP(&domains, "domain", "d", []string{}, "Domain name(s) for the certificate (required)")
	issueCmd.Flags().StringVarP(&webroot, "webroot", "w", "", "Webroot path for HTTP-01 validation")
	issueCmd.Flags().StringVar(&dnsProvider, "dns", "", "DNS provider for DNS-01 validation (e.g., dns_cf, dns_ali)")
	issueCmd.Flags().StringVar(&keyType, "key-type", "ec-256", "Key type: rsa-2048, rsa-3072, rsa-4096, ec-256, ec-384")
	issueCmd.Flags().StringVar(&email, "email", "", "Email address for account registration")
	issueCmd.Flags().StringVar(&server, "server", "", "ACME server URL (default: Let's Encrypt)")
	issueCmd.Flags().BoolVar(&staging, "staging", false, "Use staging environment")
	issueCmd.Flags().BoolVar(&force, "force", false, "Force issue even if certificate exists")

	issueCmd.MarkFlagRequired("domain")
}

func runIssue(cmd *cobra.Command, args []string) error {
	if len(domains) == 0 {
		return fmt.Errorf("at least one domain is required")
	}

	// Validate that either webroot or DNS provider is specified
	if webroot == "" && dnsProvider == "" {
		return fmt.Errorf("either --webroot or --dns must be specified")
	}

	if webroot != "" && dnsProvider != "" {
		return fmt.Errorf("cannot use both --webroot and --dns at the same time")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command line flags
	if email != "" {
		cfg.Email = email
	}
	if server != "" {
		cfg.Server = server
	}
	if staging {
		cfg.Staging = true
	}

	// Validate email
	if cfg.Email == "" {
		return fmt.Errorf("email is required for account registration")
	}

	// Create ACME client
	client, err := acme.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ACME client: %w", err)
	}

	// Prepare certificate request
	certReq := &acme.CertificateRequest{
		Domains:     domains,
		KeyType:     keyType,
		Force:       force,
	}

	// Set validation method
	if webroot != "" {
		certReq.ValidationMethod = acme.ValidationHTTP01
		certReq.Webroot = webroot
	} else if dnsProvider != "" {
		certReq.ValidationMethod = acme.ValidationDNS01
		certReq.DNSProvider = dnsProvider
	}

	log.Printf("Issuing certificate for domains: %s", strings.Join(domains, ", "))
	
	// Issue certificate
	cert, err := client.Issue(certReq)
	if err != nil {
		return fmt.Errorf("failed to issue certificate: %w", err)
	}

	log.Printf("Certificate issued successfully!")
	log.Printf("Certificate saved to: %s", cert.CertPath)
	log.Printf("Private key saved to: %s", cert.KeyPath)
	log.Printf("Certificate expires: %s", cert.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}