package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Sdk messages
	cdc.RegisterConcrete(&MsgCreateResource{}, "resource/CreateResource", nil)

	// State value data
	cdc.RegisterInterface((*Resource)(nil), nil)
	cdc.RegisterConcrete(&Resource{}, "resource/Resource", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Sdk messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateResource{},
	)

	// State value data
	registry.RegisterInterface("Resource", (*Resource)(nil))
	registry.RegisterImplementations((*Resource)(nil), &Resource{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
