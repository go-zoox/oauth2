package cmd

import (
	"fmt"
	"log"

	"github.com/acme-go/acme-client/internal/acme"
	"github.com/acme-go/acme-client/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certificates",
	Long: `List all certificates managed by acme-go.

Examples:
  # List all certificates
  acme-go list

  # List certificates in raw format
  acme-go list --raw`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd, args)
	},
}

var (
	listRaw bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listRaw, "raw", false, "Output in raw format")
}

func runList(cmd *cobra.Command, args []string) error {
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

	log.Println("Listing all certificates...")

	// List certificates
	certs, err := client.ListCertificates()
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	if len(certs) == 0 {
		log.Println("No certificates found.")
		return nil
	}

	if listRaw {
		for _, cert := range certs {
			fmt.Printf("%s|%s|%s|%s\n", cert.Domain, cert.KeyType, cert.CreatedAt.Format("2006-01-02"), cert.ExpiresAt.Format("2006-01-02"))
		}
	} else {
		fmt.Printf("%-30s %-10s %-12s %-12s %-10s\n", "Domain", "KeyType", "Created", "Expires", "Status")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, cert := range certs {
			status := "Valid"
			if cert.NeedsRenewal() {
				status = "Renewal Due"
			}
			fmt.Printf("%-30s %-10s %-12s %-12s %-10s\n", 
				cert.Domain, 
				cert.KeyType, 
				cert.CreatedAt.Format("2006-01-02"), 
				cert.ExpiresAt.Format("2006-01-02"),
				status,
			)
		}
	}

	return nil
}