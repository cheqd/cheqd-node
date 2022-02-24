package types

import "github.com/cheqd/cheqd-node/x/cheqd/utils"

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, controller := range msg.Controller {
			if controller == msg.Id {
				result[i] = Signer{
					Signer:             controller,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: controller,
				}
			}
		}

		return result
	}

	if len(msg.Authentication) > 0 {
		return []Signer{
			{
				Signer:             msg.Id,
				Authentication:     msg.Authentication,
				VerificationMethod: msg.VerificationMethod,
			},
		}
	}

	return []Signer{}
}

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateDidPayload) ValidateDynamic(namespace string) error {
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
