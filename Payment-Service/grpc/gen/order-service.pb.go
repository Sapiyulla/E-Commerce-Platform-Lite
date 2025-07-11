// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: proto/order-service.proto

package gen

import (
	common "payment-service/grpc/gen/common"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Order struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ItemId        string                 `protobuf:"bytes,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_proto_order_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_proto_order_service_proto_rawDescGZIP(), []int{0}
}

func (x *Order) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Order) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

var File_proto_order_service_proto protoreflect.FileDescriptor

const file_proto_order_service_proto_rawDesc = "" +
	"\n" +
	"\x19proto/order-service.proto\x12\forderservice\x1a\x15proto/common/id.proto\"0\n" +
	"\x05Order\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x17\n" +
	"\aitem_id\x18\x02 \x01(\tR\x06itemId2k\n" +
	"\fOrderService\x12.\n" +
	"\vCreateOrder\x12\n" +
	".common.ID\x1a\x13.orderservice.Order\x12+\n" +
	"\bPayOrder\x12\n" +
	".common.ID\x1a\x13.orderservice.OrderB\n" +
	"Z\bgrpc/genb\x06proto3"

var (
	file_proto_order_service_proto_rawDescOnce sync.Once
	file_proto_order_service_proto_rawDescData []byte
)

func file_proto_order_service_proto_rawDescGZIP() []byte {
	file_proto_order_service_proto_rawDescOnce.Do(func() {
		file_proto_order_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_order_service_proto_rawDesc), len(file_proto_order_service_proto_rawDesc)))
	})
	return file_proto_order_service_proto_rawDescData
}

var file_proto_order_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_order_service_proto_goTypes = []any{
	(*Order)(nil),     // 0: orderservice.Order
	(*common.ID)(nil), // 1: common.ID
}
var file_proto_order_service_proto_depIdxs = []int32{
	1, // 0: orderservice.OrderService.CreateOrder:input_type -> common.ID
	1, // 1: orderservice.OrderService.PayOrder:input_type -> common.ID
	0, // 2: orderservice.OrderService.CreateOrder:output_type -> orderservice.Order
	0, // 3: orderservice.OrderService.PayOrder:output_type -> orderservice.Order
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_order_service_proto_init() }
func file_proto_order_service_proto_init() {
	if File_proto_order_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_order_service_proto_rawDesc), len(file_proto_order_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_order_service_proto_goTypes,
		DependencyIndexes: file_proto_order_service_proto_depIdxs,
		MessageInfos:      file_proto_order_service_proto_msgTypes,
	}.Build()
	File_proto_order_service_proto = out.File
	file_proto_order_service_proto_goTypes = nil
	file_proto_order_service_proto_depIdxs = nil
}
