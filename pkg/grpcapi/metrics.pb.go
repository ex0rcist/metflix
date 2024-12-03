// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: metrics.proto

package grpcapi

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricExchange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Mtype string  `protobuf:"bytes,2,opt,name=mtype,proto3" json:"mtype,omitempty"`
	Delta int64   `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	Value float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"`
	Hash  string  `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *MetricExchange) Reset() {
	*x = MetricExchange{}
	mi := &file_metrics_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetricExchange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricExchange) ProtoMessage() {}

func (x *MetricExchange) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricExchange.ProtoReflect.Descriptor instead.
func (*MetricExchange) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *MetricExchange) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MetricExchange) GetMtype() string {
	if x != nil {
		return x.Mtype
	}
	return ""
}

func (x *MetricExchange) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *MetricExchange) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *MetricExchange) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type BatchUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*MetricExchange `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *BatchUpdateRequest) Reset() {
	*x = BatchUpdateRequest{}
	mi := &file_metrics_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchUpdateRequest) ProtoMessage() {}

func (x *BatchUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchUpdateRequest.ProtoReflect.Descriptor instead.
func (*BatchUpdateRequest) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *BatchUpdateRequest) GetData() []*MetricExchange {
	if x != nil {
		return x.Data
	}
	return nil
}

type BatchUpdateEncryptedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncryptedData []byte `protobuf:"bytes,1,opt,name=encrypted_data,json=encryptedData,proto3" json:"encrypted_data,omitempty"`
}

func (x *BatchUpdateEncryptedRequest) Reset() {
	*x = BatchUpdateEncryptedRequest{}
	mi := &file_metrics_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchUpdateEncryptedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchUpdateEncryptedRequest) ProtoMessage() {}

func (x *BatchUpdateEncryptedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchUpdateEncryptedRequest.ProtoReflect.Descriptor instead.
func (*BatchUpdateEncryptedRequest) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *BatchUpdateEncryptedRequest) GetEncryptedData() []byte {
	if x != nil {
		return x.EncryptedData
	}
	return nil
}

type BatchUpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*MetricExchange `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *BatchUpdateResponse) Reset() {
	*x = BatchUpdateResponse{}
	mi := &file_metrics_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchUpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchUpdateResponse) ProtoMessage() {}

func (x *BatchUpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchUpdateResponse.ProtoReflect.Descriptor instead.
func (*BatchUpdateResponse) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *BatchUpdateResponse) GetData() []*MetricExchange {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_metrics_proto protoreflect.FileDescriptor

var file_metrics_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x22, 0x76, 0x0a, 0x0e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x22, 0x44, 0x0a, 0x12, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69,
	0x78, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x45, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x1b, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x22,
	0x45, 0x0a, 0x13, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xbb, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x4e, 0x0a, 0x0b, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x1e, 0x2e, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x60, 0x0a, 0x14, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x12, 0x27, 0x2e, 0x6d, 0x65, 0x74,
	0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6d, 0x65, 0x74, 0x66, 0x6c, 0x69, 0x78, 0x2e, 0x76, 0x31,
	0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x65, 0x78, 0x30, 0x72, 0x63, 0x69, 0x73, 0x74, 0x2f, 0x6d, 0x65, 0x74, 0x66,
	0x6c, 0x69, 0x78, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_metrics_proto_rawDescOnce sync.Once
	file_metrics_proto_rawDescData = file_metrics_proto_rawDesc
)

func file_metrics_proto_rawDescGZIP() []byte {
	file_metrics_proto_rawDescOnce.Do(func() {
		file_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_metrics_proto_rawDescData)
	})
	return file_metrics_proto_rawDescData
}

var file_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_metrics_proto_goTypes = []any{
	(*MetricExchange)(nil),              // 0: metflix.v1.MetricExchange
	(*BatchUpdateRequest)(nil),          // 1: metflix.v1.BatchUpdateRequest
	(*BatchUpdateEncryptedRequest)(nil), // 2: metflix.v1.BatchUpdateEncryptedRequest
	(*BatchUpdateResponse)(nil),         // 3: metflix.v1.BatchUpdateResponse
}
var file_metrics_proto_depIdxs = []int32{
	0, // 0: metflix.v1.BatchUpdateRequest.data:type_name -> metflix.v1.MetricExchange
	0, // 1: metflix.v1.BatchUpdateResponse.data:type_name -> metflix.v1.MetricExchange
	1, // 2: metflix.v1.Metrics.BatchUpdate:input_type -> metflix.v1.BatchUpdateRequest
	2, // 3: metflix.v1.Metrics.BatchUpdateEncrypted:input_type -> metflix.v1.BatchUpdateEncryptedRequest
	3, // 4: metflix.v1.Metrics.BatchUpdate:output_type -> metflix.v1.BatchUpdateResponse
	3, // 5: metflix.v1.Metrics.BatchUpdateEncrypted:output_type -> metflix.v1.BatchUpdateResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_metrics_proto_init() }
func file_metrics_proto_init() {
	if File_metrics_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metrics_proto_goTypes,
		DependencyIndexes: file_metrics_proto_depIdxs,
		MessageInfos:      file_metrics_proto_msgTypes,
	}.Build()
	File_metrics_proto = out.File
	file_metrics_proto_rawDesc = nil
	file_metrics_proto_goTypes = nil
	file_metrics_proto_depIdxs = nil
}