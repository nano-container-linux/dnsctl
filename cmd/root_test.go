package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHCLConfig(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "dnsctl.hcl")
	content := []byte(`grpc_addr = "10.0.0.5"
grpc_port = 50099
ssh_agent = false
ssh_key = "/tmp/id_ed25519"
`)
	if err := os.WriteFile(p, content, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := parseHCLConfig(p)
	if err != nil {
		t.Fatalf("parseHCLConfig error: %v", err)
	}
	if cfg.GRPCAddr != "10.0.0.5" || cfg.GRPCPort != 50099 || cfg.SSHKey != "/tmp/id_ed25519" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
	if cfg.SSHAgent == nil || *cfg.SSHAgent {
		t.Fatalf("expected ssh_agent=false, got %+v", cfg.SSHAgent)
	}
}
