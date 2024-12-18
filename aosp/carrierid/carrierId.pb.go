// Copyright (C) 2019 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.19.6
// source: carrierId.proto

package carrierid

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

// A complete list of carriers
type CarrierList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A collection of carriers. one entry for one carrier.
	CarrierId []*CarrierId `protobuf:"bytes,1,rep,name=carrier_id,json=carrierId" json:"carrier_id,omitempty"`
	// Version number of current carrier list
	Version *int32 `protobuf:"varint,2,opt,name=version" json:"version,omitempty"`
}

func (x *CarrierList) Reset() {
	*x = CarrierList{}
	mi := &file_carrierId_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CarrierList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarrierList) ProtoMessage() {}

func (x *CarrierList) ProtoReflect() protoreflect.Message {
	mi := &file_carrierId_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarrierList.ProtoReflect.Descriptor instead.
func (*CarrierList) Descriptor() ([]byte, []int) {
	return file_carrierId_proto_rawDescGZIP(), []int{0}
}

func (x *CarrierList) GetCarrierId() []*CarrierId {
	if x != nil {
		return x.CarrierId
	}
	return nil
}

func (x *CarrierList) GetVersion() int32 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return 0
}

// CarrierId is the unique representation of a carrier in CID table.
type CarrierId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// [Optional] A unique canonical number designated to a carrier.
	CanonicalId *int32 `protobuf:"varint,1,opt,name=canonical_id,json=canonicalId" json:"canonical_id,omitempty"`
	// [Optional] A user-friendly carrier name (not localized).
	CarrierName *string `protobuf:"bytes,2,opt,name=carrier_name,json=carrierName" json:"carrier_name,omitempty"`
	// [Optional] Carrier attributes to match a carrier. At least one value is required.
	CarrierAttribute []*CarrierAttribute `protobuf:"bytes,3,rep,name=carrier_attribute,json=carrierAttribute" json:"carrier_attribute,omitempty"`
	// [Optional] A unique canonical number to represent its parent carrier. The parent-child
	// relationship can be used to differentiate a single carrier by different networks,
	// by prepaid v.s. postpaid  or even by 4G v.s. 3G plan.
	ParentCanonicalId *int32 `protobuf:"varint,4,opt,name=parent_canonical_id,json=parentCanonicalId" json:"parent_canonical_id,omitempty"`
}

func (x *CarrierId) Reset() {
	*x = CarrierId{}
	mi := &file_carrierId_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CarrierId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarrierId) ProtoMessage() {}

