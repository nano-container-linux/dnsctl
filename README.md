# dnsctl

`dnsctl` is the standalone client for `dnsd` dynamic DNS and ACME token operations.

## Features

- Submit dynamic DNS payloads over gRPC
- Create/list/revoke per-record ACME tokens
- SSH signing with private key or ssh-agent
- HCL config file support (`/etc/dnsctl.hcl`, `~/.config/dnsctl/dnsctl.hcl`)

## Build

```sh
go build -o dnsctl .
```

## Test

```sh
go test ./...
```

## Config

```hcl
grpc_addr = "127.0.0.1"
grpc_port = 50051
ssh_agent = true
# ssh_key = "/home/user/.ssh/id_ed25519"
```

You can also override via env vars:

- `DNSCTL_GRPC_ADDR`
- `DNSCTL_GRPC_PORT`
- `DNSCTL_SSH_KEY`
- `DNSCTL_SSH_AGENT`

## Usage examples

```sh
./dnsctl dyndns submit --file ./update.hcl
./dnsctl acme token create --fqdn example.com.
./dnsctl acme token list
./dnsctl acme token revoke --token <token>
```
