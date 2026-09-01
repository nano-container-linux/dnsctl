package client

import (
	lib "github.com/nano-container-linux/libdnsd"
)

func submitDynamicOverGRPC(target string, req lib.DynamicSubmitRequest) (*lib.DynamicSubmitResponse, error) {
	return lib.SubmitDynamicOverGRPC(target, req)
}

func createACMETokenOverGRPC(target string, req lib.AcmeTokenCreateRequest) (*lib.AcmeTokenCreateResponse, error) {
	return lib.CreateACMETokenOverGRPC(target, req)
}

func revokeACMETokenOverGRPC(target string, req lib.AcmeTokenRevokeRequest) (*lib.AcmeTokenRevokeResponse, error) {
	return lib.RevokeACMETokenOverGRPC(target, req)
}

func listACMETokensOverGRPC(target string, req lib.AcmeTokenListRequest) (*lib.AcmeTokenListResponse, error) {
	return lib.ListACMETokensOverGRPC(target, req)
}

func SubmitDynamicPayload(target string, payload string, privateKeyPath string, useAgent bool) (*lib.DynamicSubmitResponse, error) {
	return lib.SubmitDynamicPayload(target, payload, privateKeyPath, useAgent, lib.SubmitDynamicOverGRPC)
}
