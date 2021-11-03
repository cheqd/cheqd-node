package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateDid{}, "cheqd/CreateDid", nil)
	cdc.RegisterConcrete(&MsgUpdateDid{}, "cheqd/UpdateDid", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDid{},
		&MsgUpdateDid{},
	)

	registry.RegisterInterface(MessageCreateDid, (*IdentityMsg)(nil), &MsgCreateDidPayload{})
	registry.RegisterInterface(MessageUpdateDid, (*IdentityMsg)(nil), &MsgUpdateDidPayload{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
