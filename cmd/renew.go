package cmd

import (
	"fmt"
	"log"

	"github.com/acme-go/acme-client/internal/acme"
	"github.com/acme-go/acme-client/internal/config"
	"github.com/spf13/cobra"
)

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew an existing certificate",
	Long: `Renew an existing SSL/TLS certificate.

Examples:
  # Renew certificate for a specific domain
  acme-go renew -d example.com

  # Force renew even if not due for renewal
  acme-go renew -d example.com --force

  # Renew all certificates
  acme-go renew --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRenew(cmd, args)
	},
}

var (
	renewDomain string
	renewAll    bool
	renewForce  bool
)

func init() {
	rootCmd.AddCommand(renewCmd)

	renewCmd.Flags().StringVarP(&renewDomain, "domain", "d", "", "Domain name to renew")
	renewCmd.Flags().BoolVar(&renewAll, "all", false, "Renew all certificates")
	renewCmd.Flags().BoolVar(&renewForce, "force", false, "Force renewal even if not due")
}

func runRenew(cmd *cobra.Command, args []string) error {
	if renewDomain == "" && !renewAll {
		return fmt.Errorf("either --domain or --all must be specified")
	}

	if renewDomain != "" && renewAll {
		return fmt.Errorf("cannot use both --domain and --all at the same time")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create ACME client
	client, err := acme.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ACME client: %w", err)
	}

	if renewAll {
		log.Println("Renewing all certificates...")
		return client.RenewAll(renewForce)
	} else {
		log.Printf("Renewing certificate for domain: %s", renewDomain)
		return client.Renew(renewDomain, renewForce)
	}
}