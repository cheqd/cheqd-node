package ante

import (
	"strings"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourceutils "github.com/cheqd/cheqd-node/x/resource/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgCreateDid int = iota
	MsgUpdateDid
	MsgDeactivateDid
	MsgCreateResourceDefault
	MsgCreateResourceImage
	MsgCreateResourceJson

	TaxableMsgFeeCount
)

const (
	BurnFactorCheqd int = iota
	BurnFactorResource

	BurnFactorCount
)

type TaxableMsgFee = [TaxableMsgFeeCount]sdk.Coins

type BurnFactor = [BurnFactorCount]sdk.Dec

var TaxableMsgFees = TaxableMsgFee{
	MsgCreateDid:      			(sdk.Coins)(nil),
	MsgUpdateDid:      			(sdk.Coins)(nil),
	MsgDeactivateDid:  			(sdk.Coins)(nil),
	MsgCreateResourceDefault: 	(sdk.Coins)(nil),
	MsgCreateResourceImage: 	(sdk.Coins)(nil),
	MsgCreateResourceJson: 		(sdk.Coins)(nil),
}

var BurnFactors = BurnFactor{
	BurnFactorCheqd:    sdk.NewDec(0),
	BurnFactorResource: sdk.NewDec(0),
}

func GetTaxableMsg(msg interface{}) bool {
	switch msg.(type) {
	case *cheqdtypes.MsgCreateDid:
		return true
	case *cheqdtypes.MsgUpdateDid:
		return true
	// TODO: Add `MsgDeactivateDid` when it will be implemented
	case *resourcetypes.MsgCreateResource:
		return true
	default:
		return false
	}
}

func GetTaxableMsgFeeWithBurnPortion(ctx sdk.Context, msg interface{}) (sdk.Coins, sdk.Coins, bool) {
	switch msg := msg.(type) {
	case *cheqdtypes.MsgCreateDid:
		return TaxableMsgFees[MsgCreateDid], GetBurnFeePortion(ctx, BurnFactors[BurnFactorCheqd], TaxableMsgFees[MsgCreateDid]), true
	case *cheqdtypes.MsgUpdateDid:
		return TaxableMsgFees[MsgUpdateDid], GetBurnFeePortion(ctx, BurnFactors[BurnFactorCheqd], TaxableMsgFees[MsgUpdateDid]), true
	// TODO: Add `MsgDeactivateDid` when it will be implemented
	case *resourcetypes.MsgCreateResource:
		return GetResourceTaxableMsgFee(ctx, msg)
	default:
		return nil, nil, false
	}
}

func GetResourceTaxableMsgFee(ctx sdk.Context, msg *resourcetypes.MsgCreateResource) (sdk.Coins, sdk.Coins, bool) {
	mediaType := resourceutils.DetectMediaType(msg.GetPayload().ToResource().Data)

	// Mime type image
	if strings.HasPrefix(mediaType, "image/") {
		return TaxableMsgFees[MsgCreateResourceImage], GetBurnFeePortion(ctx, BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceImage]), true
	}

	// Mime type json
	if strings.HasPrefix(mediaType, "application/json") {
		return TaxableMsgFees[MsgCreateResourceJson], GetBurnFeePortion(ctx, BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceJson]), true
	}

	// Default mime type
	return TaxableMsgFees[MsgCreateResourceDefault], GetBurnFeePortion(ctx, BurnFactors[BurnFactorResource], TaxableMsgFees[MsgCreateResourceDefault]), true
}

func checkFeeParamsFromState(ctx sdk.Context, cheqdKeeper CheqdKeeper, resourceKeeper ResourceKeeper) (bool) {
	cheqdParams := cheqdKeeper.GetParams(ctx)
	createDidFeeCoins := sdk.NewCoins(cheqdParams.CreateDid)
	updateDidFeeCoins := sdk.NewCoins(cheqdParams.UpdateDid)
	deactivateDidFeeCoins := sdk.NewCoins(cheqdParams.DeactivateDid)

	resourceParams := resourceKeeper.GetParams(ctx)
	createResourceImageFeeCoins := sdk.NewCoins(resourceParams.Image)
	createResourceJsonFeeCoins := sdk.NewCoins(resourceParams.Json)
	createResourceDefaultFeeCoins := sdk.NewCoins(resourceParams.Default)

	if !createDidFeeCoins.IsEqual(TaxableMsgFees[MsgCreateDid]) {
		TaxableMsgFees[MsgCreateDid] = createDidFeeCoins
	}

	if !updateDidFeeCoins.IsEqual(TaxableMsgFees[MsgUpdateDid]) {
		TaxableMsgFees[MsgUpdateDid] = updateDidFeeCoins
	}

	if !deactivateDidFeeCoins.IsEqual(TaxableMsgFees[MsgDeactivateDid]) {
		TaxableMsgFees[MsgDeactivateDid] = deactivateDidFeeCoins
	}

	if !cheqdParams.BurnFactor.Equal(BurnFactors[BurnFactorCheqd]) {
		BurnFactors[BurnFactorCheqd] = cheqdParams.BurnFactor
	}

	if !createResourceImageFeeCoins.IsEqual(TaxableMsgFees[MsgCreateResourceImage]) {
		TaxableMsgFees[MsgCreateResourceImage] = createResourceImageFeeCoins
	}

	if !createResourceJsonFeeCoins.IsEqual(TaxableMsgFees[MsgCreateResourceJson]) {
		TaxableMsgFees[MsgCreateResourceJson] = createResourceJsonFeeCoins
	}

	if !createResourceDefaultFeeCoins.IsEqual(TaxableMsgFees[MsgCreateResourceDefault]) {
		TaxableMsgFees[MsgCreateResourceDefault] = createResourceDefaultFeeCoins
	}

	if !resourceParams.BurnFactor.Equal(BurnFactors[BurnFactorResource]) {
		BurnFactors[BurnFactorResource] = resourceParams.BurnFactor
	}

	return true
}


func IsIdentityTx(ctx sdk.Context, cheqdKeeper CheqdKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx) (bool, sdk.Coins, sdk.Coins) {
	_ = checkFeeParamsFromState(ctx, cheqdKeeper, resourceKeeper)
	fee := (sdk.Coins)(nil)
	burn := (sdk.Coins)(nil)
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		identityMsgFee, burnPortion, isIdentityMsg := GetTaxableMsgFeeWithBurnPortion(ctx, msg)
		if !isIdentityMsg {
			continue
		}
		if identityMsgFee != nil {
			fee = fee.Add(identityMsgFee...)
			burn = burn.Add(burnPortion...)
		}
	}

	if !fee.IsZero() {
		return true, fee, burn
	}

	return false, nil, nil
}

func IsIdentityTxLite(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		if GetTaxableMsg(msg) {
			return true
		}
	}

	return false
}
