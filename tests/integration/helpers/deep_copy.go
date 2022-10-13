package helpers

import (
	"reflect"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

// IDeepCopy is an interface for deep copy in the decorator pattern.
type IDeepCopy interface {
	DeepCopy(src interface{}) interface{}
}

// TDeepCopy is a decorator for deep copy.
type TDeepCopy struct{}

func (d *TDeepCopy) DeepCopy(src interface{}) interface{} {
	return deepCopy(src)
}

// DeepCopyCreateDid is a decorator for deep copy of type MsgCreateDidPayload.
type DeepCopyCreateDid struct {
	TDeepCopy
}

func (d *DeepCopyCreateDid) DeepCopy(src cheqdtypes.MsgCreateDidPayload) cheqdtypes.MsgCreateDidPayload {
	return deepCopy(src).(cheqdtypes.MsgCreateDidPayload)
}

// DeepCopyUpdateDid is a decorator for deep copy of type MsgUpdateDidPayload.
type DeepCopyUpdateDid struct {
	TDeepCopy
}

func (d *DeepCopyUpdateDid) DeepCopy(src cheqdtypes.MsgUpdateDidPayload) cheqdtypes.MsgUpdateDidPayload {
	return deepCopy(src).(cheqdtypes.MsgUpdateDidPayload)
}

// DeepCopyCreateResource is a decorator for deep copy of type MsgCreateResource.
type DeepCopyCreateResource struct {
	TDeepCopy
}

func (d *DeepCopyCreateResource) DeepCopy(src resourcetypes.MsgCreateResource) resourcetypes.MsgCreateResource {
	return deepCopy(src).(resourcetypes.MsgCreateResource)
}

// TODO: Add generics after bumping to Go 1.18 and remove this workaround.
func deepCopy(src interface{}) interface{} {
	var dst, reflection reflect.Value

	switch actualSrc := (src).(type) {
	case cheqdtypes.MsgCreateDidPayload:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1)
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []cheqdtypes.MsgCreateDidPayload{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(cheqdtypes.MsgCreateDidPayload)
	case cheqdtypes.MsgUpdateDidPayload:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1)
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []cheqdtypes.MsgUpdateDidPayload{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(cheqdtypes.MsgUpdateDidPayload)
	case resourcetypes.MsgCreateResource:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1)
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []resourcetypes.MsgCreateResource{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(resourcetypes.MsgCreateResource)
	default:
		panic("Unsupported type")
	}
}
