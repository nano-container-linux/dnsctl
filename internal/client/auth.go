package client

import (
	"github.com/nano-container-linux/libdnsd"
)

// Use libdnsd.NormalizeName and libdnsd.AcmeChallengeFQDN

// Utiliser directement libdnsd.BuildDynamicSubmitRequest et libdnsd.SubmitDynamicPayload

// Utilise la fonction mutualisée de libdnsd

func CreateACMEToken(target, fqdn, privateKeyPath string, useAgent bool) (*libdnsd.AcmeTokenCreateResponse, error) {
	challengeFQDN := libdnsd.AcmeChallengeFQDN(fqdn)
	signed := "acme-token-create:" + challengeFQDN
	pub, sig, err := libdnsd.SignString(signed, privateKeyPath, useAgent)
	if err != nil {
		return nil, err
	}
	return createACMETokenOverGRPC(target, libdnsd.AcmeTokenCreateRequest{FQDN: fqdn, PublicKey: pub, Signature: sig})
}

func RevokeACMEToken(target, token, privateKeyPath string, useAgent bool) (*libdnsd.AcmeTokenRevokeResponse, error) {
	signed := "acme-token-revoke:" + token
	pub, sig, err := libdnsd.SignString(signed, privateKeyPath, useAgent)
	if err != nil {
		return nil, err
	}
	return revokeACMETokenOverGRPC(target, libdnsd.AcmeTokenRevokeRequest{Token: token, PublicKey: pub, Signature: sig})
}

func ListACMETokens(target, privateKeyPath string, useAgent bool) (*libdnsd.AcmeTokenListResponse, error) {
	pub, sig, err := libdnsd.SignString("acme-token-list", privateKeyPath, useAgent)
	if err != nil {
		return nil, err
	}
	return listACMETokensOverGRPC(target, libdnsd.AcmeTokenListRequest{PublicKey: pub, Signature: sig})
}
