package	migrations

import (
	"reflect"

	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/gogo/protobuf/proto"
)

type StateValueData interface {
	proto.Message
}

func UnpackData(stateValue didtypesv1.StateValue) (StateValueData, error) {
	value, isOk := stateValue.Data.GetCachedValue().(StateValueData)
	if !isOk {
		return nil, didtypes.ErrUnpackStateValue.Wrapf("invalid type url: %s", stateValue.Data.TypeUrl)
	}

	return value, nil
}

func UnpackDataAsDid(stateValue didtypesv1.StateValue) (*didtypesv1.Did, error) {
	data, err := UnpackData(stateValue)
	if err != nil {
		return nil, err
	}

	value, isValue := data.(*didtypesv1.Did)
	if !isValue {
		return nil, didtypes.ErrUnpackStateValue.Wrap(reflect.TypeOf(data).String())
	}

	return value, nil
}

func StateValueToDIDDocWithMetadata(stateValue didtypesv1.StateValue) (didtypes.DidDocWithMetadata, error) {
	didDoc, err := UnpackDataAsDid(stateValue)
	if err != nil {
		return didtypes.DidDocWithMetadata{}, err
	}
	didDoc = didtypes.

	didtypes.NewDidDocWithMetadata(&didDoc, &metadata)
}