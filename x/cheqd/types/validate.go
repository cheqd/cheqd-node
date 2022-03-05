package types

import "github.com/cheqd/cheqd-node/x/cheqd/utils"


type DIDRule struct {
	method string
	allowedNamespaces []string
}

func NewDIDRule(method string, allowedNamespaces []string) *DIDRule {
	return &DIDRule{method: method, allowedNamespaces: allowedNamespaces}
}

func (D DIDRule) Validate(value interface{}) error {
	casted, ok := value.(string)
	if !ok {
		panic("DIDRule must be only applied on string properties")
	}

	return utils.ValidateDID(casted, D.method, D.allowedNamespaces)
}


type DIDUrlRule struct {
	method string
	allowedNamespaces []string
}

func NewDIDUrlRule(method string, allowedNamespaces []string) *DIDUrlRule {
	return &DIDUrlRule{method: method, allowedNamespaces: allowedNamespaces}
}

func (D DIDUrlRule) Validate(value interface{}) error {
	casted, ok := value.(string)
	if !ok {
		panic("DIDUrlRule must be only applied on string properties")
	}

	return utils.ValidateDIDUrl(casted, D.method, D.allowedNamespaces)
}


////const (
////	PublicKeyJwk       = "PublicKeyJwk"
////	PublicKeyMultibase = "PublicKeyMultibase"
////)
////
////var VerificationMethodType = map[string]string{
////	"JsonWebKey2020":             PublicKeyJwk,
////	"Ed25519VerificationKey2020": PublicKeyMultibase,
////}
////
////var ServiceType = []string{
////	"LinkedDomains",
////	"DIDCommMessaging",
////}
////
////func GetVerificationMethodType(vmType string) string {
////	return VerificationMethodType[vmType]
////}
////
////func IsValidDidServiceType(sType string) bool {
////	return strings.Contains(ServiceType, sType)
////}
