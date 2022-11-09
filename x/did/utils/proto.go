package utils

import "github.com/gogo/protobuf/proto"

// MsgTypeURL returns the TypeURL of a `proto.Message`.
func MsgTypeURL(msg proto.Message) string {
	return "/" + proto.MessageName(msg)
}
