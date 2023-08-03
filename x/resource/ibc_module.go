package resource

import (
	"encoding/json"

	resourceKeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"
)

type IBCModule struct {
	keeper resourceKeeper.Keeper
}

func NewIBCModule(k resourceKeeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

func validateChannelParams(
	ctx sdk.Context,
	keeper resourceKeeper.Keeper,
	order channeltypes.Order,
	version string,
) error {
	if order != channeltypes.UNORDERED {
		return sdkerrors.Wrapf(channeltypes.ErrInvalidChannelOrdering, "expected %s channel, got %s ", channeltypes.UNORDERED, order)
	}

	if version != resourcetypes.IBCVersion {
		return sdkerrors.Wrapf(resourcetypes.ErrInvalidVersion, "expected version %s , got %s ", resourcetypes.IBCVersion, version)
	}

	return nil
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	_ []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	if err := validateChannelParams(ctx, im.keeper, order, counterpartyVersion); err != nil {
		return "", err
	}

	// Claim channel capability passed back by IBC module
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	return counterpartyVersion, nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	// Require portID is the portID module is bound to
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if err := validateChannelParams(ctx, im.keeper, order, counterpartyVersion); err != nil {
		return "", err
	}

	// Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// If the module can already authenticate the capability then the module already owns it so we don't need to claim
	// Otherwise, the module does not have channel capability and we must claim it from IBC
	if !im.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		// Only claim channel capability passed back by IBC module if we do not already own it
		if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
			return "", err
		}
	}

	return counterpartyVersion, nil

}

func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	_ string,
	counterpartyVersion string,
) error {
	if counterpartyVersion != resourcetypes.IBCVersion {
		return sdkerrors.Wrapf(resourcetypes.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, resourcetypes.IBCVersion)
	}
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for transfer channels
	// Todo: what should the appropriate action be?
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface.
// A successful acknowledgement in the form of the resource is returned if
// the packet data is successfully decoded and the receive application logic returns without error.
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {

	var reqPacket resourcetypes.ResourceReqPacket

	err := json.Unmarshal(packet.GetData(), &reqPacket)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			resourcetypes.ErrUnexpectedPacket.Wrapf("Error unmarshal packet data: %s", err),
		)
	}
	resource, err := im.keeper.GetResource(&ctx, reqPacket.CollectionId, reqPacket.ResourceId)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			resourcetypes.ErrResourceNotAvail.Wrapf("Error get resource: %s", err),
		)
	}

	jsonResource, err := json.Marshal(resource)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			resourcetypes.ErrInternal.Wrapf("Error marshal ack json: %s", err),
		)
	}
	ack := channeltypes.NewResultAcknowledgement(jsonResource)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {

	return sdkerrors.Wrap(resourcetypes.ErrUnexpectedAck, "unexpected acknowledgement")
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return sdkerrors.Wrap(resourcetypes.ErrUnexpectedAck, "unexpected acknowledgement")
}
