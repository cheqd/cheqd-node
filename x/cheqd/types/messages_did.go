package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateDid{}

func NewMsgCreateDid(
	id string,
	controller []string,
	verificationMethod []*VerificationMethod,
	authentication []*types.Any,
	assertionMethod []*types.Any,
	capabilityInvocation []*types.Any,
	capabilityDelegation []*types.Any,
	keyAgreement []*types.Any,
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
	authentication []*types.Any,
	assertionMethod []*types.Any,
	capabilityInvocation []*types.Any,
	capabilityDelegation []*types.Any,
	keyAgreement []*types.Any,
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
