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
	CheckIn(ctx context.Context, in *CheckInRequest, opts ...grpc.CallOption) (*CheckInResponse, error)
	TrustedCheckIn(ctx context.Context, in *TrustedCheckInRequest, opts ...grpc.CallOption) (*TrustedCheckInResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetView(ctx context.Context, in *GetViewRequest, opts ...grpc.CallOption) (*GetViewResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	ListViews(ctx context.Context, in *ListViewsRequest, opts ...grpc.CallOption) (*ListViewsResponse, error)
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

func (c *rVerClient) CheckIn(ctx context.Context, in *CheckInRequest, opts ...grpc.CallOption) (*CheckInResponse, error) {
	out := new(CheckInResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/CheckIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) TrustedCheckIn(ctx context.Context, in *TrustedCheckInRequest, opts ...grpc.CallOption) (*TrustedCheckInResponse, error) {
	out := new(TrustedCheckInResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/TrustedCheckIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *rVerClient) GetView(ctx context.Context, in *GetViewRequest, opts ...grpc.CallOption) (*GetViewResponse, error) {
	out := new(GetViewResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/GetView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rVerClient) ListViews(ctx context.Context, in *ListViewsRequest, opts ...grpc.CallOption) (*ListViewsResponse, error) {
	out := new(ListViewsResponse)
	err := c.cc.Invoke(ctx, "/rvapi.RVer/ListViews", in, out, opts...)
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
	CheckIn(context.Context, *CheckInRequest) (*CheckInResponse, error)
	TrustedCheckIn(context.Context, *TrustedCheckInRequest) (*TrustedCheckInResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	GetView(context.Context, *GetViewRequest) (*GetViewResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	ListViews(context.Context, *ListViewsRequest) (*ListViewsResponse, error)
	Report(context.Context, *ReportRequest) (*ReportResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
	mustEmbedUnimplementedRVerServer()
}

// UnimplementedRVerServer must be embedded to have forward compatible implementations.
type UnimplementedRVerServer struct {
}

func (*UnimplementedRVerServer) CheckIn(context.Context, *CheckInRequest) (*CheckInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIn not implemented")
}
func (*UnimplementedRVerServer) TrustedCheckIn(context.Context, *TrustedCheckInRequest) (*TrustedCheckInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrustedCheckIn not implemented")
}
func (*UnimplementedRVerServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedRVerServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedRVerServer) GetView(context.Context, *GetViewRequest) (*GetViewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetView not implemented")
}
func (*UnimplementedRVerServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*UnimplementedRVerServer) ListViews(context.Context, *ListViewsRequest) (*ListViewsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListViews not implemented")
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

func _RVer_CheckIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).CheckIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/CheckIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).CheckIn(ctx, req.(*CheckInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_TrustedCheckIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrustedCheckInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).TrustedCheckIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/TrustedCheckIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).TrustedCheckIn(ctx, req.(*TrustedCheckInRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _RVer_GetView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).GetView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/GetView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).GetView(ctx, req.(*GetViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RVer_ListViews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListViewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RVerServer).ListViews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rvapi.RVer/ListViews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RVerServer).ListViews(ctx, req.(*ListViewsRequest))
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
			MethodName: "CheckIn",
			Handler:    _RVer_CheckIn_Handler,
		},
		{
			MethodName: "TrustedCheckIn",
			Handler:    _RVer_TrustedCheckIn_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _RVer_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _RVer_Get_Handler,
		},
		{
			MethodName: "GetView",
			Handler:    _RVer_GetView_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RVer_List_Handler,
		},
		{
			MethodName: "ListViews",
			Handler:    _RVer_ListViews_Handler,
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
