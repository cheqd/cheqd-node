package ante

import (
	"strings"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgCreateDidDoc int = iota
	MsgUpdateDidDoc
	MsgDeactivateDidDoc
	MsgCreateResourceDefault
	MsgCreateResourceImage
	MsgCreateResourceJson

	TaxableMsgFeeCount
)

const (
	BurnFactorDid int = iota
	BurnFactorResource

	BurnFactorCount
)

type TaxableMsgFee = [TaxableMsgFeeCount]sdk.Coins

type BurnFactor = [BurnFactorCount]sdk.Dec

var TaxableMsgFees = TaxableMsgFee{
	MsgCreateDidDoc:          (sdk.Coins)(nil),
	MsgUpdateDidDoc:          (sdk.Coins)(nil),
	MsgDeactivateDidDoc:      (sdk.Coins)(nil),
	MsgCreateResourceDefault: (sdk.Coins)(nil),
	MsgCreateResourceImage:   (sdk.Coins)(nil),
	MsgCreateResourceJson:    (sdk.Coins)(nil),
}

var BurnFactors = BurnFactor{
	BurnFactorDid:      sdk.NewDec(0),
	BurnFactorResource: sdk.NewDec(0),
}

func GetTaxableMsg(msg interface{}) bool {
	switch msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		return true
	case *didtypes.MsgUpdateDidDoc:
		return true
	case *didtypes.MsgDeactivateDidDoc:
		return true
	case *resourcetypes.MsgCreateResource:
		return true
	default:
		return false
	}
}

func GetTaxableMsgFeeWithBurnPortion(ctx sdk.Context, msg interface{}) (sdk.Coins, sdk.Coins, bool) {
	switch msg := msg.(type) {
	case *didtypes.MsgCreateDidDoc:
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], TaxableMsgFees[MsgCreateDidDoc])
		return GetRewardPortion(TaxableMsgFees[MsgCreateDidDoc], burnPortion), burnPortion, true
	case *didtypes.MsgUpdateDidDoc:
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], TaxableMsgFees[MsgUpdateDidDoc])
		return GetRewardPortion(TaxableMsgFees[MsgUpdateDidDoc], burnPortion), burnPortion, true
	case *didtypes.MsgDeactivateDidDoc:
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorDid], TaxableMsgFees[MsgDeactivateDidDoc])
		return GetRewardPortion(TaxableMsgFees[MsgDeactivateDidDoc], burnPortion), burnPortion, true
	case *resourcetypes.MsgCreateResource:
		return GetResourceTaxableMsgFee(ctx, msg)
	default:
		return nil, nil, false
	}
}

func GetRewardPortion(total sdk.Coins, burnPortion sdk.Coins) sdk.Coins {
	if burnPortion.IsZero() {
		return total
	}

	return total.Sub(burnPortion...)
}

func GetResourceTaxableMsgFee(ctx sdk.Context, msg *resourcetypes.MsgCreateResource) (sdk.Coins, sdk.Coins, bool) {
	mediaType := resourceutils.DetectMediaType(msg.GetPayload().ToResource().Resource.Data)

	// Mime type image
	if strings.HasPrefix(mediaType, "image/") {
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceImage])
		return GetRewardPortion(TaxableMsgFees[MsgCreateResourceImage], burnPortion), burnPortion, true
	}

	// Mime type json
	if strings.HasPrefix(mediaType, "application/json") {
		burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceJson])
		return GetRewardPortion(TaxableMsgFees[MsgCreateResourceJson], burnPortion), burnPortion, true
	}

	// Default mime type
	burnPortion := GetBurnFeePortion(BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceDefault])
	return GetRewardPortion(TaxableMsgFees[MsgCreateResourceDefault], burnPortion), burnPortion, true
}

func checkFeeParamsFromSubspace(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper) bool {
	didParams := didKeeper.GetParams(ctx)
	TaxableMsgFees[MsgCreateDidDoc] = sdk.Coins{sdk.Coin{Denom: didParams.TxTypes[didtypes.DefaultKeyCreateDid].Denom, Amount: didParams.TxTypes[didtypes.DefaultKeyCreateDid].Amount}}
	TaxableMsgFees[MsgUpdateDidDoc] = sdk.Coins{sdk.Coin{Denom: didParams.TxTypes[didtypes.DefaultKeyUpdateDid].Denom, Amount: didParams.TxTypes[didtypes.DefaultKeyUpdateDid].Amount}}
	TaxableMsgFees[MsgDeactivateDidDoc] = sdk.Coins{sdk.Coin{Denom: didParams.TxTypes[didtypes.DefaultKeyDeactivateDid].Denom, Amount: didParams.TxTypes[didtypes.DefaultKeyDeactivateDid].Amount}}

	resourceParams := resourceKeeper.GetParams(ctx)
	TaxableMsgFees[MsgCreateResourceImage] = sdk.Coins{sdk.Coin{Denom: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResourceImage].Denom, Amount: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResourceImage].Amount}}
	TaxableMsgFees[MsgCreateResourceJson] = sdk.Coins{sdk.Coin{Denom: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResourceJson].Denom, Amount: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResourceJson].Amount}}
	TaxableMsgFees[MsgCreateResourceDefault] = sdk.Coins{sdk.Coin{Denom: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResource].Denom, Amount: resourceParams.MediaTypes[resourcetypes.DefaultKeyCreateResource].Amount}}

	BurnFactors[BurnFactorDid] = didParams.BurnFactor
	BurnFactors[BurnFactorResource] = resourceParams.BurnFactor

	return true
}

func IsTaxableTx(ctx sdk.Context, didKeeper DidKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx) (bool, sdk.Coins, sdk.Coins) {
	_ = checkFeeParamsFromSubspace(ctx, didKeeper, resourceKeeper)
	reward := (sdk.Coins)(nil)
	burn := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		rewardPortion, burnPortion, isIdentityMsg := GetTaxableMsgFeeWithBurnPortion(ctx, msg)
		if !isIdentityMsg {
			continue
		}
		if rewardPortion != nil {
			reward = reward.Add(rewardPortion...)
			burn = burn.Add(burnPortion...)
		}
	}

	if !reward.IsZero() {
		return true, reward, burn
	}

	return false, nil, nil
}

func IsTaxableTxLite(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		if GetTaxableMsg(msg) {
			return true
		}
	}

	return false
}
