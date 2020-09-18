// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rvapi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RVerClient is the client API for RVer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RVerClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Overview(ctx context.Context, in *OverviewRequest, opts ...grpc.CallOption) (*OverviewResponse, error)
	Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
}

type rVerClient struct {
	cc grpc.ClientConnInterface
}

func NewRVerClient(cc grpc.ClientConnInterface) RVerClient {
	return &rVerClient{cc}
}

func (c *rVerClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) Overview(ctx context.Context, in *OverviewRequest, opts ...grpc.CallOption) (*OverviewResponse, error) {
	out := new(OverviewResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Overview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error) {
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Report", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RVerServer is the server API for RVer service.
// All implementations must embed UnimplementedRVerServer
// for forward compatibility
type RVerServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Overview(context.Context, *OverviewRequest) (*OverviewResponse, error)
	Report(context.Context, *ReportRequest) (*ReportResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
	mustEmbedUnimplementedRVerServer()
}

// UnimplementedRVerServer must be embedded to have forward compatible implementations.
type UnimplementedRVerServer struct {
}

func (*UnimplementedRVerServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedRVerServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedRVerServer) Overview(context.Context, *OverviewRequest) (*OverviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Overview not implemented")
}
func (*UnimplementedRVerServer) Report(context.Context, *ReportRequest) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (*UnimplementedRVerServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedRVerServer) Vote(context.Context, *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (*UnimplementedRVerServer) mustEmbedUnimplementedRVerServer() {}

func RegisterRVerServer(s *grpc.Server, srv RVerServer) {
	s.RegisterService(&_RVer_serviceDesc, srv)
}

func _RVer_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_Overview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OverviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Overview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Overview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Overview(ctx, req.(*OverviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Report",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Report(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RVer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rvapi.RVer",
	HandlerType: (*RVerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RVer_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _RVer_Get_Handler,
		},
		{
			MethodName: "Overview",
			Handler:    _RVer_Overview_Handler,
		},
		{
			MethodName: "Report",
			Handler:    _RVer_Report_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _RVer_Update_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _RVer_Vote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/rvapi/rvapi.proto",
}
