package keeper

import (
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// GetPort returns the portID for the resource module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	byteKey := didutils.StrBytes(types.ResourcePortIdKey)
	return string(store.Get(byteKey))
}

// SetPort sets the portID for the resource module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	byteKey := didutils.StrBytes(types.ResourcePortIdKey)
	store.Set(byteKey, []byte(portID))
}

// IsBound checks if the  module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to the module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// ClaimCapability allows the resource module to claim a capability that IBC module passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// AuthenticateCapability attempts to authenticate a given capability and name
// from a caller. It allows for a caller to check that a capability does in fact
// correspond to a particular name.
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}
