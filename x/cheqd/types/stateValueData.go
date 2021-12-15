package types

import "github.com/gogo/protobuf/proto"

// StateValueData is interface uniting possible types to be used for stateValue.data field
type StateValueData interface {
	proto.Message
}
