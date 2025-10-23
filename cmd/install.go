package cmd

import (
	"fmt"
	"log"

	"github.com/acme-go/acme-client/internal/acme"
	"github.com/acme-go/acme-client/internal/config"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install-cert",
	Short: "Install certificate to specified locations",
	Long: `Install the issued certificate to specified file locations and reload services.

Examples:
  # Install certificate for nginx
  acme-go install-cert -d example.com \
    --cert-file /etc/nginx/ssl/cert.pem \
    --key-file /etc/nginx/ssl/key.pem \
    --fullchain-file /etc/nginx/ssl/fullchain.pem \
    --reload-cmd "systemctl reload nginx"

  # Install certificate for apache
  acme-go install-cert -d example.com \
    --cert-file /etc/apache2/ssl/cert.pem \
    --key-file /etc/apache2/ssl/key.pem \
    --ca-file /etc/apache2/ssl/ca.pem \
    --reload-cmd "systemctl reload apache2"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstall(cmd, args)
	},
}

var (
	installDomain     string
	certFile          string
	keyFile           string
	caFile            string
	fullchainFile     string
	reloadCmd         string
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installDomain, "domain", "d", "", "Domain name (required)")
	installCmd.Flags().StringVar(&certFile, "cert-file", "", "Path to install certificate file")
	installCmd.Flags().StringVar(&keyFile, "key-file", "", "Path to install private key file")
	installCmd.Flags().StringVar(&caFile, "ca-file", "", "Path to install CA certificate file")
	installCmd.Flags().StringVar(&fullchainFile, "fullchain-file", "", "Path to install full certificate chain file")
	installCmd.Flags().StringVar(&reloadCmd, "reload-cmd", "", "Command to reload the service after installation")

	installCmd.MarkFlagRequired("domain")
}

func runInstall(cmd *cobra.Command, args []string) error {
	if installDomain == "" {
		return fmt.Errorf("domain is required")
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

	// Prepare install request
	installReq := &acme.InstallRequest{
		Domain:        installDomain,
		CertFile:      certFile,
		KeyFile:       keyFile,
		CAFile:        caFile,
		FullchainFile: fullchainFile,
		ReloadCmd:     reloadCmd,
	}

	log.Printf("Installing certificate for domain: %s", installDomain)

	// Install certificate
	err = client.InstallCert(installReq)
	if err != nil {
		return fmt.Errorf("failed to install certificate: %w", err)
	}

	log.Println("Certificate installed successfully!")
	return nil
}