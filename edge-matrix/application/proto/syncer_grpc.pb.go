// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: application/proto/syncer.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SyncAppClient is the client API for SyncApp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SyncAppClient interface {
	// Returns stream of data beginning specified from
	GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (SyncApp_GetDataClient, error)
	// Returns app peer's status
	GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AppStatus, error)
}

type syncAppClient struct {
	cc grpc.ClientConnInterface
}

func NewSyncAppClient(cc grpc.ClientConnInterface) SyncAppClient {
	return &syncAppClient{cc}
}

func (c *syncAppClient) GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (SyncApp_GetDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &SyncApp_ServiceDesc.Streams[0], "/v1.SyncApp/GetData", opts...)
	if err != nil {
		return nil, err
	}
	x := &syncAppGetDataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SyncApp_GetDataClient interface {
	Recv() (*Data, error)
	grpc.ClientStream
}

type syncAppGetDataClient struct {
	grpc.ClientStream
}

func (x *syncAppGetDataClient) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *syncAppClient) GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AppStatus, error) {
	out := new(AppStatus)
	err := c.cc.Invoke(ctx, "/v1.SyncApp/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncAppServer is the server API for SyncApp service.
// All implementations must embed UnimplementedSyncAppServer
// for forward compatibility
type SyncAppServer interface {
	// Returns stream of data beginning specified from
	GetData(*GetDataRequest, SyncApp_GetDataServer) error
	// Returns app peer's status
	GetStatus(context.Context, *emptypb.Empty) (*AppStatus, error)
	mustEmbedUnimplementedSyncAppServer()
}

// UnimplementedSyncAppServer must be embedded to have forward compatible implementations.
type UnimplementedSyncAppServer struct {
}

func (UnimplementedSyncAppServer) GetData(*GetDataRequest, SyncApp_GetDataServer) error {
	return status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedSyncAppServer) GetStatus(context.Context, *emptypb.Empty) (*AppStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedSyncAppServer) mustEmbedUnimplementedSyncAppServer() {}

// UnsafeSyncAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SyncAppServer will
// result in compilation errors.
type UnsafeSyncAppServer interface {
	mustEmbedUnimplementedSyncAppServer()
}

func RegisterSyncAppServer(s grpc.ServiceRegistrar, srv SyncAppServer) {
	s.RegisterService(&SyncApp_ServiceDesc, srv)
}

func _SyncApp_GetData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SyncAppServer).GetData(m, &syncAppGetDataServer{stream})
}

type SyncApp_GetDataServer interface {
	Send(*Data) error
	grpc.ServerStream
}

type syncAppGetDataServer struct {
	grpc.ServerStream
}

func (x *syncAppGetDataServer) Send(m *Data) error {
	return x.ServerStream.SendMsg(m)
}

func _SyncApp_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncAppServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.SyncApp/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncAppServer).GetStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SyncApp_ServiceDesc is the grpc.ServiceDesc for SyncApp service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SyncApp_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.SyncApp",
	HandlerType: (*SyncAppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _SyncApp_GetStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetData",
			Handler:       _SyncApp_GetData_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "application/proto/syncer.proto",
}
