// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: createsmcontext.proto

package create_sm_context_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SendDataForCreateSmContext_SendDataForCreateSmContext_FullMethodName = "/create_sm_context.SendDataForCreateSmContext/SendDataForCreateSmContext"
)

// SendDataForCreateSmContextClient is the client API for SendDataForCreateSmContext service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SendDataForCreateSmContextClient interface {
	SendDataForCreateSmContext(ctx context.Context, in *CreateSmContextDataFromNasMod, opts ...grpc.CallOption) (*CreateSmContextRespToNasMod, error)
}

type sendDataForCreateSmContextClient struct {
	cc grpc.ClientConnInterface
}

func NewSendDataForCreateSmContextClient(cc grpc.ClientConnInterface) SendDataForCreateSmContextClient {
	return &sendDataForCreateSmContextClient{cc}
}

func (c *sendDataForCreateSmContextClient) SendDataForCreateSmContext(ctx context.Context, in *CreateSmContextDataFromNasMod, opts ...grpc.CallOption) (*CreateSmContextRespToNasMod, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateSmContextRespToNasMod)
	err := c.cc.Invoke(ctx, SendDataForCreateSmContext_SendDataForCreateSmContext_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendDataForCreateSmContextServer is the server API for SendDataForCreateSmContext service.
// All implementations must embed UnimplementedSendDataForCreateSmContextServer
// for forward compatibility.
type SendDataForCreateSmContextServer interface {
	SendDataForCreateSmContext(context.Context, *CreateSmContextDataFromNasMod) (*CreateSmContextRespToNasMod, error)
	mustEmbedUnimplementedSendDataForCreateSmContextServer()
}

// UnimplementedSendDataForCreateSmContextServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSendDataForCreateSmContextServer struct{}

func (UnimplementedSendDataForCreateSmContextServer) SendDataForCreateSmContext(context.Context, *CreateSmContextDataFromNasMod) (*CreateSmContextRespToNasMod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDataForCreateSmContext not implemented")
}
func (UnimplementedSendDataForCreateSmContextServer) mustEmbedUnimplementedSendDataForCreateSmContextServer() {
}
func (UnimplementedSendDataForCreateSmContextServer) testEmbeddedByValue() {}

// UnsafeSendDataForCreateSmContextServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SendDataForCreateSmContextServer will
// result in compilation errors.
type UnsafeSendDataForCreateSmContextServer interface {
	mustEmbedUnimplementedSendDataForCreateSmContextServer()
}

func RegisterSendDataForCreateSmContextServer(s grpc.ServiceRegistrar, srv SendDataForCreateSmContextServer) {
	// If the following call pancis, it indicates UnimplementedSendDataForCreateSmContextServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SendDataForCreateSmContext_ServiceDesc, srv)
}

func _SendDataForCreateSmContext_SendDataForCreateSmContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSmContextDataFromNasMod)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendDataForCreateSmContextServer).SendDataForCreateSmContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SendDataForCreateSmContext_SendDataForCreateSmContext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendDataForCreateSmContextServer).SendDataForCreateSmContext(ctx, req.(*CreateSmContextDataFromNasMod))
	}
	return interceptor(ctx, in, info, handler)
}

// SendDataForCreateSmContext_ServiceDesc is the grpc.ServiceDesc for SendDataForCreateSmContext service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SendDataForCreateSmContext_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "create_sm_context.SendDataForCreateSmContext",
	HandlerType: (*SendDataForCreateSmContextServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendDataForCreateSmContext",
			Handler:    _SendDataForCreateSmContext_SendDataForCreateSmContext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "createsmcontext.proto",
}
