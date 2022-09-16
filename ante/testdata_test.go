package ante_test

import (
	"crypto/ed25519"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	cheqdtests "github.com/cheqd/cheqd-node/x/cheqd/tests"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetests "github.com/cheqd/cheqd-node/x/resource/tests"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTestFeeAmount is a test fee amount in `cheq`.
func NewTestFeeAmount() sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("cheq", 150))
}

// NewTestFeeAmount is a test fee amount lower than the fixed fee in `ncheq`.
func NewTestFeeAmountMinimalDenomLTFixedFee() sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("ncheq", 150*1e9))
}

// NewTestFeeAmount is a test fee amount 5x greater than the fixed fee in `ncheq`.
func NewTestFeeAmountMinimalDenomEFixedFee() sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("ncheq", cheqdante.MinimalIdentityFee*cheqdante.CheqFactor))
}

// NewTestFeeAmount is a test fee amount 5x greater than the fixed fee in `ncheq`.
func NewTestFeeAmountMinimalDenomGTFixedFee() sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("ncheq", 5*cheqdante.MinimalIdentityFee*cheqdante.CheqFactor))
}

func NewTestDidMsg() *cheqdtypes.MsgCreateDid {
	payload := &cheqdtypes.MsgCreateDidPayload{
		Id:             cheqdtests.ImposterDID,
		Authentication: []string{cheqdtests.ImposterKey1},
		VerificationMethod: []*cheqdtypes.VerificationMethod{
			{
				Id:         cheqdtests.ImposterKey1,
				Type:       cheqdtests.Ed25519VerificationKey2020,
				Controller: cheqdtests.ImposterDID,
			},
		},
	}
	signInput := &cheqdtypes.SignInfo{
		VerificationMethodId: cheqdtests.ImposterKey1,
		Signature:            string(ed25519.Sign(cheqdtests.GenerateKeyPair().PrivateKey, payload.GetSignBytes())),
	}
	return &cheqdtypes.MsgCreateDid{
		Payload:    payload,
		Signatures: []*cheqdtypes.SignInfo{signInput},
	}
}

func NewTestResourceMsg() *resourcetypes.MsgCreateResource {
	payload := &resourcetypes.MsgCreateResourcePayload{
		CollectionId: cheqdtests.ImposterDID,
		Id:           resourcetests.ResourceId,
		Name:         resourcetests.TestResourceName,
		ResourceType: resourcetests.CLSchemaType,
		Data:         []byte(resourcetests.SchemaData),
	}
	signInput := &cheqdtypes.SignInfo{
		VerificationMethodId: cheqdtests.ImposterKey1,
		Signature:            string(ed25519.Sign(cheqdtests.GenerateKeyPair().PrivateKey, payload.GetSignBytes())),
	}
	return &resourcetypes.MsgCreateResource{
		Payload:    payload,
		Signatures: []*cheqdtypes.SignInfo{signInput},
	}
}
