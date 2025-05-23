// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package didv2

import (
	_ "cosmossdk.io/api/amino"
	v1beta1 "cosmossdk.io/api/cosmos/base/v1beta1"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_FeeParams                protoreflect.MessageDescriptor
	fd_FeeParams_create_did     protoreflect.FieldDescriptor
	fd_FeeParams_update_did     protoreflect.FieldDescriptor
	fd_FeeParams_deactivate_did protoreflect.FieldDescriptor
	fd_FeeParams_burn_factor    protoreflect.FieldDescriptor
)

func init() {
	file_cheqd_did_v2_fee_proto_init()
	md_FeeParams = File_cheqd_did_v2_fee_proto.Messages().ByName("FeeParams")
	fd_FeeParams_create_did = md_FeeParams.Fields().ByName("create_did")
	fd_FeeParams_update_did = md_FeeParams.Fields().ByName("update_did")
	fd_FeeParams_deactivate_did = md_FeeParams.Fields().ByName("deactivate_did")
	fd_FeeParams_burn_factor = md_FeeParams.Fields().ByName("burn_factor")
}

var _ protoreflect.Message = (*fastReflection_FeeParams)(nil)

type fastReflection_FeeParams FeeParams

func (x *FeeParams) ProtoReflect() protoreflect.Message {
	return (*fastReflection_FeeParams)(x)
}

