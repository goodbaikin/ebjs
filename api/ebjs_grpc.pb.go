// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/ebjs.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Encoder_Encode_FullMethodName = "/Encoder/Encode"
)

// EncoderClient is the client API for Encoder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EncoderClient interface {
	Encode(ctx context.Context, in *EncodeRequest, opts ...grpc.CallOption) (Encoder_EncodeClient, error)
}

type encoderClient struct {
	cc grpc.ClientConnInterface
}

func NewEncoderClient(cc grpc.ClientConnInterface) EncoderClient {
	return &encoderClient{cc}
}

func (c *encoderClient) Encode(ctx context.Context, in *EncodeRequest, opts ...grpc.CallOption) (Encoder_EncodeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Encoder_ServiceDesc.Streams[0], Encoder_Encode_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &encoderEncodeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Encoder_EncodeClient interface {
	Recv() (*EncodeProgress, error)
	grpc.ClientStream
}

type encoderEncodeClient struct {
	grpc.ClientStream
}

func (x *encoderEncodeClient) Recv() (*EncodeProgress, error) {
	m := new(EncodeProgress)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EncoderServer is the server API for Encoder service.
// All implementations must embed UnimplementedEncoderServer
// for forward compatibility
type EncoderServer interface {
	Encode(*EncodeRequest, Encoder_EncodeServer) error
	mustEmbedUnimplementedEncoderServer()
}

// UnimplementedEncoderServer must be embedded to have forward compatible implementations.
type UnimplementedEncoderServer struct {
}

func (UnimplementedEncoderServer) Encode(*EncodeRequest, Encoder_EncodeServer) error {
	return status.Errorf(codes.Unimplemented, "method Encode not implemented")
}
func (UnimplementedEncoderServer) mustEmbedUnimplementedEncoderServer() {}

// UnsafeEncoderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EncoderServer will
// result in compilation errors.
type UnsafeEncoderServer interface {
	mustEmbedUnimplementedEncoderServer()
}

func RegisterEncoderServer(s grpc.ServiceRegistrar, srv EncoderServer) {
	s.RegisterService(&Encoder_ServiceDesc, srv)
}

func _Encoder_Encode_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EncodeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EncoderServer).Encode(m, &encoderEncodeServer{stream})
}

type Encoder_EncodeServer interface {
	Send(*EncodeProgress) error
	grpc.ServerStream
}

type encoderEncodeServer struct {
	grpc.ServerStream
}

func (x *encoderEncodeServer) Send(m *EncodeProgress) error {
	return x.ServerStream.SendMsg(m)
}

// Encoder_ServiceDesc is the grpc.ServiceDesc for Encoder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Encoder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Encoder",
	HandlerType: (*EncoderServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Encode",
			Handler:       _Encoder_Encode_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/ebjs.proto",
}
