package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var VerificationMethodType = []string{
	"JsonWebKey2020",
	"Ed25519VerificationKey2020",
}

var ServiceType = []string{
	"LinkedDomains",
	"DIDCommMessaging",
}

func IsVerificationMethodType(vmType string) bool {
	return strings.Include(VerificationMethodType, vmType)
}

func IsDidServiceType(sType string) bool {
	return strings.Include(ServiceType, sType)
}
