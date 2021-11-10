package v1

import (
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	StateValueDid = "/cheqdid.cheqdnode.cheqd.v1.Did"
)

func NewStateValue(msg proto.Message, metadata *Metadata) (*StateValue, error) {
	data, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, ErrInvalidDidStateValue.Wrap(err.Error())
	}

	return &StateValue{Data: data, Metadata: metadata}, nil
}

func NewMetadata(ctx sdk.Context) Metadata {
	created := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	return Metadata{Created: created, Updated: created, Deactivated: false, VersionId: txHash}
}

func (m StateValue) GetDid() (*Did, error) {
	value, isValue := m.Data.GetCachedValue().(Did)
	if isValue {
		return &value, nil
	}

	if m.Data.TypeUrl != StateValueDid {
		return nil, ErrInvalidDidStateValue.Wrap(m.Data.TypeUrl)
	}

	state := Did{}
	err := state.Unmarshal(m.Data.Value)
	if err != nil {
		return nil, ErrInvalidDidStateValue.Wrap(err.Error())
	}

	return &state, nil
}
