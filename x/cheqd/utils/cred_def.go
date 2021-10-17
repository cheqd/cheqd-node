package utils

import "github.com/cheqd/cheqd-node/x/cheqd/utils/strings"

var AllowedCredDefSignatureType = []string{"CL-Sig-Cred_def"}

func IsNotCredDefSignatureType(signatureType string) bool {
	return !strings.Include(AllowedCredDefSignatureType, signatureType)
}
