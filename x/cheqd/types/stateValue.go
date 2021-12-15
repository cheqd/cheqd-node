package types

import (
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"reflect"
)

var _ types.UnpackInterfacesMessage = &StateValue{}

func (m *StateValue) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var data StateValueData
	return unpacker.UnpackAny(m.Data, &data)
}

func NewStateValue(data StateValueData, metadata *Metadata) (*StateValue, error) {
	any, err := types.NewAnyWithValue(data)
	if err != nil {
		return nil, ErrInvalidDidStateValue.Wrap(err.Error())
	}

	return &StateValue{Data: any, Metadata: metadata}, nil
}

func NewMetadata(ctx sdk.Context) Metadata {
	created := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	return Metadata{Created: created, Updated: created, Deactivated: false, VersionId: txHash}
}

func (m StateValue) UnpackData() (StateValueData, error) {
	value, isOk := m.Data.GetCachedValue().(StateValueData)
	if !isOk {
		return nil, ErrInvalidDidStateValue.Wrap(m.Data.TypeUrl)
	}

	return value, nil
}

func (m StateValue) UnpackDataAsDid() (*Did, error) {
	data, err := m.UnpackData()
	if err != nil {
		return nil, err
	}

	value, isValue := data.(*Did)
	if !isValue {
		return nil, ErrInvalidDidStateValue.Wrap(reflect.TypeOf(data).String())
	}

	return value, nil
}
