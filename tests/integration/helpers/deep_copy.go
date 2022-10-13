package helpers

import (
	"reflect"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	/* resourcetypes "github.com/cheqd/cheqd-node/x/resource/types" */)

// TODO: Add generics after bumping to Go 1.18
/* type PayloadMsg interface{
	cheqdtypes.MsgCreateDidPayload | cheqdtypes.MsgUpdateDidPayload | resourcetypes.MsgCreateResourcePayload
} */

func DeepCopy(src cheqdtypes.MsgUpdateDidPayload) cheqdtypes.MsgUpdateDidPayload {
	dst := reflect.ValueOf(src).Elem()
	reflect.Copy(dst, reflect.ValueOf(src))
	return dst.Interface().(cheqdtypes.MsgUpdateDidPayload)
}
