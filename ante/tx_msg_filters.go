package ante

import (
	"cosmossdk.io/math"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgCreateDid int = iota
	MsgUpdateDid
	MsgCreateResource

	TaxableMsgFeeCount
)

type TaxableMsgFee = [TaxableMsgFeeCount]sdk.Coins

var (
	TaxableMsgFees = TaxableMsgFee{
		MsgCreateDid:      sdk.NewCoins(sdk.Coin{Denom: "ncheq", Amount: math.NewInt(int64(MinimalIdentityFee * CheqFactor))}),
		MsgUpdateDid:      sdk.NewCoins(sdk.Coin{Denom: "ncheq", Amount: math.NewInt(int64(MinimalIdentityFee * CheqFactor))}),
		MsgCreateResource: sdk.NewCoins(sdk.Coin{Denom: "ncheq", Amount: math.NewInt(int64(MinimalIdentityFee * CheqFactor))}),
	}
)

func GetTaxableMsgFee(msg interface{}) (sdk.Coins, bool) {
	switch msg.(type) {
	case *cheqdtypes.MsgCreateDid:
		return TaxableMsgFees[MsgCreateDid], true
	case *cheqdtypes.MsgUpdateDid:
		return TaxableMsgFees[MsgUpdateDid], true
	case *resourcetypes.MsgCreateResource:
		return TaxableMsgFees[MsgCreateResource], true
	default:
		return nil, false
	}
}

func IsIdentityTx(tx sdk.Tx) (bool, sdk.Coins) {
	fee := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		identityMsgFee, isIdentityMsg := GetTaxableMsgFee(msg)
		if !isIdentityMsg {
			continue
		}
		if identityMsgFee != nil {
			fee = fee.Add(identityMsgFee...)
		}
	}

	if !fee.IsZero() {
		return true, fee
	}

	return false, nil
}
