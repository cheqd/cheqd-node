package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

// DID rule
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

// DIDURL rule
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

// List of DID rule
type DIDListRule struct {
	method string
	allowedNamespaces []string
}

func NewDIDListRule(method string, allowedNamespaces []string) *DIDListRule {
	return &DIDListRule{method: method, allowedNamespaces: allowedNamespaces}
}

func (D DIDListRule) Validate(value interface{}) error {
	casted, ok := value.([]string)
	if !ok {
		panic("DIDListRule must be only applied on list of strings")
	}
	return utils.ValidateDIDList(casted, D.method, D.allowedNamespaces)
}

// PublicKeyMultibase
type PublicKeyMultibaseRule struct {
}

func NewPublicKeyMultibaseRule() *PublicKeyMultibaseRule {
	return &PublicKeyMultibaseRule{}
}

func (D PublicKeyMultibaseRule) Validate(value interface{}) error {
	casted, ok := value.(string)
	if !ok {
		panic("PublicKeyMultibaseRule must be only applied on string properties")
	}

	return utils.ValidatePublicKeyMultibase(casted)
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
