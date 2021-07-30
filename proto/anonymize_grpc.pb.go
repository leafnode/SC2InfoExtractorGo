// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// AnonymizeServiceClient is the client API for AnonymizeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnonymizeServiceClient interface {
	GetAnonymizedID(ctx context.Context, in *SendNickname, opts ...grpc.CallOption) (*ReceiveID, error)
}

type anonymizeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAnonymizeServiceClient(cc grpc.ClientConnInterface) AnonymizeServiceClient {
	return &anonymizeServiceClient{cc}
}

func (c *anonymizeServiceClient) GetAnonymizedID(ctx context.Context, in *SendNickname, opts ...grpc.CallOption) (*ReceiveID, error) {
	out := new(ReceiveID)
	err := c.cc.Invoke(ctx, "/AnonymizeService/getAnonymizedID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnonymizeServiceServer is the server API for AnonymizeService service.
// All implementations must embed UnimplementedAnonymizeServiceServer
// for forward compatibility
type AnonymizeServiceServer interface {
	GetAnonymizedID(context.Context, *SendNickname) (*ReceiveID, error)
	mustEmbedUnimplementedAnonymizeServiceServer()
}

// UnimplementedAnonymizeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAnonymizeServiceServer struct {
}

func (UnimplementedAnonymizeServiceServer) GetAnonymizedID(context.Context, *SendNickname) (*ReceiveID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAnonymizedID not implemented")
}
func (UnimplementedAnonymizeServiceServer) mustEmbedUnimplementedAnonymizeServiceServer() {}

// UnsafeAnonymizeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnonymizeServiceServer will
// result in compilation errors.
type UnsafeAnonymizeServiceServer interface {
	mustEmbedUnimplementedAnonymizeServiceServer()
}

func RegisterAnonymizeServiceServer(s grpc.ServiceRegistrar, srv AnonymizeServiceServer) {
	s.RegisterService(&AnonymizeService_ServiceDesc, srv)
}

func _AnonymizeService_GetAnonymizedID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendNickname)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnonymizeServiceServer).GetAnonymizedID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AnonymizeService/getAnonymizedID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnonymizeServiceServer).GetAnonymizedID(ctx, req.(*SendNickname))
	}
	return interceptor(ctx, in, info, handler)
}

// AnonymizeService_ServiceDesc is the grpc.ServiceDesc for AnonymizeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AnonymizeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AnonymizeService",
	HandlerType: (*AnonymizeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getAnonymizedID",
			Handler:    _AnonymizeService_GetAnonymizedID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "anonymize.proto",
}