func (x *CarrierId) ProtoReflect() protoreflect.Message {
	mi := &file_carrierId_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarrierId.ProtoReflect.Descriptor instead.
func (*CarrierId) Descriptor() ([]byte, []int) {
	return file_carrierId_proto_rawDescGZIP(), []int{1}
}

func (x *CarrierId) GetCanonicalId() int32 {
	if x != nil && x.CanonicalId != nil {
		return *x.CanonicalId
	}
	return 0
}

func (x *CarrierId) GetCarrierName() string {
	if x != nil && x.CarrierName != nil {
		return *x.CarrierName
	}
	return ""
}

func (x *CarrierId) GetCarrierAttribute() []*CarrierAttribute {
	if x != nil {
		return x.CarrierAttribute
	}
	return nil
}

func (x *CarrierId) GetParentCanonicalId() int32 {
	if x != nil && x.ParentCanonicalId != nil {
		return *x.ParentCanonicalId
	}
	return 0
}

// Attributes used to match a carrier.
// For each field within this message:
//   - if not set, the attribute is ignored;
//   - if set, the device must have one of the specified values to match.
//
// Match is based on AND between any field that is set and OR for values within a repeated field.
type CarrierAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// [Optional] The MCC and MNC that map to this carrier. At least one value is required.
	MccmncTuple []string `protobuf:"bytes,1,rep,name=mccmnc_tuple,json=mccmncTuple" json:"mccmnc_tuple,omitempty"`
	// [Optional] Prefix of IMSI (International Mobile Subscriber Identity) in
	// decimal format. Some digits can be replaced with "x" symbols matching any digit.
	// Sample values: 20404794, 21670xx2xxx.
	ImsiPrefixXpattern []string `protobuf:"bytes,2,rep,name=imsi_prefix_xpattern,json=imsiPrefixXpattern" json:"imsi_prefix_xpattern,omitempty"`
	// [Optional] The Service Provider Name. Read from subscription EF_SPN.
	// Sample values: C Spire, LeclercMobile
	Spn []string `protobuf:"bytes,3,rep,name=spn" json:"spn,omitempty"`
	// [Optional] PLMN network name. Read from subscription EF_PNN.
	// Sample values:
	Plmn []string `protobuf:"bytes,4,rep,name=plmn" json:"plmn,omitempty"`
	// [Optional] Group Identifier Level1 for a GSM phone. Read from subscription EF_GID1.
	// Sample values: 6D, BAE0000000000000
	Gid1 []string `protobuf:"bytes,5,rep,name=gid1" json:"gid1,omitempty"`
	// [Optional] Group Identifier Level2 for a GSM phone. Read from subscription EF_GID2.
	// Sample values: 6D, BAE0000000000000
	Gid2 []string `protobuf:"bytes,6,rep,name=gid2" json:"gid2,omitempty"`
	// [Optional] The Access Point Name, corresponding to "apn" field returned by
	// "content://telephony/carriers/preferapn" on device.
	// Sample values: fast.t-mobile.com, internet
	PreferredApn []string `protobuf:"bytes,7,rep,name=preferred_apn,json=preferredApn" json:"preferred_apn,omitempty"`
	// [Optional] Prefix of Integrated Circuit Card Identifier. Read from subscription EF_ICCID.
	// Sample values: 894430, 894410
	IccidPrefix []string `protobuf:"bytes,8,rep,name=iccid_prefix,json=iccidPrefix" json:"iccid_prefix,omitempty"`
	// [Optional] Carrier Privilege Access Rule in hex string.
	// Sample values: 61ed377e85d386a8dfee6b864bd85b0bfaa5af88
	PrivilegeAccessRule []string `protobuf:"bytes,9,rep,name=privilege_access_rule,json=privilegeAccessRule" json:"privilege_access_rule,omitempty"`
}

func (x *CarrierAttribute) Reset() {
	*x = CarrierAttribute{}
	mi := &file_carrierId_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CarrierAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarrierAttribute) ProtoMessage() {}

