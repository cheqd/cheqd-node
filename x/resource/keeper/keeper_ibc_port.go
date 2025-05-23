package keeper

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
)

// GetPort returns the portID for the resource module. Used in ExportGenesis
func (k Keeper) GetPort(ctx context.Context) (string, error) {
	return k.Port.Get(ctx)
}

// SetPort sets the portID for the resource module. Used in InitGenesis
func (k Keeper) SetPort(ctx context.Context, portID string) error {
	return k.Port.Set(ctx, portID)
}

// IsBound checks if the  module is already bound to the desired port
func (k Keeper) IsBound(ctx context.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(sdktypes.UnwrapSDKContext(ctx), host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to the module's InitGenesis function
func (k Keeper) BindPort(ctx context.Context, portID string) error {
	cap := k.portKeeper.BindPort(sdktypes.UnwrapSDKContext(ctx), portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// ClaimCapability allows the resource module to claim a capability that IBC module passes to it
func (k Keeper) ClaimCapability(ctx context.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(sdktypes.UnwrapSDKContext(ctx), cap, name)
}

// AuthenticateCapability attempts to authenticate a given capability and name
// from a caller. It allows for a caller to check that a capability does in fact
// correspond to a particular name.
func (k Keeper) AuthenticateCapability(ctx context.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(sdktypes.UnwrapSDKContext(ctx), cap, name)
}
