// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: api/bmi.proto

package api

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

type UserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *UserInfoRequest) Reset() {
	*x = UserInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bmi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfoRequest) ProtoMessage() {}

func (x *UserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_bmi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfoRequest.ProtoReflect.Descriptor instead.
func (*UserInfoRequest) Descriptor() ([]byte, []int) {
	return file_api_bmi_proto_rawDescGZIP(), []int{0}
}

func (x *UserInfoRequest) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type BMIInfoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bmi *BMI `protobuf:"bytes,1,opt,name=bmi,proto3" json:"bmi,omitempty"`
}

func (x *BMIInfoReply) Reset() {
	*x = BMIInfoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bmi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BMIInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BMIInfoReply) ProtoMessage() {}

func (x *BMIInfoReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_bmi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BMIInfoReply.ProtoReflect.Descriptor instead.
func (*BMIInfoReply) Descriptor() ([]byte, []int) {
	return file_api_bmi_proto_rawDescGZIP(), []int{1}
}

func (x *BMIInfoReply) GetBmi() *BMI {
	if x != nil {
		return x.Bmi
	}
	return nil
}

type BMI struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nickname string `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Height   uint64 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Weight   uint64 `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Uid      uint64 `protobuf:"varint,4,opt,name=uid,proto3" json:"uid,omitempty"`
	Bmi      uint64 `protobuf:"varint,5,opt,name=bmi,proto3" json:"bmi,omitempty"`
}

func (x *BMI) Reset() {
	*x = BMI{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_bmi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BMI) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BMI) ProtoMessage() {}

func (x *BMI) ProtoReflect() protoreflect.Message {
	mi := &file_api_bmi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BMI.ProtoReflect.Descriptor instead.
func (*BMI) Descriptor() ([]byte, []int) {
	return file_api_bmi_proto_rawDescGZIP(), []int{2}
}

func (x *BMI) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *BMI) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BMI) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *BMI) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *BMI) GetBmi() uint64 {
	if x != nil {
		return x.Bmi
	}
	return 0
}

var File_api_bmi_proto protoreflect.FileDescriptor

var file_api_bmi_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6d, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x34, 0x22, 0x23, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x2d, 0x0a, 0x0c,
	0x42, 0x4d, 0x49, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x03,
	0x62, 0x6d, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x34, 0x2e, 0x42, 0x4d, 0x49, 0x52, 0x03, 0x62, 0x6d, 0x69, 0x22, 0x75, 0x0a, 0x03, 0x42,
	0x4d, 0x49, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x62, 0x6d, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x62,
	0x6d, 0x69, 0x32, 0x48, 0x0a, 0x0a, 0x42, 0x4d, 0x49, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x3a, 0x0a, 0x07, 0x42, 0x4d, 0x49, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x34, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x34, 0x2e, 0x42, 0x4d,
	0x49, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x34, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_api_bmi_proto_rawDescOnce sync.Once
	file_api_bmi_proto_rawDescData = file_api_bmi_proto_rawDesc
)

func file_api_bmi_proto_rawDescGZIP() []byte {
	file_api_bmi_proto_rawDescOnce.Do(func() {
		file_api_bmi_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_bmi_proto_rawDescData)
	})
	return file_api_bmi_proto_rawDescData
}

var file_api_bmi_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_bmi_proto_goTypes = []interface{}{
	(*UserInfoRequest)(nil), // 0: model4.UserInfoRequest
	(*BMIInfoReply)(nil),    // 1: model4.BMIInfoReply
	(*BMI)(nil),             // 2: model4.BMI
}
var file_api_bmi_proto_depIdxs = []int32{
	2, // 0: model4.BMIInfoReply.bmi:type_name -> model4.BMI
	0, // 1: model4.BMIService.BMIInfo:input_type -> model4.UserInfoRequest
	1, // 2: model4.BMIService.BMIInfo:output_type -> model4.BMIInfoReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_bmi_proto_init() }
func file_api_bmi_proto_init() {
	if File_api_bmi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_bmi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfoRequest); i {
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
		file_api_bmi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BMIInfoReply); i {
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
		file_api_bmi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BMI); i {
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
			RawDescriptor: file_api_bmi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_bmi_proto_goTypes,
		DependencyIndexes: file_api_bmi_proto_depIdxs,
		MessageInfos:      file_api_bmi_proto_msgTypes,
	}.Build()
	File_api_bmi_proto = out.File
	file_api_bmi_proto_rawDesc = nil
	file_api_bmi_proto_goTypes = nil
	file_api_bmi_proto_depIdxs = nil
}
