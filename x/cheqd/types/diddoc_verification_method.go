package types

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/multiformats/go-multibase"
)

var SupportedMethodTypes = []string{
	JsonWebKey2020{}.Type(),
	Ed25519VerificationKey2020{}.Type(),
}

func NewVerificationMethod(id string, type_ string, controller string, verificationMaterial string) *VerificationMethod {
	return &VerificationMethod{
		Id:                   id,
		Type:                 type_,
		Controller:           controller,
		VerificationMaterial: verificationMaterial,
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

	switch vm.Type {
	case Ed25519VerificationKey2020{}.Type():
		var ed25519VerificationKey2020 Ed25519VerificationKey2020
		err := json.Unmarshal([]byte(vm.VerificationMaterial), &ed25519VerificationKey2020)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to unmarshal verification material for %s", vm.Id)
		}

		_, keyBytes, err := multibase.Decode(ed25519VerificationKey2020.PublicKeyMultibase)
		if err != nil {
			return err
		}

		verificationError = utils.VerifyED25519Signature(keyBytes, message, signature)

	case JsonWebKey2020{}.Type():
		var jsonWebKey2020 JsonWebKey2020
		err := json.Unmarshal([]byte(vm.VerificationMaterial), &jsonWebKey2020)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to unmarshal verification material for %s", vm.Id)
		}

		var raw interface{}
		err = jwk.ParseRawKey([]byte(jsonWebKey2020.PublicKeyJwk), &raw)
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
	did, path, query, fragment := utils.MustSplitDIDUrl(vm.Id)
	if did == old {
		did = new
	}

	vm.Id = utils.JoinDIDUrl(did, path, query, fragment)
}

// Validation

func (vm VerificationMethod) Validate(baseDid string, allowedNamespaces []string) error {
	return validation.ValidateStruct(&vm,
		validation.Field(&vm.Id, validation.Required, IsDIDUrl(allowedNamespaces, Empty, Empty, Required), HasPrefix(baseDid)),
		validation.Field(&vm.Controller, validation.Required, IsDID(allowedNamespaces)),
		validation.Field(&vm.Type, validation.Required, validation.In(utils.ToInterfaces(SupportedMethodTypes)...)),
		validation.Field(&vm.VerificationMaterial,
			validation.When(vm.Type == Ed25519VerificationKey2020{}.Type(), validation.Required, ValidEd25519VerificationKey2020Rule()),
		),
		validation.Field(&vm.VerificationMaterial,
			validation.When(vm.Type == JsonWebKey2020{}.Type(), validation.Required, ValidJsonWebKey2020Rule()),
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

func IsUniqueVerificationMethodListByIdRule() *CustomErrorRule {
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
