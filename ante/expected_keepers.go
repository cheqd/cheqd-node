package ante

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

type DidKeeper interface {
	GetParams(ctx context.Context) (didtypes.FeeParams, error)
}

type ResourceKeeper interface {
	GetParams(ctx context.Context) (params resourcetypes.FeeParams, err error)
}
type AccountKeeper interface {
	GetParams(ctx context.Context) (params authtypes.Params)
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	AddressCodec() address.Codec
}
type FeeGrantKeeper interface {
	UseGrantedFees(ctx context.Context, granter, grantee sdk.AccAddress, fee sdk.Coins, msgs []sdk.Msg) error
}

type FeeMarketKeeper interface {
	GetState(ctx context.Context) (feemarkettypes.State, error)
	GetMinGasPrice(ctx context.Context, denom string) (sdk.DecCoin, error)
	GetParams(ctx context.Context) (feemarkettypes.Params, error)
	SetState(ctx context.Context, state feemarkettypes.State) error
	SetParams(ctx context.Context, params feemarkettypes.Params) error
	ResolveToDenom(ctx context.Context, coin sdk.DecCoin, denom string) (sdk.DecCoin, error)
}

type OracleKeeper interface {
	GetEMA(ctx sdk.Context, denom string) (math.LegacyDec, bool)
	GetExchangeRate(ctx sdk.Context, denom string) (math.LegacyDec, error)
	GetWMA(ctx sdk.Context, denom string, strategy string) (math.LegacyDec, bool)
}