func (x *FeeParams) slowProtoReflect() protoreflect.Message {
	mi := &file_cheqd_did_v2_fee_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_FeeParams_messageType fastReflection_FeeParams_messageType
var _ protoreflect.MessageType = fastReflection_FeeParams_messageType{}

type fastReflection_FeeParams_messageType struct{}

func (x fastReflection_FeeParams_messageType) Zero() protoreflect.Message {
	return (*fastReflection_FeeParams)(nil)
}
func (x fastReflection_FeeParams_messageType) New() protoreflect.Message {
	return new(fastReflection_FeeParams)
}
func (x fastReflection_FeeParams_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_FeeParams
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_FeeParams) Descriptor() protoreflect.MessageDescriptor {
	return md_FeeParams
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_FeeParams) Type() protoreflect.MessageType {
	return _fastReflection_FeeParams_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_FeeParams) New() protoreflect.Message {
	return new(fastReflection_FeeParams)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_FeeParams) Interface() protoreflect.ProtoMessage {
	return (*FeeParams)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_FeeParams) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.CreateDid != nil {
		value := protoreflect.ValueOfMessage(x.CreateDid.ProtoReflect())
		if !f(fd_FeeParams_create_did, value) {
			return
		}
	}
	if x.UpdateDid != nil {
		value := protoreflect.ValueOfMessage(x.UpdateDid.ProtoReflect())
		if !f(fd_FeeParams_update_did, value) {
			return
		}
	}
	if x.DeactivateDid != nil {
		value := protoreflect.ValueOfMessage(x.DeactivateDid.ProtoReflect())
		if !f(fd_FeeParams_deactivate_did, value) {
			return
		}
	}
	if x.BurnFactor != "" {
		value := protoreflect.ValueOfString(x.BurnFactor)
		if !f(fd_FeeParams_burn_factor, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_FeeParams) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		return x.CreateDid != nil
	case "cheqd.did.v2.FeeParams.update_did":
		return x.UpdateDid != nil
	case "cheqd.did.v2.FeeParams.deactivate_did":
		return x.DeactivateDid != nil
	case "cheqd.did.v2.FeeParams.burn_factor":
		return x.BurnFactor != ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_FeeParams) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		x.CreateDid = nil
	case "cheqd.did.v2.FeeParams.update_did":
		x.UpdateDid = nil
	case "cheqd.did.v2.FeeParams.deactivate_did":
		x.DeactivateDid = nil
	case "cheqd.did.v2.FeeParams.burn_factor":
		x.BurnFactor = ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_FeeParams) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		value := x.CreateDid
		return protoreflect.ValueOfMessage(value.ProtoReflect())
	case "cheqd.did.v2.FeeParams.update_did":
		value := x.UpdateDid
		return protoreflect.ValueOfMessage(value.ProtoReflect())
	case "cheqd.did.v2.FeeParams.deactivate_did":
		value := x.DeactivateDid
		return protoreflect.ValueOfMessage(value.ProtoReflect())
	case "cheqd.did.v2.FeeParams.burn_factor":
		value := x.BurnFactor
		return protoreflect.ValueOfString(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_FeeParams) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		x.CreateDid = value.Message().Interface().(*v1beta1.Coin)
	case "cheqd.did.v2.FeeParams.update_did":
		x.UpdateDid = value.Message().Interface().(*v1beta1.Coin)
	case "cheqd.did.v2.FeeParams.deactivate_did":
		x.DeactivateDid = value.Message().Interface().(*v1beta1.Coin)
	case "cheqd.did.v2.FeeParams.burn_factor":
		x.BurnFactor = value.Interface().(string)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_FeeParams) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		if x.CreateDid == nil {
			x.CreateDid = new(v1beta1.Coin)
		}
		return protoreflect.ValueOfMessage(x.CreateDid.ProtoReflect())
	case "cheqd.did.v2.FeeParams.update_did":
		if x.UpdateDid == nil {
			x.UpdateDid = new(v1beta1.Coin)
		}
		return protoreflect.ValueOfMessage(x.UpdateDid.ProtoReflect())
	case "cheqd.did.v2.FeeParams.deactivate_did":
		if x.DeactivateDid == nil {
			x.DeactivateDid = new(v1beta1.Coin)
		}
		return protoreflect.ValueOfMessage(x.DeactivateDid.ProtoReflect())
	case "cheqd.did.v2.FeeParams.burn_factor":
		panic(fmt.Errorf("field burn_factor of message cheqd.did.v2.FeeParams is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_FeeParams) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "cheqd.did.v2.FeeParams.create_did":
		m := new(v1beta1.Coin)
		return protoreflect.ValueOfMessage(m.ProtoReflect())
	case "cheqd.did.v2.FeeParams.update_did":
		m := new(v1beta1.Coin)
		return protoreflect.ValueOfMessage(m.ProtoReflect())
	case "cheqd.did.v2.FeeParams.deactivate_did":
		m := new(v1beta1.Coin)
		return protoreflect.ValueOfMessage(m.ProtoReflect())
	case "cheqd.did.v2.FeeParams.burn_factor":
		return protoreflect.ValueOfString("")
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cheqd.did.v2.FeeParams"))
		}
		panic(fmt.Errorf("message cheqd.did.v2.FeeParams does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_FeeParams) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in cheqd.did.v2.FeeParams", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_FeeParams) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_FeeParams) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_FeeParams) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_FeeParams) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*FeeParams)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		if x.CreateDid != nil {
			l = options.Size(x.CreateDid)
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.UpdateDid != nil {
			l = options.Size(x.UpdateDid)
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.DeactivateDid != nil {
			l = options.Size(x.DeactivateDid)
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.BurnFactor)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*FeeParams)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if len(x.BurnFactor) > 0 {
			i -= len(x.BurnFactor)
			copy(dAtA[i:], x.BurnFactor)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.BurnFactor)))
			i--
			dAtA[i] = 0x22
		}
		if x.DeactivateDid != nil {
			encoded, err := options.Marshal(x.DeactivateDid)
			if err != nil {
				return protoiface.MarshalOutput{
					NoUnkeyedLiterals: input.NoUnkeyedLiterals,
					Buf:               input.Buf,
				}, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
			i--
			dAtA[i] = 0x1a
		}
		if x.UpdateDid != nil {
			encoded, err := options.Marshal(x.UpdateDid)
			if err != nil {
				return protoiface.MarshalOutput{
					NoUnkeyedLiterals: input.NoUnkeyedLiterals,
					Buf:               input.Buf,
				}, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
			i--
			dAtA[i] = 0x12
		}
		if x.CreateDid != nil {
			encoded, err := options.Marshal(x.CreateDid)
			if err != nil {
				return protoiface.MarshalOutput{
					NoUnkeyedLiterals: input.NoUnkeyedLiterals,
					Buf:               input.Buf,
				}, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*FeeParams)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				wire |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			fieldNum := int32(wire >> 3)
			wireType := int(wire & 0x7)
			if wireType == 4 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: FeeParams: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: FeeParams: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field CreateDid", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if x.CreateDid == nil {
					x.CreateDid = &v1beta1.Coin{}
				}
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.CreateDid); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field UpdateDid", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if x.UpdateDid == nil {
					x.UpdateDid = &v1beta1.Coin{}
				}
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.UpdateDid); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field DeactivateDid", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if x.DeactivateDid == nil {
					x.DeactivateDid = &v1beta1.Coin{}
				}
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.DeactivateDid); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field BurnFactor", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.BurnFactor = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: cheqd/did/v2/fee.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// FeeParams defines the parameters for the cheqd DID module fixed fee
type FeeParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Fixed fee for creating a DID
	//
	// Default: 50 CHEQ or 50000000000ncheq
	CreateDid *v1beta1.Coin `protobuf:"bytes,1,opt,name=create_did,json=createDid,proto3" json:"create_did,omitempty"`
	// Fixed fee for updating a DID
	//
	// Default: 25 CHEQ or 25000000000ncheq
	UpdateDid *v1beta1.Coin `protobuf:"bytes,2,opt,name=update_did,json=updateDid,proto3" json:"update_did,omitempty"`
	// Fixed fee for deactivating a DID
	//
	// Default: 10 CHEQ or 10000000000ncheq
	DeactivateDid *v1beta1.Coin `protobuf:"bytes,3,opt,name=deactivate_did,json=deactivateDid,proto3" json:"deactivate_did,omitempty"`
	// Percentage of the fixed fee that will be burned
	//
	// Default: 0.5 (50%)
	BurnFactor string `protobuf:"bytes,4,opt,name=burn_factor,json=burnFactor,proto3" json:"burn_factor,omitempty"`
}

