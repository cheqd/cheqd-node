package ante_tests

import (
	"crypto/ed25519"

	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewTestFeeAmount() sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(didtypes.BaseMinimalDenom, 150*didtypes.BaseFactor))
}

func SandboxDidDoc() *didtypes.MsgCreateDidDoc {
	setup := didtestssetup.Setup()
	didDocInfo := setup.BuildSimpleDidDoc()

	signBytes := didDocInfo.Msg.GetSignBytes()
	var signatures []*didtypes.SignInfo

	signatures = append(signatures, &didtypes.SignInfo{
		VerificationMethodId: didDocInfo.SignInput.VerificationMethodId,
		Signature:            ed25519.Sign(didDocInfo.SignInput.Key, signBytes),
	})

	return &didtypes.MsgCreateDidDoc{
		Payload:    didDocInfo.Msg,
		Signatures: signatures,
	}
}

func SandboxResource() *resourcetypes.MsgCreateResource {
	setup := resourcetestssetup.Setup()
	didDocInfo := setup.BuildSimpleDidDoc()
	resource := setup.BuildSimpleResource(didDocInfo.Did, `{"message": "test"}`, "Test Name", "Test Type")

	signBytes := resource.GetSignBytes()
	var signatures []*didtypes.SignInfo

	signatures = append(signatures, &didtypes.SignInfo{
		VerificationMethodId: didDocInfo.SignInput.VerificationMethodId,
		Signature:            ed25519.Sign(didDocInfo.SignInput.Key, signBytes),
	})

	return &resourcetypes.MsgCreateResource{
		Payload:    &resource,
		Signatures: signatures,
	}
}
