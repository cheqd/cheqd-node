package legacy

import (
	"reflect"
	"time"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

// StateValueData is interface uniting possible types to be used for stateValue.data field
type StateValueData interface {
	proto.Message
}

var _ types.UnpackInterfacesMessage = &StateValue{}

func (m *StateValue) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var data StateValueData
	return unpacker.UnpackAny(m.Data, &data)
}

func NewStateValue(data StateValueData, metadata *Metadata) (StateValue, error) {
	any, err := types.NewAnyWithValue(data)
	if err != nil {
		return StateValue{}, err
	}

	return StateValue{Data: any, Metadata: metadata}, nil
}

func NewMetadataFromContext(ctx sdk.Context) Metadata {
	created := ctx.BlockTime().Format(time.RFC3339)
	txHash := utils.GetTxHash(ctx.TxBytes())

	return Metadata{Created: created, Deactivated: false, VersionId: txHash}
}

func (m *Metadata) Update(ctx sdk.Context) {
	m.Updated = ctx.BlockTime().Format(time.RFC3339)
}

func (m StateValue) UnpackData() (StateValueData, error) {
	value, isOk := m.Data.GetCachedValue().(StateValueData)
	if !isOk {
		return nil, ErrUnpackStateValue.Wrapf("invalid type url: %s", m.Data.TypeUrl)
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
		return nil, ErrUnpackStateValue.Wrap(reflect.TypeOf(data).String())
	}

	return value, nil
}