func (x *CarrierAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_carrierId_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarrierAttribute.ProtoReflect.Descriptor instead.
func (*CarrierAttribute) Descriptor() ([]byte, []int) {
	return file_carrierId_proto_rawDescGZIP(), []int{2}
}

func (x *CarrierAttribute) GetMccmncTuple() []string {
	if x != nil {
		return x.MccmncTuple
	}
	return nil
}

func (x *CarrierAttribute) GetImsiPrefixXpattern() []string {
	if x != nil {
		return x.ImsiPrefixXpattern
	}
	return nil
}

func (x *CarrierAttribute) GetSpn() []string {
	if x != nil {
		return x.Spn
	}
	return nil
}

func (x *CarrierAttribute) GetPlmn() []string {
	if x != nil {
		return x.Plmn
	}
	return nil
}

func (x *CarrierAttribute) GetGid1() []string {
	if x != nil {
		return x.Gid1
	}
	return nil
}

func (x *CarrierAttribute) GetGid2() []string {
	if x != nil {
		return x.Gid2
	}
	return nil
}

func (x *CarrierAttribute) GetPreferredApn() []string {
	if x != nil {
		return x.PreferredApn
	}
	return nil
}

func (x *CarrierAttribute) GetIccidPrefix() []string {
	if x != nil {
		return x.IccidPrefix
	}
	return nil
}

func (x *CarrierAttribute) GetPrivilegeAccessRule() []string {
	if x != nil {
		return x.PrivilegeAccessRule
	}
	return nil
}

var File_carrierId_proto protoreflect.FileDescriptor

var file_carrierId_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x15, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x68, 0x0a, 0x0b, 0x43, 0x61, 0x72, 0x72,
	0x69, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x0a, 0x63, 0x61, 0x72, 0x72, 0x69,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x61,
	0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x52, 0x09, 0x63,
	0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0xd7, 0x01, 0x0a, 0x09, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61,
	0x6c, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x72, 0x72, 0x69,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x54, 0x0a, 0x11, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65,
	0x72, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x27, 0x2e, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65,
	0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x52, 0x10, 0x63, 0x61, 0x72, 0x72,
	0x69, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x13,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x22, 0xb1, 0x02, 0x0a,
	0x10, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x63, 0x63, 0x6d, 0x6e, 0x63, 0x5f, 0x74, 0x75, 0x70, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x63, 0x63, 0x6d, 0x6e, 0x63, 0x54,
	0x75, 0x70, 0x6c, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x69, 0x6d, 0x73, 0x69, 0x5f, 0x70, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x5f, 0x78, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x12, 0x69, 0x6d, 0x73, 0x69, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x58, 0x70,
	0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x70, 0x6e, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x03, 0x73, 0x70, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6c, 0x6d, 0x6e,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6c, 0x6d, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x67, 0x69, 0x64, 0x31, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x67, 0x69, 0x64, 0x31,
	0x12, 0x12, 0x0a, 0x04, 0x67, 0x69, 0x64, 0x32, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x67, 0x69, 0x64, 0x32, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65,
	0x64, 0x5f, 0x61, 0x70, 0x6e, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x41, 0x70, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x63, 0x63,
	0x69, 0x64, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0b, 0x69, 0x63, 0x63, 0x69, 0x64, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x32, 0x0a, 0x15,
	0x70, 0x72, 0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x13, 0x70, 0x72, 0x69,
	0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x75, 0x6c, 0x65,
	0x42, 0x31, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x68,
	0x6f, 0x6e, 0x79, 0x42, 0x0e, 0x43, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x49, 0x64, 0x50, 0x72,
	0x6f, 0x74, 0x6f,
}

var (
	file_carrierId_proto_rawDescOnce sync.Once
	file_carrierId_proto_rawDescData = file_carrierId_proto_rawDesc
)

func file_carrierId_proto_rawDescGZIP() []byte {
	file_carrierId_proto_rawDescOnce.Do(func() {
		file_carrierId_proto_rawDescData = protoimpl.X.CompressGZIP(file_carrierId_proto_rawDescData)
	})
	return file_carrierId_proto_rawDescData
}

var file_carrierId_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_carrierId_proto_goTypes = []any{
	(*CarrierList)(nil),      // 0: carrierIdentification.CarrierList
	(*CarrierId)(nil),        // 1: carrierIdentification.CarrierId
	(*CarrierAttribute)(nil), // 2: carrierIdentification.CarrierAttribute
}
var file_carrierId_proto_depIdxs = []int32{
	1, // 0: carrierIdentification.CarrierList.carrier_id:type_name -> carrierIdentification.CarrierId
	2, // 1: carrierIdentification.CarrierId.carrier_attribute:type_name -> carrierIdentification.CarrierAttribute
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_carrierId_proto_init() }
func file_carrierId_proto_init() {
	if File_carrierId_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_carrierId_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_carrierId_proto_goTypes,
		DependencyIndexes: file_carrierId_proto_depIdxs,
		MessageInfos:      file_carrierId_proto_msgTypes,
	}.Build()
	File_carrierId_proto = out.File
	file_carrierId_proto_rawDesc = nil
	file_carrierId_proto_goTypes = nil
	file_carrierId_proto_depIdxs = nil
}
