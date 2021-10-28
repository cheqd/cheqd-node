package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

const (
	PublicKeyJwk       = "PublicKeyJwk"
	PublicKeyMultibase = "PublicKeyMultibase"
)

var VerificationMethodType = map[string]string{
	"JsonWebKey2020":             PublicKeyJwk,
	"Ed25519VerificationKey2020": PublicKeyMultibase,
}

var ServiceType = []string{
	"LinkedDomains",
	"DIDCommMessaging",
}

func IsVerificationMethodType(vmType string) string {
	return VerificationMethodType[vmType]
}

func IsDidServiceType(sType string) bool {
	return strings.Include(ServiceType, sType)
}
