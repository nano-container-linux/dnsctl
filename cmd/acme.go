package cmd

import (
	"dnsctl/internal/client"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	acmeCmd := &cobra.Command{Use: "acme", Short: "ACME DNS challenge management"}
	tokenCmd := &cobra.Command{Use: "token", Short: "Manage per-record ACME bearer tokens"}

	var createFQDN string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new per-record ACME bearer token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(createFQDN) == "" {
				return fmt.Errorf("--fqdn is required")
			}
			target, err := grpcTargetFromConfig()
			if err != nil {
				return err
			}
			privateKey, useAgent := sshConfigFromConfig()
			resp, err := client.CreateACMEToken(target, createFQDN, privateKey, useAgent)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "token: %s\n", resp.Token)
			fmt.Fprintf(cmd.OutOrStdout(), "fqdn:  %s\n", resp.FQDN)
			return nil
		},
	}
	createCmd.Flags().StringVar(&createFQDN, "fqdn", "", "Domain (or wildcard, e.g. *.example.com.) to create a token for")

	var revokeToken string
	revokeCmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a per-record ACME bearer token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(revokeToken) == "" {
				return fmt.Errorf("--token is required")
			}
			target, err := grpcTargetFromConfig()
			if err != nil {
				return err
			}
			privateKey, useAgent := sshConfigFromConfig()
			resp, err := client.RevokeACMEToken(target, revokeToken, privateKey, useAgent)
			if err != nil {
				return err
			}
			if resp.Revoked {
				fmt.Fprintln(cmd.OutOrStdout(), "token revoked")
			}
			return nil
		},
	}
	revokeCmd.Flags().StringVar(&revokeToken, "token", "", "Token to revoke")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all per-record ACME bearer tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := grpcTargetFromConfig()
			if err != nil {
				return err
			}
			privateKey, useAgent := sshConfigFromConfig()
			resp, err := client.ListACMETokens(target, privateKey, useAgent)
			if err != nil {
				return err
			}
			if len(resp.Tokens) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no tokens")
				return nil
			}
			for _, tok := range resp.Tokens {
				fmt.Fprintf(cmd.OutOrStdout(), "%-64s  %s\n", tok.Token, tok.FQDN)
			}
			return nil
		},
	}

	tokenCmd.AddCommand(createCmd, revokeCmd, listCmd)
	acmeCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(acmeCmd)
}
