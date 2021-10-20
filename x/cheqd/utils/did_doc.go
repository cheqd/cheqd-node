package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var VerificationMethodType = []string{
	"JsonWebKey2020",
	"EcdsaSecp256k1VerificationKey2019",
	"Ed25519VerificationKey2018",
	"Ed25519VerificationKey2020",
	"Bls12381G1Key2020",
	"Bls12381G2Key2020",
	"PgpVerificationKey2021",
	"RsaVerificationKey2018",
	"X25519KeyAgreementKey2019",
	"SchnorrSecp256k1VerificationKey2019",
	"EcdsaSecp256k1RecoveryMethod2020",
	"VerifiableCondition2021",
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
