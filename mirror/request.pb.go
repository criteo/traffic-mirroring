// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: request.proto

package mirror

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Method int32

const (
	Method_GET     Method = 0
	Method_HEAD    Method = 1
	Method_POST    Method = 2
	Method_PUT     Method = 3
	Method_PATCH   Method = 4
	Method_DELETE  Method = 5
	Method_CONNECT Method = 6
	Method_OPTIONS Method = 7
	Method_TRACE   Method = 8
)

// Enum value maps for Method.
var (
	Method_name = map[int32]string{
		0: "GET",
		1: "HEAD",
		2: "POST",
		3: "PUT",
		4: "PATCH",
		5: "DELETE",
		6: "CONNECT",
		7: "OPTIONS",
		8: "TRACE",
	}
	Method_value = map[string]int32{
		"GET":     0,
		"HEAD":    1,
		"POST":    2,
		"PUT":     3,
		"PATCH":   4,
		"DELETE":  5,
		"CONNECT": 6,
		"OPTIONS": 7,
		"TRACE":   8,
	}
)

func (x Method) Enum() *Method {
	p := new(Method)
	*p = x
	return p
}

func (x Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Method) Descriptor() protoreflect.EnumDescriptor {
	return file_request_proto_enumTypes[0].Descriptor()
}

func (Method) Type() protoreflect.EnumType {
	return &file_request_proto_enumTypes[0]
}

func (x Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Method.Descriptor instead.
func (Method) EnumDescriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{0}
}

type HTTPVersion int32

const (
	HTTPVersion_HTTP1_0 HTTPVersion = 0
	HTTPVersion_HTTP1_1 HTTPVersion = 1
	HTTPVersion_HTTP2   HTTPVersion = 2
)

// Enum value maps for HTTPVersion.
var (
	HTTPVersion_name = map[int32]string{
		0: "HTTP1_0",
		1: "HTTP1_1",
		2: "HTTP2",
	}
	HTTPVersion_value = map[string]int32{
		"HTTP1_0": 0,
		"HTTP1_1": 1,
		"HTTP2":   2,
	}
)

func (x HTTPVersion) Enum() *HTTPVersion {
	p := new(HTTPVersion)
	*p = x
	return p
}

func (x HTTPVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HTTPVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_request_proto_enumTypes[1].Descriptor()
}

func (HTTPVersion) Type() protoreflect.EnumType {
	return &file_request_proto_enumTypes[1]
}

func (x HTTPVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HTTPVersion.Descriptor instead.
func (HTTPVersion) EnumDescriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{1}
}

type HeaderValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *HeaderValue) Reset() {
	*x = HeaderValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeaderValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderValue) ProtoMessage() {}

func (x *HeaderValue) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderValue.ProtoReflect.Descriptor instead.
func (*HeaderValue) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{0}
}

func (x *HeaderValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type MetaValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//	*MetaValue_String_
	//	*MetaValue_Int
	//	*MetaValue_Bool
	Value isMetaValue_Value `protobuf_oneof:"value"`
}

func (x *MetaValue) Reset() {
	*x = MetaValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaValue) ProtoMessage() {}

func (x *MetaValue) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaValue.ProtoReflect.Descriptor instead.
func (*MetaValue) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{1}
}

func (m *MetaValue) GetValue() isMetaValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *MetaValue) GetString_() string {
	if x, ok := x.GetValue().(*MetaValue_String_); ok {
		return x.String_
	}
	return ""
}

func (x *MetaValue) GetInt() int64 {
	if x, ok := x.GetValue().(*MetaValue_Int); ok {
		return x.Int
	}
	return 0
}

func (x *MetaValue) GetBool() bool {
	if x, ok := x.GetValue().(*MetaValue_Bool); ok {
		return x.Bool
	}
	return false
}

type isMetaValue_Value interface {
	isMetaValue_Value()
}

type MetaValue_String_ struct {
	String_ string `protobuf:"bytes,1,opt,name=string,proto3,oneof"`
}

type MetaValue_Int struct {
	Int int64 `protobuf:"varint,2,opt,name=int,proto3,oneof"`
}

type MetaValue_Bool struct {
	Bool bool `protobuf:"varint,3,opt,name=bool,proto3,oneof"`
}

func (*MetaValue_String_) isMetaValue_Value() {}

func (*MetaValue_Int) isMetaValue_Value() {}

