package helpers

import (
	"github.com/cheqd/cheqd-node/app/params"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

var (
	Codec    codec.Codec
	Registry types.InterfaceRegistry
)

func init() {
	encodingConfig := params.MakeEncodingConfig()
	Codec = encodingConfig.Codec
	Registry = encodingConfig.InterfaceRegistry
}
