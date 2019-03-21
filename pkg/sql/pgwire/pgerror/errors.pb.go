// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/pgwire/pgerror/errors.proto

package pgerror

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Error contains all Postgres wire protocol error fields.
// See https://www.postgresql.org/docs/current/static/protocol-error-fields.html
// for a list of all Postgres error fields, most of which are optional and can
// be used to provide auxiliary error information.
type Error struct {
	// standard pg error fields. This can be passed
	// over the pg wire protocol.
	Code    string        `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string        `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Detail  string        `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	Hint    string        `protobuf:"bytes,4,opt,name=hint,proto3" json:"hint,omitempty"`
	Source  *Error_Source `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
	// a telemetry key, used as telemetry counter name.
	// Typically of the form [<prefix>.]#issuenum[.details]
	TelemetryKey         string   `protobuf:"bytes,6,opt,name=telemetry_key,json=telemetryKey,proto3" json:"telemetry_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_e4ea9b89df718931, []int{0}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
}
func (m *Error) XXX_Size() int {
	return m.Size()
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

type Error_Source struct {
	File                 string   `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Line                 int32    `protobuf:"varint,2,opt,name=line,proto3" json:"line,omitempty"`
	Function             string   `protobuf:"bytes,3,opt,name=function,proto3" json:"function,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error_Source) Reset()         { *m = Error_Source{} }
func (m *Error_Source) String() string { return proto.CompactTextString(m) }
func (*Error_Source) ProtoMessage()    {}
func (*Error_Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_e4ea9b89df718931, []int{0, 0}
}
func (m *Error_Source) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Error_Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *Error_Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error_Source.Merge(dst, src)
}
func (m *Error_Source) XXX_Size() int {
	return m.Size()
}
func (m *Error_Source) XXX_DiscardUnknown() {
	xxx_messageInfo_Error_Source.DiscardUnknown(m)
}

var xxx_messageInfo_Error_Source proto.InternalMessageInfo

// SafeDetailPayload encapsulates safe (PII-free) additional details.
type SafeDetailPayload struct {
	// safe_message encodes the message as per log.ReportablesToSafeError.
	SafeMessage string `protobuf:"bytes,1,opt,name=safe_message,json=safeMessage,proto3" json:"safe_message,omitempty"`
	// encoded_stack_trace encodes the stack trace as per log.EncodeStackTrace.
	EncodedStackTrace string `protobuf:"bytes,2,opt,name=encoded_stack_trace,json=encodedStackTrace,proto3" json:"encoded_stack_trace,omitempty"`
	// source is the head of the stack trace.
	Source               *Error_Source `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SafeDetailPayload) Reset()         { *m = SafeDetailPayload{} }
func (m *SafeDetailPayload) String() string { return proto.CompactTextString(m) }
func (*SafeDetailPayload) ProtoMessage()    {}
func (*SafeDetailPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_e4ea9b89df718931, []int{1}
}
func (m *SafeDetailPayload) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SafeDetailPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *SafeDetailPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SafeDetailPayload.Merge(dst, src)
}
func (m *SafeDetailPayload) XXX_Size() int {
	return m.Size()
}
func (m *SafeDetailPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_SafeDetailPayload.DiscardUnknown(m)
}

var xxx_messageInfo_SafeDetailPayload proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Error)(nil), "cockroach.pgerror.Error")
	proto.RegisterType((*Error_Source)(nil), "cockroach.pgerror.Error.Source")
	proto.RegisterType((*SafeDetailPayload)(nil), "cockroach.pgerror.SafeDetailPayload")
}
func (m *Error) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Error) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Code) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Code)))
		i += copy(dAtA[i:], m.Code)
	}
	if len(m.Message) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	if len(m.Detail) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Detail)))
		i += copy(dAtA[i:], m.Detail)
	}
	if len(m.Hint) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Hint)))
		i += copy(dAtA[i:], m.Hint)
	}
	if m.Source != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.Source.Size()))
		n1, err := m.Source.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.TelemetryKey) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.TelemetryKey)))
		i += copy(dAtA[i:], m.TelemetryKey)
	}
	return i, nil
}

func (m *Error_Source) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Error_Source) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.File) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.File)))
		i += copy(dAtA[i:], m.File)
	}
	if m.Line != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.Line))
	}
	if len(m.Function) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Function)))
		i += copy(dAtA[i:], m.Function)
	}
	return i, nil
}

func (m *SafeDetailPayload) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SafeDetailPayload) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.SafeMessage) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.SafeMessage)))
		i += copy(dAtA[i:], m.SafeMessage)
	}
	if len(m.EncodedStackTrace) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.EncodedStackTrace)))
		i += copy(dAtA[i:], m.EncodedStackTrace)
	}
	if m.Source != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.Source.Size()))
		n2, err := m.Source.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func encodeVarintErrors(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Error) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.Detail)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.Hint)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	if m.Source != nil {
		l = m.Source.Size()
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.TelemetryKey)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}

func (m *Error_Source) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.File)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	if m.Line != 0 {
		n += 1 + sovErrors(uint64(m.Line))
	}
	l = len(m.Function)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}

func (m *SafeDetailPayload) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SafeMessage)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.EncodedStackTrace)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	if m.Source != nil {
		l = m.Source.Size()
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}

