package client

import (
"context"
"fmt"

"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure"
"google.golang.org/grpc/encoding"
)

func submitDynamicOverGRPC(target string, req DynamicSubmitRequest) (*DynamicSubmitResponse, error) {
codec := jsonCodec{}
encoding.RegisterCodec(codec)
conn, err := grpc.Dial(
target,
grpc.WithTransportCredentials(insecure.NewCredentials()),
grpc.WithDefaultCallOptions(grpc.ForceCodec(codec)),
)
if err != nil {
return nil, fmt.Errorf("failed to connect grpc target %s: %w", target, err)
}
defer conn.Close()

resp := new(DynamicSubmitResponse)
if err := conn.Invoke(context.Background(), "/dnsd.DynamicDNS/Submit", &req, resp); err != nil {
return nil, fmt.Errorf("grpc submit failed: %w", err)
}
return resp, nil
}

func acmeTokenGRPCConn(target string) (*grpc.ClientConn, error) {
codec := jsonCodec{}
encoding.RegisterCodec(codec)
return grpc.Dial(
target,
grpc.WithTransportCredentials(insecure.NewCredentials()),
grpc.WithDefaultCallOptions(grpc.ForceCodec(codec)),
)
}

func createACMETokenOverGRPC(target string, req AcmeTokenCreateRequest) (*AcmeTokenCreateResponse, error) {
conn, err := acmeTokenGRPCConn(target)
if err != nil {
return nil, fmt.Errorf("failed to connect grpc target %s: %w", target, err)
}
defer conn.Close()
resp := new(AcmeTokenCreateResponse)
if err := conn.Invoke(context.Background(), "/dnsd.ACMETokens/Create", &req, resp); err != nil {
return nil, fmt.Errorf("grpc acme token create failed: %w", err)
}
return resp, nil
}

func revokeACMETokenOverGRPC(target string, req AcmeTokenRevokeRequest) (*AcmeTokenRevokeResponse, error) {
conn, err := acmeTokenGRPCConn(target)
if err != nil {
return nil, fmt.Errorf("failed to connect grpc target %s: %w", target, err)
}
defer conn.Close()
resp := new(AcmeTokenRevokeResponse)
if err := conn.Invoke(context.Background(), "/dnsd.ACMETokens/Revoke", &req, resp); err != nil {
return nil, fmt.Errorf("grpc acme token revoke failed: %w", err)
}
return resp, nil
}

func listACMETokensOverGRPC(target string, req AcmeTokenListRequest) (*AcmeTokenListResponse, error) {
conn, err := acmeTokenGRPCConn(target)
if err != nil {
return nil, fmt.Errorf("failed to connect grpc target %s: %w", target, err)
}
defer conn.Close()
resp := new(AcmeTokenListResponse)
if err := conn.Invoke(context.Background(), "/dnsd.ACMETokens/List", &req, resp); err != nil {
return nil, fmt.Errorf("grpc acme token list failed: %w", err)
}
return resp, nil
}
