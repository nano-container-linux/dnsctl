package client

type DynamicSubmitRequest struct {
PayloadHCL string `json:"payload_hcl"`
PublicKey  string `json:"public_key"`
Signature  string `json:"signature"`
}

type DynamicSubmitResponse struct {
ID      string `json:"id"`
Path    string `json:"path"`
Message string `json:"message"`
}

type AcmeTokenCreateRequest struct {
FQDN      string `json:"fqdn"`
PublicKey string `json:"public_key"`
Signature string `json:"signature"`
}

type AcmeTokenCreateResponse struct {
Token string `json:"token"`
FQDN  string `json:"fqdn"`
}

type AcmeTokenRevokeRequest struct {
Token     string `json:"token"`
PublicKey string `json:"public_key"`
Signature string `json:"signature"`
}

type AcmeTokenRevokeResponse struct {
Revoked bool `json:"revoked"`
}

type AcmeTokenListRequest struct {
PublicKey string `json:"public_key"`
Signature string `json:"signature"`
}

type ACMETokenEntry struct {
Token string `json:"token"`
FQDN  string `json:"fqdn"`
}

type AcmeTokenListResponse struct {
Tokens []ACMETokenEntry `json:"tokens"`
}
