package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var AllowedCredDefSignatureType = []string{"CL-CredDef"}

var CredDefSuffix = "?service=CL-CredDef"
var CredDefSuffixLen = len(CredDefSuffix)

func IsNotCredDefSignatureType(signatureType string) bool {
	return !strings.Include(AllowedCredDefSignatureType, signatureType)
}

func IsCredDef(did string) bool {
	return len(did) >= CredDefSuffixLen && did[len(did)-CredDefSuffixLen:] == CredDefSuffix
}

func GetDidFromCredDef(credDef string) string {
	return credDef[:len(credDef)-CredDefSuffixLen]
}

func GetCredDefFromDid(did string) string {
	return did + CredDefSuffix
}
