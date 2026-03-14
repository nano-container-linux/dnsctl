package cmd

import (
"dnsctl/client"
"fmt"
"os"
"strings"

"github.com/spf13/cobra"
)

func init() {
var payloadFile string

dyndnsCmd := &cobra.Command{Use: "dyndns", Short: "Dynamic DNS operations"}
submitCmd := &cobra.Command{
Use:   "submit",
Short: "Submit a dynamic DNS payload over gRPC",
RunE: func(cmd *cobra.Command, args []string) error {
if strings.TrimSpace(payloadFile) == "" {
return fmt.Errorf("--file is required")
}
payloadBytes, err := os.ReadFile(payloadFile)
if err != nil {
return fmt.Errorf("failed to read payload file %s: %w", payloadFile, err)
}
target, err := grpcTargetFromConfig()
if err != nil {
return err
}
privateKey, useAgent := sshConfigFromConfig()
resp, err := client.SubmitDynamicPayload(target, string(payloadBytes), privateKey, useAgent)
if err != nil {
return err
}
fmt.Fprintf(cmd.OutOrStdout(), "id: %s\n", resp.ID)
fmt.Fprintf(cmd.OutOrStdout(), "path: %s\n", resp.Path)
fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
return nil
},
}

submitCmd.Flags().StringVar(&payloadFile, "file", "", "Path to .hcl payload file")
dyndnsCmd.AddCommand(submitCmd)
rootCmd.AddCommand(dyndnsCmd)
}
