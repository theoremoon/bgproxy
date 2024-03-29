// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bgproxy.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Result struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_0deaf48897fdb7c5, []int{0}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type Target struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ExpectedStatus       int32    `protobuf:"varint,2,opt,name=expected_status,json=expectedStatus,proto3" json:"expected_status,omitempty"`
	HealthcheckInterval  int32    `protobuf:"varint,3,opt,name=healthcheck_interval,json=healthcheckInterval,proto3" json:"healthcheck_interval,omitempty"`
	UnhealthyLimit       int32    `protobuf:"varint,4,opt,name=unhealthy_limit,json=unhealthyLimit,proto3" json:"unhealthy_limit,omitempty"`
	WaitingTime          int32    `protobuf:"varint,5,opt,name=waiting_time,json=waitingTime,proto3" json:"waiting_time,omitempty"`
	StopCommand          string   `protobuf:"bytes,6,opt,name=stop_command,json=stopCommand,proto3" json:"stop_command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Target) Reset()         { *m = Target{} }
func (m *Target) String() string { return proto.CompactTextString(m) }
func (*Target) ProtoMessage()    {}
func (*Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_0deaf48897fdb7c5, []int{1}
}

func (m *Target) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Target.Unmarshal(m, b)
}
func (m *Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Target.Marshal(b, m, deterministic)
}
func (m *Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Target.Merge(m, src)
}
func (m *Target) XXX_Size() int {
	return xxx_messageInfo_Target.Size(m)
}
func (m *Target) XXX_DiscardUnknown() {
	xxx_messageInfo_Target.DiscardUnknown(m)
}

var xxx_messageInfo_Target proto.InternalMessageInfo

func (m *Target) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Target) GetExpectedStatus() int32 {
	if m != nil {
		return m.ExpectedStatus
	}
	return 0
}

func (m *Target) GetHealthcheckInterval() int32 {
	if m != nil {
		return m.HealthcheckInterval
	}
	return 0
}

func (m *Target) GetUnhealthyLimit() int32 {
	if m != nil {
		return m.UnhealthyLimit
	}
	return 0
}

func (m *Target) GetWaitingTime() int32 {
	if m != nil {
		return m.WaitingTime
	}
	return 0
}

func (m *Target) GetStopCommand() string {
	if m != nil {
		return m.StopCommand
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_0deaf48897fdb7c5, []int{2}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Result)(nil), "pb.Result")
	proto.RegisterType((*Target)(nil), "pb.Target")
	proto.RegisterType((*Empty)(nil), "pb.Empty")
}

func init() { proto.RegisterFile("bgproxy.proto", fileDescriptor_0deaf48897fdb7c5) }

var fileDescriptor_0deaf48897fdb7c5 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0xe9, 0xe6, 0xea, 0xf6, 0x4d, 0xa7, 0x44, 0x0f, 0x65, 0xa7, 0xad, 0x17, 0x77, 0x1a,
	0xa8, 0xff, 0x40, 0x91, 0x21, 0x78, 0x90, 0x6e, 0xf7, 0x92, 0x76, 0x1f, 0x5d, 0x58, 0xd2, 0x84,
	0xf4, 0xeb, 0x5c, 0xf1, 0xf7, 0xfa, 0x3f, 0x24, 0x69, 0x27, 0x0a, 0xde, 0x3e, 0x9e, 0xf7, 0xe1,
	0x85, 0xbc, 0x81, 0xcb, 0xac, 0x30, 0x56, 0x1f, 0x9b, 0xa5, 0xb1, 0x9a, 0x34, 0xeb, 0x99, 0x2c,
	0x9e, 0x42, 0x98, 0x60, 0x55, 0x4b, 0x62, 0xd7, 0xd0, 0x57, 0x55, 0x11, 0x05, 0xb3, 0x60, 0x31,
	0x4a, 0xdc, 0x19, 0x7f, 0x05, 0x10, 0x6e, 0xb8, 0x2d, 0xd0, 0x87, 0xb5, 0x95, 0xa7, 0xb0, 0xb6,
	0x92, 0xdd, 0xc1, 0x15, 0x1e, 0x0d, 0xe6, 0x84, 0xdb, 0xb4, 0x22, 0x4e, 0x75, 0x15, 0xf5, 0x66,
	0xc1, 0x62, 0x90, 0x4c, 0x4e, 0x78, 0xed, 0x29, 0xbb, 0x87, 0xdb, 0x1d, 0x72, 0x49, 0xbb, 0x7c,
	0x87, 0xf9, 0x3e, 0x15, 0x25, 0xa1, 0x3d, 0x70, 0x19, 0xf5, 0xbd, 0x7d, 0xf3, 0x2b, 0x7b, 0xed,
	0x22, 0xd7, 0x5d, 0x97, 0x6d, 0xd0, 0xa4, 0x52, 0x28, 0x41, 0xd1, 0x59, 0xdb, 0xfd, 0x83, 0xdf,
	0x1c, 0x65, 0x73, 0xb8, 0xf8, 0xe0, 0x82, 0x44, 0x59, 0xa4, 0x24, 0x14, 0x46, 0x03, 0x6f, 0x8d,
	0x3b, 0xb6, 0x11, 0x0a, 0x9d, 0x52, 0x91, 0x36, 0x69, 0xae, 0x95, 0xe2, 0xe5, 0x36, 0x0a, 0xfd,
	0x13, 0xc6, 0x8e, 0x3d, 0xb7, 0x28, 0x3e, 0x87, 0xc1, 0x8b, 0x32, 0xd4, 0x3c, 0x7c, 0xc2, 0xe4,
	0x69, 0xf5, 0xee, 0x16, 0x5a, 0xa3, 0x3d, 0x88, 0x1c, 0x59, 0x0c, 0xc3, 0x35, 0xd2, 0xca, 0x22,
	0x96, 0x0c, 0x96, 0x26, 0x5b, 0xb6, 0x7b, 0x4c, 0xfd, 0xdd, 0x0d, 0x37, 0x87, 0x61, 0xa2, 0xa5,
	0xcc, 0x78, 0xbe, 0x67, 0x23, 0xc7, 0x7d, 0xd9, 0x1f, 0x25, 0x86, 0xd1, 0x0a, 0xa9, 0x1b, 0xe4,
	0x7f, 0x27, 0x0b, 0xfd, 0xa7, 0x3c, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0xaf, 0x73, 0x09, 0xc9,
	0xa5, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BGProxyServiceClient is the client API for BGProxyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BGProxyServiceClient interface {
	SetGreen(ctx context.Context, in *Target, opts ...grpc.CallOption) (*Result, error)
	Rollback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error)
	GetStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error)
}

type bGProxyServiceClient struct {
	cc *grpc.ClientConn
}

func NewBGProxyServiceClient(cc *grpc.ClientConn) BGProxyServiceClient {
	return &bGProxyServiceClient{cc}
}

func (c *bGProxyServiceClient) SetGreen(ctx context.Context, in *Target, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/pb.BGProxyService/SetGreen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bGProxyServiceClient) Rollback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/pb.BGProxyService/Rollback", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bGProxyServiceClient) GetStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/pb.BGProxyService/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BGProxyServiceServer is the server API for BGProxyService service.
type BGProxyServiceServer interface {
	SetGreen(context.Context, *Target) (*Result, error)
	Rollback(context.Context, *Empty) (*Result, error)
	GetStatus(context.Context, *Empty) (*Result, error)
}

// UnimplementedBGProxyServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBGProxyServiceServer struct {
}

func (*UnimplementedBGProxyServiceServer) SetGreen(ctx context.Context, req *Target) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGreen not implemented")
}
func (*UnimplementedBGProxyServiceServer) Rollback(ctx context.Context, req *Empty) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rollback not implemented")
}
func (*UnimplementedBGProxyServiceServer) GetStatus(ctx context.Context, req *Empty) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}

func RegisterBGProxyServiceServer(s *grpc.Server, srv BGProxyServiceServer) {
	s.RegisterService(&_BGProxyService_serviceDesc, srv)
}

func _BGProxyService_SetGreen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Target)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BGProxyServiceServer).SetGreen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BGProxyService/SetGreen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BGProxyServiceServer).SetGreen(ctx, req.(*Target))
	}
	return interceptor(ctx, in, info, handler)
}

func _BGProxyService_Rollback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BGProxyServiceServer).Rollback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BGProxyService/Rollback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BGProxyServiceServer).Rollback(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BGProxyService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BGProxyServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BGProxyService/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BGProxyServiceServer).GetStatus(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _BGProxyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BGProxyService",
	HandlerType: (*BGProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetGreen",
			Handler:    _BGProxyService_SetGreen_Handler,
		},
		{
			MethodName: "Rollback",
			Handler:    _BGProxyService_Rollback_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _BGProxyService_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bgproxy.proto",
}
