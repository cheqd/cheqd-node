package types

import "github.com/multiformats/go-multibase"

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateDidPayload) ToDID() Did {
	return Did{
		Context:              msg.Context,
		Id:                   msg.Id,
		Controller:           msg.Controller,
		VerificationMethod:   msg.VerificationMethod,
		Authentication:       msg.Authentication,
		AssertionMethod:      msg.AssertionMethod,
		CapabilityInvocation: msg.CapabilityInvocation,
		CapabilityDelegation: msg.CapabilityDelegation,
		KeyAgreement:         msg.KeyAgreement,
		AlsoKnownAs:          msg.AlsoKnownAs,
		Service:              msg.Service,
	}
}

// AggregateControllerDIDs returns all controller DIDs used in the message
func (msg *MsgCreateDidPayload) AggregateControllerDIDs() []string {
	result := msg.Controller

	for _, vm := range msg.VerificationMethod {
		result = append(result, vm.Controller)
	}




	if len(v.PublicKeyMultibase) > 0 {
		_, key, err := multibase.Decode(v.PublicKeyMultibase)
		if err != nil {
			return nil, ErrInvalidPublicKey.Wrapf("Cannot decode verification method '%s' public key", v.Id)
		}
		return key, nil
	}

	if len(v.PublicKeyJwk) > 0 {
		return nil, ErrInvalidPublicKey.Wrap("JWK format not supported")
	}

	return nil, ErrInvalidPublicKey.Wrapf("verification method '%s' public key not found", v.Id)
}
