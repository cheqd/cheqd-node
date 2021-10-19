package utils

import (
	ustring "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"
	"regexp"
	"strings"
)

var DidForbiddenSymbolsRegexp, _ = regexp.Compile("^[^#?&/\\\\]+$")

var VerificationMethodType = []string{
	"JsonWebKey2020",
	"EcdsaSecp256k1VerificationKey2019",
	"Ed25519VerificationKey2018",
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

func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
	fragments := strings.Split(didUrl, "#")
	return fragments[0], fragments[1]
}

func IsDidFragment(didUrl string) bool {
	if !strings.Contains(didUrl, "#") {
		return false
	}

	if didUrl[0] == '#' {
		return true
	}

	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
	return IsDid(did)
}

func IsFullDidFragment(didUrl string) bool {
	if !strings.Contains(didUrl, "#") {
		return false
	}

	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
	return IsDid(did)
}

func IsVerificationMethodType(vmType string) bool {
	return ustring.Include(VerificationMethodType, vmType)
}

func IsDidServiceType(sType string) bool {
	return ustring.Include(ServiceType, sType)
}

func ResolveId(did string, methodId string) string {
	result := methodId

	methodDid, methodFragment := SplitDidUrlIntoDidAndFragment(methodId)
	if len(methodDid) == 0 {
		result = did + "#" + methodFragment
	}

	return result
}

func ArrayContainsNotDid(array []string) (bool, int) {
	for i, did := range array {
		if IsNotDid(did) {
			return true, i
		}
	}

	return false, 0
}

func ArrayContainsNotDidFragment(array []string) (bool, int) {
	for i, did := range array {
		if !IsDidFragment(did) {
			return true, i
		}
	}

	return false, 0
}

func IsNotDid(did string) bool {
	return IsDid(did) == false
}

func IsDid(did string) bool {
	if len(did) == 0 {
		return false
	}

	if !DidForbiddenSymbolsRegexp.MatchString(did) {
		return false
	}

	fragments := strings.Split(did, ":")

	if len(fragments) <= 3 {
		return false
	}

	if fragments[0] != "did" {
		return false
	}

	if fragments[1] != MethodName {
		return false
	}

	if fragments[2] != MethodSpecificId {
		return false
	}

	return true
}
