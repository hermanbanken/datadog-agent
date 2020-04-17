// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/security/api/api.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	pkg/security/api/api.proto

It has these top-level messages:
	GetParams
	SecurityEventMessage
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetParams struct {
}

func (m *GetParams) Reset()                    { *m = GetParams{} }
func (m *GetParams) String() string            { return proto.CompactTextString(m) }
func (*GetParams) ProtoMessage()               {}
func (*GetParams) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SecurityEventMessage struct {
	RuleName string `protobuf:"bytes,1,opt,name=RuleName,json=ruleName" json:"RuleName,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=Data,json=data,proto3" json:"Data,omitempty"`
}

func (m *SecurityEventMessage) Reset()                    { *m = SecurityEventMessage{} }
func (m *SecurityEventMessage) String() string            { return proto.CompactTextString(m) }
func (*SecurityEventMessage) ProtoMessage()               {}
func (*SecurityEventMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SecurityEventMessage) GetRuleName() string {
	if m != nil {
		return m.RuleName
	}
	return ""
}

func (m *SecurityEventMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*GetParams)(nil), "api.GetParams")
	proto.RegisterType((*SecurityEventMessage)(nil), "api.SecurityEventMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SecurityModule service

type SecurityModuleClient interface {
	GetEvents(ctx context.Context, in *GetParams, opts ...grpc.CallOption) (SecurityModule_GetEventsClient, error)
}

type securityModuleClient struct {
	cc *grpc.ClientConn
}

func NewSecurityModuleClient(cc *grpc.ClientConn) SecurityModuleClient {
	return &securityModuleClient{cc}
}

func (c *securityModuleClient) GetEvents(ctx context.Context, in *GetParams, opts ...grpc.CallOption) (SecurityModule_GetEventsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SecurityModule_serviceDesc.Streams[0], c.cc, "/api.SecurityModule/GetEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &securityModuleGetEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SecurityModule_GetEventsClient interface {
	Recv() (*SecurityEventMessage, error)
	grpc.ClientStream
}

type securityModuleGetEventsClient struct {
	grpc.ClientStream
}

func (x *securityModuleGetEventsClient) Recv() (*SecurityEventMessage, error) {
	m := new(SecurityEventMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SecurityModule service

type SecurityModuleServer interface {
	GetEvents(*GetParams, SecurityModule_GetEventsServer) error
}

func RegisterSecurityModuleServer(s *grpc.Server, srv SecurityModuleServer) {
	s.RegisterService(&_SecurityModule_serviceDesc, srv)
}

func _SecurityModule_GetEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetParams)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SecurityModuleServer).GetEvents(m, &securityModuleGetEventsServer{stream})
}

type SecurityModule_GetEventsServer interface {
	Send(*SecurityEventMessage) error
	grpc.ServerStream
}

type securityModuleGetEventsServer struct {
	grpc.ServerStream
}

func (x *securityModuleGetEventsServer) Send(m *SecurityEventMessage) error {
	return x.ServerStream.SendMsg(m)
}

var _SecurityModule_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.SecurityModule",
	HandlerType: (*SecurityModuleServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetEvents",
			Handler:       _SecurityModule_GetEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/security/api/api.proto",
}

func init() { proto.RegisterFile("pkg/security/api/api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0xc8, 0x4e, 0xd7,
	0x2f, 0x4e, 0x4d, 0x2e, 0x2d, 0xca, 0x2c, 0xa9, 0xd4, 0x4f, 0x2c, 0xc8, 0x04, 0x61, 0xbd, 0x82,
	0xa2, 0xfc, 0x92, 0x7c, 0x21, 0xe6, 0xc4, 0x82, 0x4c, 0x25, 0x6e, 0x2e, 0x4e, 0xf7, 0xd4, 0x92,
	0x80, 0xc4, 0xa2, 0xc4, 0xdc, 0x62, 0x25, 0x37, 0x2e, 0x91, 0x60, 0xa8, 0x5a, 0xd7, 0xb2, 0xd4,
	0xbc, 0x12, 0xdf, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x21, 0x29, 0x2e, 0x8e, 0xa0, 0xd2, 0x9c,
	0x54, 0xbf, 0xc4, 0xdc, 0x54, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x8e, 0x22, 0x28, 0x5f,
	0x48, 0x88, 0x8b, 0xc5, 0x25, 0xb1, 0x24, 0x51, 0x82, 0x49, 0x81, 0x51, 0x83, 0x27, 0x88, 0x25,
	0x25, 0xb1, 0x24, 0xd1, 0xc8, 0x87, 0x8b, 0x0f, 0x66, 0x8e, 0x6f, 0x7e, 0x4a, 0x69, 0x4e, 0xaa,
	0x90, 0x15, 0xd8, 0x1a, 0xb0, 0xa1, 0xc5, 0x42, 0x7c, 0x7a, 0x20, 0x47, 0xc0, 0xad, 0x95, 0x92,
	0x04, 0xf3, 0xb1, 0xd9, 0xac, 0xc4, 0x60, 0xc0, 0x98, 0xc4, 0x06, 0x76, 0xae, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0xc6, 0x16, 0xa8, 0x2b, 0xcc, 0x00, 0x00, 0x00,
}