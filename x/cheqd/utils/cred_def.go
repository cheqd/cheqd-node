package utils

var AllowedCredDefSignatureType = []string{"CL-Sig-Cred_def"}

func IsNotCredDefSignatureType(signatureType string) bool {
	return !StringArrayContains(AllowedCredDefSignatureType, signatureType)
}
