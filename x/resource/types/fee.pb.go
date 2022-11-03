// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cheqd/resource/v2/fee.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// FeeParams defines the parameters for the `resource` module fixed fee.
type FeeParams struct {
	// Media types define the fixed fee each for the `resource` module.
	MediaTypes map[string]types.Coin                  `protobuf:"bytes,1,rep,name=media_types,json=mediaTypes,proto3" json:"media_types" yaml:"media_types" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	BurnFactor github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=burn_factor,json=burnFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"burn_factor"`
}

func (m *FeeParams) Reset()         { *m = FeeParams{} }
func (m *FeeParams) String() string { return proto.CompactTextString(m) }
func (*FeeParams) ProtoMessage()    {}
func (*FeeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_133abe56c2e24f1e, []int{0}
}
func (m *FeeParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeeParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeeParams.Merge(m, src)
}
func (m *FeeParams) XXX_Size() int {
	return m.Size()
}
func (m *FeeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_FeeParams.DiscardUnknown(m)
}

var xxx_messageInfo_FeeParams proto.InternalMessageInfo

func (m *FeeParams) GetMediaTypes() map[string]types.Coin {
	if m != nil {
		return m.MediaTypes
	}
	return nil
}

func init() {
	proto.RegisterType((*FeeParams)(nil), "cheqd.resource.v2.FeeParams")
	proto.RegisterMapType((map[string]types.Coin)(nil), "cheqd.resource.v2.FeeParams.MediaTypesEntry")
}

func init() { proto.RegisterFile("cheqd/resource/v2/fee.proto", fileDescriptor_133abe56c2e24f1e) }

var fileDescriptor_133abe56c2e24f1e = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0xce, 0x48, 0x2d,
	0x4c, 0xd1, 0x2f, 0x4a, 0x2d, 0xce, 0x2f, 0x2d, 0x4a, 0x4e, 0xd5, 0x2f, 0x33, 0xd2, 0x4f, 0x4b,
	0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x04, 0x4b, 0xea, 0xc1, 0x24, 0xf5, 0xca,
	0x8c, 0xa4, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0xb2, 0xfa, 0x20, 0x16, 0x44, 0xa1, 0x94, 0x5c,
	0x72, 0x7e, 0x71, 0x6e, 0x7e, 0xb1, 0x7e, 0x52, 0x62, 0x71, 0xaa, 0x7e, 0x99, 0x61, 0x52, 0x6a,
	0x49, 0xa2, 0xa1, 0x7e, 0x72, 0x7e, 0x66, 0x1e, 0x54, 0x5e, 0x12, 0x22, 0x1f, 0x0f, 0xd1, 0x08,
	0xe1, 0x40, 0xa4, 0x94, 0xd6, 0x33, 0x71, 0x71, 0xba, 0xa5, 0xa6, 0x06, 0x24, 0x16, 0x25, 0xe6,
	0x16, 0x0b, 0x65, 0x72, 0x71, 0xe7, 0xa6, 0xa6, 0x64, 0x26, 0xc6, 0x97, 0x54, 0x16, 0xa4, 0x16,
	0x4b, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0xe9, 0xe8, 0x61, 0xb8, 0x43, 0x0f, 0xae, 0x45, 0xcf,
	0x17, 0xa4, 0x3e, 0x04, 0xa4, 0xdc, 0x35, 0xaf, 0xa4, 0xa8, 0xd2, 0x49, 0xea, 0xc4, 0x3d, 0x79,
	0x86, 0x4f, 0xf7, 0xe4, 0x85, 0x2a, 0x13, 0x73, 0x73, 0xac, 0x94, 0x90, 0x8c, 0x53, 0x0a, 0xe2,
	0xca, 0x85, 0x2b, 0x16, 0x8a, 0xe5, 0xe2, 0x4e, 0x2a, 0x2d, 0xca, 0x8b, 0x4f, 0x4b, 0x4c, 0x2e,
	0xc9, 0x2f, 0x92, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x74, 0xb2, 0x01, 0x69, 0xbe, 0x75, 0x4f, 0x5e,
	0x2d, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x17, 0xea, 0x5c, 0x28, 0xa5, 0x5b,
	0x9c, 0x92, 0xad, 0x0f, 0x36, 0x4d, 0xcf, 0x25, 0x35, 0xf9, 0xd2, 0x16, 0x5d, 0x2e, 0xa8, 0x6f,
	0x5c, 0x52, 0x93, 0x83, 0xb8, 0x40, 0x06, 0xba, 0x81, 0xcd, 0x93, 0x8a, 0xe0, 0xe2, 0x47, 0x73,
	0x99, 0x90, 0x00, 0x17, 0x73, 0x76, 0x6a, 0xa5, 0x04, 0x23, 0xc8, 0xa6, 0x20, 0x10, 0x53, 0x48,
	0x9f, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34, 0x15, 0x6c, 0x3b, 0xb7, 0x91, 0xa4, 0x1e, 0xd4, 0x30,
	0x50, 0x38, 0xea, 0x41, 0xc3, 0x51, 0xcf, 0x39, 0x3f, 0x33, 0x2f, 0x08, 0xa2, 0xce, 0x8a, 0xc9,
	0x82, 0xd1, 0xc9, 0x6b, 0xc5, 0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63,
	0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96,
	0x63, 0x88, 0xd2, 0x41, 0x76, 0x39, 0x38, 0x6e, 0xc1, 0xa4, 0x6e, 0x5e, 0x7e, 0x4a, 0xaa, 0x7e,
	0x05, 0x22, 0xa2, 0xc1, 0x7e, 0x48, 0x62, 0x03, 0x47, 0x82, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x7b, 0xda, 0x05, 0x38, 0x07, 0x02, 0x00, 0x00,
}

