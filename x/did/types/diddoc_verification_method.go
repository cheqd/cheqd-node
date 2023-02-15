package types

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/cheqd/cheqd-node/x/did/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
)

const (
	JSONWebKey2020Type             = "JsonWebKey2020"
	Ed25519VerificationKey2020Type = "Ed25519VerificationKey2020"
	Ed25519VerificationKey2018Type = "Ed25519VerificationKey2018"
)

var SupportedMethodTypes = []string{
	JSONWebKey2020Type,
	Ed25519VerificationKey2020Type,
	Ed25519VerificationKey2018Type,
}

func NewVerificationMethod(id string, vmType string, controller string, verificationMaterial string) *VerificationMethod {
	return &VerificationMethod{
		Id:                     id,
		VerificationMethodType: vmType,
		Controller:             controller,
		VerificationMaterial:   verificationMaterial,
	}
}

// Helpers

func FindVerificationMethod(vms []VerificationMethod, id string) (VerificationMethod, bool) {
	for _, vm := range vms {
		if vm.Id == id {
			return vm, true
		}
	}

	return VerificationMethod{}, false
}

func GetVerificationMethodIds(vms []*VerificationMethod) []string {
	res := make([]string, len(vms))

	for i := range vms {
		res[i] = vms[i].Id
	}

	return res
}

func VerifySignature(vm VerificationMethod, message []byte, signature []byte) error {
	var verificationError error

	switch vm.VerificationMethodType {
	case Ed25519VerificationKey2020Type:

		_, multibaseBytes, err := multibase.Decode(vm.VerificationMaterial)
		if err != nil {
			return err
		}

		keyBytes := utils.GetEd25519VerificationKey2020(multibaseBytes)
		verificationError = utils.VerifyED25519Signature(keyBytes, message, signature)

	case JSONWebKey2020Type:
		var raw interface{}
		err := jwk.ParseRawKey([]byte(vm.VerificationMaterial), &raw)
		if err != nil {
			return fmt.Errorf("can't parse jwk: %s", err.Error())
		}

		switch pubKey := raw.(type) {
		case *rsa.PublicKey:
			verificationError = utils.VerifyRSASignature(*pubKey, message, signature)
		case *ecdsa.PublicKey:
			verificationError = utils.VerifyECDSASignature(*pubKey, message, signature)
		case ed25519.PublicKey:
			verificationError = utils.VerifyED25519Signature(pubKey, message, signature)
		default:
			panic("unsupported jwk key") // This should have been checked during basic validation
		}

	case Ed25519VerificationKey2018Type:
		publicKeyBytes, err := base58.Decode(vm.VerificationMaterial)
		if err != nil {
			return err
		}

		verificationError = utils.VerifyED25519Signature(publicKeyBytes, message, signature)

	default:
		panic("unsupported verification method type") // This should have also been checked during basic validation
	}

	if verificationError != nil {
		return ErrInvalidSignature.Wrapf("verification method: %s, err: %s", vm.Id, verificationError.Error())
	}

	return nil
}

func VerificationMethodListToMapByFragment(vms []*VerificationMethod) map[string]VerificationMethod {
	result := map[string]VerificationMethod{}

	for _, vm := range vms {
		_, _, _, fragment := utils.MustSplitDIDUrl(vm.Id)
		result[fragment] = *vm
	}

	return result
}

// ReplaceDids replaces ids in all fields
func (vm *VerificationMethod) ReplaceDids(old, new string) {
	// Controller
	if vm.Controller == old {
		vm.Controller = new
	}

	// Id
	vm.Id = utils.ReplaceDidInDidURL(vm.Id, old, new)
}

// Validation
func (vm VerificationMethod) Validate(baseDid string, allowedNamespaces []string) error {
	return validation.ValidateStruct(&vm,
		validation.Field(&vm.Id, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(baseDid)),
		validation.Field(&vm.Controller, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&vm.VerificationMethodType, validation.Required, validation.In(utils.ToInterfaces(SupportedMethodTypes)...)),
		validation.Field(&vm.VerificationMaterial,
			validation.When(vm.VerificationMethodType == Ed25519VerificationKey2020Type, validation.Required, IsMultibaseEd25519VerificationKey2020()),
		),
		validation.Field(&vm.VerificationMaterial,
			validation.When(vm.VerificationMethodType == Ed25519VerificationKey2018Type, validation.Required, IsBase58Ed25519VerificationKey2018()),
		),
		validation.Field(&vm.VerificationMaterial,
			validation.When(vm.VerificationMethodType == JSONWebKey2020Type, validation.Required, IsJWK()),
		),
	)
}

func ValidVerificationMethodRule(baseDid string, allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(VerificationMethod)
		if !ok {
			panic("ValidVerificationMethodRule must be only applied on verification methods")
		}

		return casted.Validate(baseDid, allowedNamespaces)
	})
}

func IsUniqueVerificationMethodListByIDRule() *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.([]*VerificationMethod)
		if !ok {
			panic("IsUniqueVerificationMethodListByIdRule must be only applied on VM lists")
		}

		ids := GetVerificationMethodIds(casted)
		if !utils.IsUnique(ids) {
			return errors.New("there are verification method duplicates")
		}

		return nil
	})
}
