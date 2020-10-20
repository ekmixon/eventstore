// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package eventstore

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EventStoreClient is the client API for EventStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventStoreClient interface {
	// Save variable to storage
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error)
	// Load variable from storage
	Load(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error)
	// Delete variable from storage
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type eventStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewEventStoreClient(cc grpc.ClientConnInterface) EventStoreClient {
	return &eventStoreClient{cc}
}

func (c *eventStoreClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error) {
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, "/eventstore.EventStore/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventStoreClient) Load(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error) {
	out := new(LoadResponse)
	err := c.cc.Invoke(ctx, "/eventstore.EventStore/Load", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventStoreClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/eventstore.EventStore/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventStoreServer is the server API for EventStore service.
// All implementations must embed UnimplementedEventStoreServer
// for forward compatibility
type EventStoreServer interface {
	// Save variable to storage
	Save(context.Context, *SaveRequest) (*SaveResponse, error)
	// Load variable from storage
	Load(context.Context, *LoadRequest) (*LoadResponse, error)
	// Delete variable from storage
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	mustEmbedUnimplementedEventStoreServer()
}

// UnimplementedEventStoreServer must be embedded to have forward compatible implementations.
type UnimplementedEventStoreServer struct {
}

func (UnimplementedEventStoreServer) Save(context.Context, *SaveRequest) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedEventStoreServer) Load(context.Context, *LoadRequest) (*LoadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Load not implemented")
}
func (UnimplementedEventStoreServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEventStoreServer) mustEmbedUnimplementedEventStoreServer() {}

// UnsafeEventStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventStoreServer will
// result in compilation errors.
type UnsafeEventStoreServer interface {
	mustEmbedUnimplementedEventStoreServer()
}

func RegisterEventStoreServer(s *grpc.Server, srv EventStoreServer) {
	s.RegisterService(&_EventStore_serviceDesc, srv)
}

func _EventStore_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventStoreServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventstore.EventStore/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventStoreServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventStore_Load_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventStoreServer).Load(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventstore.EventStore/Load",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventStoreServer).Load(ctx, req.(*LoadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventStore_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventStoreServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventstore.EventStore/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventStoreServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eventstore.EventStore",
	HandlerType: (*EventStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _EventStore_Save_Handler,
		},
		{
			MethodName: "Load",
			Handler:    _EventStore_Load_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EventStore_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/eventstore/api/eventstore/eventstore.proto",
}
