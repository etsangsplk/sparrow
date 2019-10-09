package base

import (
	"crypto/rand"
	"encoding/base64"
)

type PublicKeyCredentialCreationOptions struct {
	Attestation     string
	Challenge       string
	RpId            string
	RpName          string
	UserId          string
	UserName        string
	UserDisplayName string
	Timeout         uint64
}

type PublicKeyCredentialRquestOptions struct {
	Challenge        string   `json:"challenge"`
	RpId             string   `json:"rpId"`
	Timeout          uint64   `json:"timeout"`
	CredIds          []string `json:"credIds"`
	UserVerification string   `json:"userVerification"`
}

type CollectedClientData struct {
	Type             string                 `json:"type"`
	Challenge        string                 `json:"challenge"`
	Origin           string                 `json:"origin"`
	HashAlgorithm    string                 `json:"hashAlgorithm"`
	TokenBinding     TokenBinding           `json:"tokenBinding"`
	ClientExtensions map[string]interface{} `codec:"clientExtensions"`
	RawBytes         []byte
}

type TokenBinding struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

type AttestationObject struct {
	AttStmt  map[string]interface{} `codec:"attStmt"`
	Fmt      string                 `codec:"fmt"`
	AuthData []byte                 `codec:authData`
}

func NewPubKeyCred(rpId string) PublicKeyCredentialCreationOptions {
	pkco := PublicKeyCredentialCreationOptions{}
	challenge := make([]byte, 16)
	rand.Read(challenge)
	pkco.Challenge = base64.StdEncoding.EncodeToString(challenge)
	pkco.RpName = "Sparrow Server"
	pkco.RpId = rpId
	pkco.UserName = "kayyagari"
	pkco.UserDisplayName = pkco.UserName
	pkco.Attestation = "none"
	pkco.Timeout = 90000
	pkco.UserId = base64.StdEncoding.EncodeToString([]byte(pkco.UserName))

	return pkco
}

func NewPubKeyAuthReq(credId string) PublicKeyCredentialRquestOptions {
	pkcro := PublicKeyCredentialRquestOptions{}
	challenge := make([]byte, 16)
	rand.Read(challenge)
	pkcro.Challenge = base64.StdEncoding.EncodeToString(challenge)
	pkcro.CredIds = make([]string, 1)
	pkcro.CredIds[0] = credId
	pkcro.UserVerification = "preferred"
	pkcro.Timeout = 90000

	return pkcro
}