func (this *FeeParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*FeeParams)
	if !ok {
		that2, ok := that.(FeeParams)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.MediaTypes) != len(that1.MediaTypes) {
		return false
	}
	for i := range this.MediaTypes {
		a := this.MediaTypes[i]
		b := that1.MediaTypes[i]
		if !(&a).Equal(&b) {
			return false
		}
	}
	if !this.BurnFactor.Equal(that1.BurnFactor) {
		return false
	}
	return true
}
func (m *FeeParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeeParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeeParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BurnFactor.Size()
		i -= size
		if _, err := m.BurnFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.MediaTypes) > 0 {
		for k := range m.MediaTypes {
			v := m.MediaTypes[k]
			baseI := i
			{
				size, err := (&v).MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFee(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintFee(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintFee(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintFee(dAtA []byte, offset int, v uint64) int {
	offset -= sovFee(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FeeParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MediaTypes) > 0 {
		for k, v := range m.MediaTypes {
			_ = k
			_ = v
			l = v.Size()
			mapEntrySize := 1 + len(k) + sovFee(uint64(len(k))) + 1 + l + sovFee(uint64(l))
			n += mapEntrySize + 1 + sovFee(uint64(mapEntrySize))
		}
	}
	l = m.BurnFactor.Size()
	n += 1 + l + sovFee(uint64(l))
	return n
}

func sovFee(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFee(x uint64) (n int) {
	return sovFee(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FeeParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFee
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
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
			return fmt.Errorf("proto: FeeParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeeParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MediaTypes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFee
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MediaTypes == nil {
				m.MediaTypes = make(map[string]types.Coin)
			}
			var mapkey string
			mapvalue := &types.Coin{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFee
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFee
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthFee
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthFee
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFee
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthFee
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthFee
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &types.Coin{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipFee(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthFee
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.MediaTypes[mapkey] = *mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BurnFactor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BurnFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFee
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipFee(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFee
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFee
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFee
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthFee
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFee
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFee
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFee        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFee          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFee = fmt.Errorf("proto: unexpected end of group")
)