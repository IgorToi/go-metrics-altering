// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: go-metrics-altering.proto

package proto

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
	Metrics_AddGaugeMetric_FullMethodName   = "/metrics.Metrics/AddGaugeMetric"
	Metrics_AddCounterMetric_FullMethodName = "/metrics.Metrics/AddCounterMetric"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	AddGaugeMetric(ctx context.Context, in *AddGaugeRequest, opts ...grpc.CallOption) (*AddGaugeResponse, error)
	AddCounterMetric(ctx context.Context, in *AddCounterRequest, opts ...grpc.CallOption) (*AddCounterResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) AddGaugeMetric(ctx context.Context, in *AddGaugeRequest, opts ...grpc.CallOption) (*AddGaugeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddGaugeResponse)
	err := c.cc.Invoke(ctx, Metrics_AddGaugeMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) AddCounterMetric(ctx context.Context, in *AddCounterRequest, opts ...grpc.CallOption) (*AddCounterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddCounterResponse)
	err := c.cc.Invoke(ctx, Metrics_AddCounterMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility.
type MetricsServer interface {
	AddGaugeMetric(context.Context, *AddGaugeRequest) (*AddGaugeResponse, error)
	AddCounterMetric(context.Context, *AddCounterRequest) (*AddCounterResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServer struct{}

func (UnimplementedMetricsServer) AddGaugeMetric(context.Context, *AddGaugeRequest) (*AddGaugeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGaugeMetric not implemented")
}
func (UnimplementedMetricsServer) AddCounterMetric(context.Context, *AddCounterRequest) (*AddCounterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCounterMetric not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}
func (UnimplementedMetricsServer) testEmbeddedByValue()                 {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_AddGaugeMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddGaugeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddGaugeMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddGaugeMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddGaugeMetric(ctx, req.(*AddGaugeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_AddCounterMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCounterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).AddCounterMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_AddCounterMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).AddCounterMetric(ctx, req.(*AddCounterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddGaugeMetric",
			Handler:    _Metrics_AddGaugeMetric_Handler,
		},
		{
			MethodName: "AddCounterMetric",
			Handler:    _Metrics_AddCounterMetric_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "go-metrics-altering.proto",
}