func (x *FeeParams) Reset() {
	*x = FeeParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cheqd_did_v2_fee_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FeeParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FeeParams) ProtoMessage() {}

// Deprecated: Use FeeParams.ProtoReflect.Descriptor instead.
func (*FeeParams) Descriptor() ([]byte, []int) {
	return file_cheqd_did_v2_fee_proto_rawDescGZIP(), []int{0}
}

func (x *FeeParams) GetCreateDid() *v1beta1.Coin {
	if x != nil {
		return x.CreateDid
	}
	return nil
}

func (x *FeeParams) GetUpdateDid() *v1beta1.Coin {
	if x != nil {
		return x.UpdateDid
	}
	return nil
}

func (x *FeeParams) GetDeactivateDid() *v1beta1.Coin {
	if x != nil {
		return x.DeactivateDid
	}
	return nil
}

func (x *FeeParams) GetBurnFactor() string {
	if x != nil {
		return x.BurnFactor
	}
	return ""
}

var File_cheqd_did_v2_fee_proto protoreflect.FileDescriptor

var file_cheqd_did_v2_fee_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x68, 0x65, 0x71, 0x64, 0x2f, 0x64, 0x69, 0x64, 0x2f, 0x76, 0x32, 0x2f, 0x66,
	0x65, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x68, 0x65, 0x71, 0x64, 0x2e,
	0x64, 0x69, 0x64, 0x2e, 0x76, 0x32, 0x1a, 0x11, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x6d,
	0x69, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x63, 0x6f, 0x73, 0x6d, 0x6f,
	0x73, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63,
	0x6f, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x63, 0x6f, 0x73, 0x6d, 0x6f,
	0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x02, 0x0a, 0x09, 0x46,
	0x65, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x3e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x5f, 0x64, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x69, 0x64, 0x12, 0x3e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x64, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x69, 0x64, 0x12, 0x46, 0x0a, 0x0e, 0x64, 0x65, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x64, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x42, 0x04, 0xc8, 0xde, 0x1f,
	0x00, 0x52, 0x0d, 0x64, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x44, 0x69, 0x64,
	0x12, 0x57, 0x0a, 0x0b, 0x62, 0x75, 0x72, 0x6e, 0x5f, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0xc8, 0xde, 0x1f, 0x00, 0xda, 0xde, 0x1f, 0x1b, 0x63,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x6d, 0x61, 0x74, 0x68,
	0x2e, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x44, 0x65, 0x63, 0xd2, 0xb4, 0x2d, 0x0a, 0x63, 0x6f,
	0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x63, 0xa8, 0xe7, 0xb0, 0x2a, 0x01, 0x52, 0x0a, 0x62,
	0x75, 0x72, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x42, 0xa9, 0x01, 0xa8, 0xe2, 0x1e, 0x01,
	0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x68, 0x65, 0x71, 0x64, 0x2e, 0x64, 0x69, 0x64, 0x2e,
	0x76, 0x32, 0x42, 0x08, 0x46, 0x65, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x35,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x65, 0x71, 0x64,
	0x2f, 0x63, 0x68, 0x65, 0x71, 0x64, 0x2d, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x32, 0x2f, 0x63, 0x68, 0x65, 0x71, 0x64, 0x2f, 0x64, 0x69, 0x64, 0x2f, 0x76, 0x32, 0x3b,
	0x64, 0x69, 0x64, 0x76, 0x32, 0xa2, 0x02, 0x03, 0x43, 0x44, 0x58, 0xaa, 0x02, 0x0c, 0x43, 0x68,
	0x65, 0x71, 0x64, 0x2e, 0x44, 0x69, 0x64, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x0c, 0x43, 0x68, 0x65,
	0x71, 0x64, 0x5c, 0x44, 0x69, 0x64, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x18, 0x43, 0x68, 0x65, 0x71,
	0x64, 0x5c, 0x44, 0x69, 0x64, 0x5c, 0x56, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x43, 0x68, 0x65, 0x71, 0x64, 0x3a, 0x3a, 0x44, 0x69,
	0x64, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cheqd_did_v2_fee_proto_rawDescOnce sync.Once
	file_cheqd_did_v2_fee_proto_rawDescData = file_cheqd_did_v2_fee_proto_rawDesc
)

