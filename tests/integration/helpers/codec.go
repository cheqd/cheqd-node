package helpers

import (
	"github.com/cheqd/cheqd-node/app/params"
	"github.com/cosmos/cosmos-sdk/codec"
)

var Codec codec.Codec

func init() {
	Codec = params.MakeEncodingConfig().Codec
}