func sovErrors(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozErrors(x uint64) (n int) {
	return sovErrors(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Error) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Error: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Error: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Detail", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Detail = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Source == nil {
				m.Source = &Error_Source{}
			}
			if err := m.Source.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TelemetryKey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TelemetryKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthErrors
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
func (m *Error_Source) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Source: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Source: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field File", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.File = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Line", wireType)
			}
			m.Line = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Line |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Function", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Function = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthErrors
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
func (m *SafeDetailPayload) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SafeDetailPayload: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SafeDetailPayload: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SafeMessage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SafeMessage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncodedStackTrace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncodedStackTrace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Source == nil {
				m.Source = &Error_Source{}
			}
			if err := m.Source.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthErrors
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
func skipErrors(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowErrors
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
					return 0, ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowErrors
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthErrors
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowErrors
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipErrors(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthErrors = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowErrors   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("sql/pgwire/pgerror/errors.proto", fileDescriptor_errors_e4ea9b89df718931)
}

var fileDescriptor_errors_e4ea9b89df718931 = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x51, 0xbd, 0x4e, 0xf3, 0x30,
	0x14, 0xad, 0xfb, 0x93, 0xf6, 0x73, 0xfb, 0x49, 0x5f, 0xfd, 0x49, 0xc8, 0xea, 0x90, 0xfe, 0xb0,
	0x74, 0x4a, 0x25, 0x18, 0xd8, 0x11, 0x2c, 0xfc, 0x48, 0xa8, 0x65, 0x62, 0x89, 0x8c, 0x73, 0xd3,
	0x46, 0x4d, 0xe3, 0x62, 0xbb, 0x42, 0x79, 0x0b, 0x24, 0x1e, 0x80, 0xd7, 0xe9, 0xc8, 0xc8, 0x08,
	0xe1, 0x45, 0x90, 0x1d, 0xb7, 0x42, 0x62, 0x62, 0x89, 0xce, 0x3d, 0xc7, 0xf7, 0xde, 0x73, 0x4f,
	0x70, 0x5f, 0x3d, 0xa4, 0x93, 0xf5, 0xfc, 0x31, 0x91, 0x30, 0x59, 0xcf, 0x41, 0x4a, 0x21, 0x27,
	0xf6, 0xab, 0x82, 0xb5, 0x14, 0x5a, 0x90, 0x2e, 0x17, 0x7c, 0x29, 0x05, 0xe3, 0x8b, 0xc0, 0xe9,
	0xa3, 0xe7, 0x2a, 0x6e, 0x9c, 0x1b, 0x44, 0x08, 0xae, 0x73, 0x11, 0x01, 0x45, 0x03, 0x34, 0xfe,
	0x33, 0xb5, 0x98, 0x50, 0xdc, 0x5c, 0x81, 0x52, 0x6c, 0x0e, 0xb4, 0x6a, 0xe9, 0x5d, 0x49, 0x0e,
	0xb0, 0x17, 0x81, 0x66, 0x49, 0x4a, 0x6b, 0x56, 0x70, 0x95, 0x99, 0xb2, 0x48, 0x32, 0x4d, 0xeb,
	0xe5, 0x14, 0x83, 0xc9, 0x09, 0xf6, 0x94, 0xd8, 0x48, 0x0e, 0xb4, 0x31, 0x40, 0xe3, 0xf6, 0x51,
	0x3f, 0xf8, 0xe1, 0x23, 0xb0, 0x1e, 0x82, 0x99, 0x7d, 0x36, 0x75, 0xcf, 0xc9, 0x21, 0xfe, 0xab,
	0x21, 0x85, 0x15, 0x68, 0x99, 0x87, 0x4b, 0xc8, 0xa9, 0x67, 0xa7, 0x76, 0xf6, 0xe4, 0x25, 0xe4,
	0xbd, 0x2b, 0xec, 0x95, 0x6d, 0x66, 0x77, 0x9c, 0xa4, 0xfb, 0x0b, 0x0c, 0x36, 0x5c, 0x9a, 0x64,
	0xa5, 0xfd, 0xc6, 0xd4, 0x62, 0xd2, 0xc3, 0xad, 0x78, 0x93, 0x71, 0x9d, 0x88, 0xcc, 0xb9, 0xdf,
	0xd7, 0x17, 0xf5, 0x56, 0xf3, 0x5f, 0x6b, 0xf4, 0x82, 0x70, 0x77, 0xc6, 0x62, 0x38, 0xb3, 0x47,
	0xdd, 0xb0, 0x3c, 0x15, 0x2c, 0x22, 0x43, 0xdc, 0x51, 0x2c, 0x86, 0x70, 0x17, 0x49, 0xb9, 0xa7,
	0x6d, 0xb8, 0x6b, 0x17, 0x4b, 0x80, 0xff, 0x43, 0x66, 0xa2, 0x8b, 0x42, 0xa5, 0x19, 0x5f, 0x86,
	0x5a, 0x32, 0xbe, 0x0b, 0xaf, 0xeb, 0xa4, 0x99, 0x51, 0x6e, 0x8d, 0xf0, 0x2d, 0x9a, 0xda, 0xaf,
	0xa2, 0x39, 0x1d, 0x6e, 0x3f, 0xfc, 0xca, 0xb6, 0xf0, 0xd1, 0x6b, 0xe1, 0xa3, 0xb7, 0xc2, 0x47,
	0xef, 0x85, 0x8f, 0x9e, 0x3e, 0xfd, 0xca, 0x5d, 0xd3, 0xf5, 0xdd, 0x7b, 0xf6, 0xa7, 0x1f, 0x7f,
	0x05, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x70, 0x0b, 0x13, 0x17, 0x02, 0x00, 0x00,
}
