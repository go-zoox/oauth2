package cmd

import (
	"fmt"
	"log"

	"github.com/acme-go/acme-client/internal/acme"
	"github.com/acme-go/acme-client/internal/config"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy certificate using deployment hooks",
	Long: `Deploy certificate to services using deployment hooks.

Examples:
  # Deploy to nginx
  acme-go deploy -d example.com --deploy-hook nginx

  # Deploy to docker container
  acme-go deploy -d example.com --deploy-hook docker

  # Deploy to custom script
  acme-go deploy -d example.com --deploy-hook custom_script`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeploy(cmd, args)
	},
}

var (
	deployDomain string
	deployHook   string
)

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployDomain, "domain", "d", "", "Domain name (required)")
	deployCmd.Flags().StringVar(&deployHook, "deploy-hook", "", "Deployment hook name (required)")

	deployCmd.MarkFlagRequired("domain")
	deployCmd.MarkFlagRequired("deploy-hook")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	if deployDomain == "" {
		return fmt.Errorf("domain is required")
	}

	if deployHook == "" {
		return fmt.Errorf("deploy-hook is required")
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

	log.Printf("Deploying certificate for domain: %s using hook: %s", deployDomain, deployHook)

	// Deploy certificate
	err = client.Deploy(deployDomain, deployHook)
	if err != nil {
		return fmt.Errorf("failed to deploy certificate: %w", err)
	}

	log.Println("Certificate deployed successfully!")
	return nil
}