func (*MetaValue_Bool) isMetaValue_Value() {}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time        *timestamppb.Timestamp  `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Method      Method                  `protobuf:"varint,2,opt,name=method,proto3,enum=mirror.Method" json:"method,omitempty"`
	Path        string                  `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	HttpVersion HTTPVersion             `protobuf:"varint,4,opt,name=http_version,json=httpVersion,proto3,enum=mirror.HTTPVersion" json:"http_version,omitempty"`
	Headers     map[string]*HeaderValue `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body        []byte                  `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	Meta        map[string]*MetaValue   `protobuf:"bytes,7,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{2}
}

func (x *Request) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Request) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_GET
}

func (x *Request) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Request) GetHttpVersion() HTTPVersion {
	if x != nil {
		return x.HttpVersion
	}
	return HTTPVersion_HTTP1_0
}

func (x *Request) GetHeaders() map[string]*HeaderValue {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Request) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Request) GetMeta() map[string]*MetaValue {
	if x != nil {
		return x.Meta
	}
	return nil
}

var File_request_proto protoreflect.FileDescriptor

var file_request_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22,
	0x58, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x06,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x03, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x03, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x04, 0x62, 0x6f,
	0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c,
	0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xc5, 0x03, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x12, 0x36, 0x0a, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72,
	0x2e, 0x48, 0x54, 0x54, 0x50, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x68, 0x74,
	0x74, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x69, 0x72,
	0x72, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x2d, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04,
	0x6d, 0x65, 0x74, 0x61, 0x1a, 0x4f, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4a, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x69, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x2a, 0x6a, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x07, 0x0a, 0x03, 0x47,
	0x45, 0x54, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x45, 0x41, 0x44, 0x10, 0x01, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x55, 0x54, 0x10,
	0x03, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x41, 0x54, 0x43, 0x48, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x4e,
	0x45, 0x43, 0x54, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x53,
	0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x52, 0x41, 0x43, 0x45, 0x10, 0x08, 0x2a, 0x32, 0x0a,
	0x0b, 0x48, 0x54, 0x54, 0x50, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0b, 0x0a, 0x07,
	0x48, 0x54, 0x54, 0x50, 0x31, 0x5f, 0x30, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x48, 0x54, 0x54,
	0x50, 0x31, 0x5f, 0x31, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x54, 0x54, 0x50, 0x32, 0x10,
	0x02, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_request_proto_rawDescOnce sync.Once
	file_request_proto_rawDescData = file_request_proto_rawDesc
)

func file_request_proto_rawDescGZIP() []byte {
	file_request_proto_rawDescOnce.Do(func() {
		file_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_request_proto_rawDescData)
	})
	return file_request_proto_rawDescData
}

var file_request_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_request_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_request_proto_goTypes = []interface{}{
	(Method)(0),                   // 0: mirror.Method
	(HTTPVersion)(0),              // 1: mirror.HTTPVersion
	(*HeaderValue)(nil),           // 2: mirror.HeaderValue
	(*MetaValue)(nil),             // 3: mirror.MetaValue
	(*Request)(nil),               // 4: mirror.Request
	nil,                           // 5: mirror.Request.HeadersEntry
	nil,                           // 6: mirror.Request.MetaEntry
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_request_proto_depIdxs = []int32{
	7, // 0: mirror.Request.time:type_name -> google.protobuf.Timestamp
	0, // 1: mirror.Request.method:type_name -> mirror.Method
	1, // 2: mirror.Request.http_version:type_name -> mirror.HTTPVersion
	5, // 3: mirror.Request.headers:type_name -> mirror.Request.HeadersEntry
	6, // 4: mirror.Request.meta:type_name -> mirror.Request.MetaEntry
	2, // 5: mirror.Request.HeadersEntry.value:type_name -> mirror.HeaderValue
	3, // 6: mirror.Request.MetaEntry.value:type_name -> mirror.MetaValue
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_request_proto_init() }
func file_request_proto_init() {
	if File_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeaderValue); i {
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
		file_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaValue); i {
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
		file_request_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
	file_request_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*MetaValue_String_)(nil),
		(*MetaValue_Int)(nil),
		(*MetaValue_Bool)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_request_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_request_proto_goTypes,
		DependencyIndexes: file_request_proto_depIdxs,
		EnumInfos:         file_request_proto_enumTypes,
		MessageInfos:      file_request_proto_msgTypes,
	}.Build()
	File_request_proto = out.File
	file_request_proto_rawDesc = nil
	file_request_proto_goTypes = nil
	file_request_proto_depIdxs = nil
}