func file_cheqd_did_v2_fee_proto_rawDescGZIP() []byte {
	file_cheqd_did_v2_fee_proto_rawDescOnce.Do(func() {
		file_cheqd_did_v2_fee_proto_rawDescData = protoimpl.X.CompressGZIP(file_cheqd_did_v2_fee_proto_rawDescData)
	})
	return file_cheqd_did_v2_fee_proto_rawDescData
}

var file_cheqd_did_v2_fee_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_cheqd_did_v2_fee_proto_goTypes = []interface{}{
	(*FeeParams)(nil),    // 0: cheqd.did.v2.FeeParams
	(*v1beta1.Coin)(nil), // 1: cosmos.base.v1beta1.Coin
}
var file_cheqd_did_v2_fee_proto_depIdxs = []int32{
	1, // 0: cheqd.did.v2.FeeParams.create_did:type_name -> cosmos.base.v1beta1.Coin
	1, // 1: cheqd.did.v2.FeeParams.update_did:type_name -> cosmos.base.v1beta1.Coin
	1, // 2: cheqd.did.v2.FeeParams.deactivate_did:type_name -> cosmos.base.v1beta1.Coin
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_cheqd_did_v2_fee_proto_init() }
func file_cheqd_did_v2_fee_proto_init() {
	if File_cheqd_did_v2_fee_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cheqd_did_v2_fee_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FeeParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cheqd_did_v2_fee_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cheqd_did_v2_fee_proto_goTypes,
		DependencyIndexes: file_cheqd_did_v2_fee_proto_depIdxs,
		MessageInfos:      file_cheqd_did_v2_fee_proto_msgTypes,
	}.Build()
	File_cheqd_did_v2_fee_proto = out.File
	file_cheqd_did_v2_fee_proto_rawDesc = nil
	file_cheqd_did_v2_fee_proto_goTypes = nil
	file_cheqd_did_v2_fee_proto_depIdxs = nil
}
