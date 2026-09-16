package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/cobra"
)

type dnsctlConfigFile struct {
	GRPCAddr string `hcl:"grpc_addr,optional"`
	GRPCPort uint16 `hcl:"grpc_port,optional"`
	SSHKey   string `hcl:"ssh_key,optional"`
	SSHAgent *bool  `hcl:"ssh_agent,optional"`
}

var resolved dnsctlConfigFile

func boolPtr(b bool) *bool { return &b }

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dnsctl",
	Short: "dnsd client",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadDNSCTLConfig(cmd)
	},
}

func init() {
	resolved = dnsctlConfigFile{GRPCAddr: "127.0.0.1", GRPCPort: 50051, SSHAgent: boolPtr(true)}
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "Path to dnsctl HCL config file")
	flags.String("grpc-addr", "127.0.0.1", "gRPC listener address")
	flags.Uint16("grpc-port", 50051, "gRPC listener port")
	flags.String("ssh-key", "", "Path to SSH private key (if omitted, uses ssh-agent)")
	flags.Bool("ssh-agent", true, "Allow using ssh-agent when --ssh-key is not provided")
}

func configFileCandidates() []string {
	if strings.TrimSpace(cfgFile) != "" {
		return []string{strings.TrimSpace(cfgFile)}
	}
	candidates := []string{"/etc/dnsctl.hcl"}
	home, err := os.UserHomeDir()
	if err == nil && strings.TrimSpace(home) != "" {
		candidates = append(candidates, filepath.Join(home, ".config", "dnsctl", "dnsctl.hcl"))
	}
	return candidates
}

func loadDNSCTLConfig(cmd *cobra.Command) error {
	cfg := dnsctlConfigFile{GRPCAddr: "127.0.0.1", GRPCPort: 50051, SSHAgent: boolPtr(true)}
	for _, path := range configFileCandidates() {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		fileCfg, err := parseHCLConfig(path)
		if err != nil {
			return err
		}
		if strings.TrimSpace(fileCfg.GRPCAddr) != "" {
			cfg.GRPCAddr = fileCfg.GRPCAddr
		}
		if fileCfg.GRPCPort != 0 {
			cfg.GRPCPort = fileCfg.GRPCPort
		}
		if strings.TrimSpace(fileCfg.SSHKey) != "" {
			cfg.SSHKey = fileCfg.SSHKey
		}
		if fileCfg.SSHAgent != nil {
			cfg.SSHAgent = fileCfg.SSHAgent
		}
		break
	}
	if v := strings.TrimSpace(os.Getenv("DNSCTL_GRPC_ADDR")); v != "" {
		cfg.GRPCAddr = v
	}
	if v := strings.TrimSpace(os.Getenv("DNSCTL_GRPC_PORT")); v != "" {
		p, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return fmt.Errorf("DNSCTL_GRPC_PORT: %w", err)
		}
		cfg.GRPCPort = uint16(p)
	}
	if v := strings.TrimSpace(os.Getenv("DNSCTL_SSH_KEY")); v != "" {
		cfg.SSHKey = v
	}
	if v := strings.TrimSpace(os.Getenv("DNSCTL_SSH_AGENT")); v != "" {
		b := v != "0" && strings.ToLower(v) != "false"
		cfg.SSHAgent = boolPtr(b)
	}

	pf := cmd.Root().PersistentFlags()
	if pf.Changed("grpc-addr") {
		cfg.GRPCAddr, _ = pf.GetString("grpc-addr")
	}
	if pf.Changed("grpc-port") {
		p, _ := pf.GetUint16("grpc-port")
		cfg.GRPCPort = p
	}
	if pf.Changed("ssh-key") {
		cfg.SSHKey, _ = pf.GetString("ssh-key")
	}
	if pf.Changed("ssh-agent") {
		b, _ := pf.GetBool("ssh-agent")
		cfg.SSHAgent = boolPtr(b)
	}

	resolved = cfg
	return nil
}

func parseHCLConfig(path string) (dnsctlConfigFile, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return dnsctlConfigFile{}, fmt.Errorf("dnsctl config %s: %w", path, diags)
	}
	var cfg dnsctlConfigFile
	diags = gohcl.DecodeBody(f.Body, &hcl.EvalContext{}, &cfg)
	if diags.HasErrors() {
		return dnsctlConfigFile{}, fmt.Errorf("dnsctl config %s: %w", path, diags)
	}
	return cfg, nil
}

func grpcTargetFromConfig() (string, error) {
	addr := strings.TrimSpace(resolved.GRPCAddr)
	if addr == "" {
		return "", fmt.Errorf("grpc_addr is required (set via --grpc-addr, DNSCTL_GRPC_ADDR, or config file)")
	}
	if resolved.GRPCPort == 0 {
		return "", fmt.Errorf("grpc_port must be > 0 (set via --grpc-port, DNSCTL_GRPC_PORT, or config file)")
	}
	return net.JoinHostPort(addr, fmt.Sprintf("%d", resolved.GRPCPort)), nil
}

func sshConfigFromConfig() (privateKey string, useAgent bool) {
	ua := true
	if resolved.SSHAgent != nil {
		ua = *resolved.SSHAgent
	}
	return strings.TrimSpace(resolved.SSHKey), ua
}

func Execute() error {
	return rootCmd.Execute()
}
