package validate

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

func (msg *MsgUpdateDidPayload) Validate(namespace string) error {
	if !utils.IsValidDid(namespace, msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if notValid, i := utils.IsNotValidDIDArray(namespace, msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s at position %d", msg.Controller[i], i)
	}

	if err := ValidateVerificationMethods(namespace, msg.Id, msg.VerificationMethod); err != nil {
		return err
	}

	if err := ValidateServices(namespace, msg.Id, msg.Service); err != nil {
		return err
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.Authentication); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("Authentication item %s", msg.Authentication[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityInvocation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityInvocation item %s", msg.CapabilityInvocation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.CapabilityDelegation); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("CapabilityDelegation item %s", msg.CapabilityDelegation[i])
	}

	if notValid, i := utils.IsNotValidDIDArrayFragment(namespace, msg.KeyAgreement); notValid {
		return ErrBadRequestIsNotDidFragment.Wrapf("KeyAgreement item %s", msg.KeyAgreement[i])
	}

	if len(msg.Authentication) == 0 && len(msg.Controller) == 0 {
		return ErrBadRequest.Wrap("The message must contain either a Controller or a Authentication")
	}

	for _, i := range msg.Authentication {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.KeyAgreement {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityDelegation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	for _, i := range msg.CapabilityInvocation {
		if !IncludeVerificationMethod(msg.Id, msg.VerificationMethod, i) {
			return ErrVerificationMethodNotFound.Wrap(i)
		}
	}

	return nil
}
