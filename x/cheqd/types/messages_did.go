package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*DidService,
) *MsgCreateDid {
	return &MsgCreateDid{
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		AlsoKnownAs:          alsoKnownAs,
		Service:              service,
	}
}

/*
func (msg *MsgCreateDid) Verify(keeper *keeper.Keeper, ctx *sdk.Context, request *MsgWriteRequest) (bool, error) {
	signingInput, err:=utils.BuildSigningInput(request)
	if err!=nil {
		return false, err
	}

	// if controller is present
	if len(msg.Controller) > 0 {
		for _, controller := range msg.Controller {
			var authentication []string
			var verificationMethod []*VerificationMethod

			// if self-signed
			if controller == msg.Id {
				authentication=msg.Authentication
				verificationMethod=msg.VerificationMethod
			} else {
				didDoc, _, err := keeper.GetDid(ctx, controller)
				if err != nil {
					return false, ErrDidDocNotFound.Wrap(controller)
				}

				authentication=didDoc.Authentication
				verificationMethod=didDoc.VerificationMethod
			}

			// check all controller signatures
			return utils.VerifyIdentitySignature(controller, authentication, verificationMethod, request.Signatures, signingInput)
		}
	}

	// controller is not present but there are authentications
	if len(msg.Authentication) > 0 {
		return utils.VerifyIdentitySignature(msg.Id, msg.Authentication, msg.VerificationMethod, request.Signatures, signingInput)
	}

	return false, ErrInvalidSignature.Wrap("At least DID Doc should contain `controller` or `authentication`")
}*/

func (msg *MsgCreateDid) Route() string {
	return RouterKey
}

func (msg *MsgCreateDid) Type() string {
	return "CreateDid"
}

func (msg *MsgCreateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgCreateDid) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgCreateDid) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgUpdateDid{}

func NewMsgUpdateDid(
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*DidService,
) *MsgUpdateDid {
	return &MsgUpdateDid{
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		AlsoKnownAs:          alsoKnownAs,
		Service:              service,
	}
}

func (msg *MsgUpdateDid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDid) Type() string {
	return "UpdateDid"
}

func (msg *MsgUpdateDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateDid) GetSignBytes() []byte {
	return []byte{}
}

func (msg *MsgUpdateDid) ValidateBasic() error {
	return nil
}
