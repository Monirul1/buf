// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.19.1
// source: buf/alpha/registry/v1alpha1/generate.proto

package registryv1alpha1

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

// GenerateServiceClient is the client API for GenerateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GenerateServiceClient interface {
	// GeneratePlugins generates an array of files given the provided
	// module reference and plugin version and option tuples. No attempt
	// is made at merging insertion points.
	GeneratePlugins(ctx context.Context, in *GeneratePluginsRequest, opts ...grpc.CallOption) (*GeneratePluginsResponse, error)
	// GenerateTemplate generates an array of files given the provided
	// module reference and template version.
	GenerateTemplate(ctx context.Context, in *GenerateTemplateRequest, opts ...grpc.CallOption) (*GenerateTemplateResponse, error)
}

type generateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGenerateServiceClient(cc grpc.ClientConnInterface) GenerateServiceClient {
	return &generateServiceClient{cc}
}

func (c *generateServiceClient) GeneratePlugins(ctx context.Context, in *GeneratePluginsRequest, opts ...grpc.CallOption) (*GeneratePluginsResponse, error) {
	out := new(GeneratePluginsResponse)
	err := c.cc.Invoke(ctx, "/buf.alpha.registry.v1alpha1.GenerateService/GeneratePlugins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *generateServiceClient) GenerateTemplate(ctx context.Context, in *GenerateTemplateRequest, opts ...grpc.CallOption) (*GenerateTemplateResponse, error) {
	out := new(GenerateTemplateResponse)
	err := c.cc.Invoke(ctx, "/buf.alpha.registry.v1alpha1.GenerateService/GenerateTemplate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GenerateServiceServer is the server API for GenerateService service.
// All implementations should embed UnimplementedGenerateServiceServer
// for forward compatibility
type GenerateServiceServer interface {
	// GeneratePlugins generates an array of files given the provided
	// module reference and plugin version and option tuples. No attempt
	// is made at merging insertion points.
	GeneratePlugins(context.Context, *GeneratePluginsRequest) (*GeneratePluginsResponse, error)
	// GenerateTemplate generates an array of files given the provided
	// module reference and template version.
	GenerateTemplate(context.Context, *GenerateTemplateRequest) (*GenerateTemplateResponse, error)
}

// UnimplementedGenerateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedGenerateServiceServer struct {
}

func (UnimplementedGenerateServiceServer) GeneratePlugins(context.Context, *GeneratePluginsRequest) (*GeneratePluginsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratePlugins not implemented")
}
func (UnimplementedGenerateServiceServer) GenerateTemplate(context.Context, *GenerateTemplateRequest) (*GenerateTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateTemplate not implemented")
}

// UnsafeGenerateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GenerateServiceServer will
// result in compilation errors.
type UnsafeGenerateServiceServer interface {
	mustEmbedUnimplementedGenerateServiceServer()
}

func RegisterGenerateServiceServer(s grpc.ServiceRegistrar, srv GenerateServiceServer) {
	s.RegisterService(&GenerateService_ServiceDesc, srv)
}

func _GenerateService_GeneratePlugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneratePluginsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenerateServiceServer).GeneratePlugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buf.alpha.registry.v1alpha1.GenerateService/GeneratePlugins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenerateServiceServer).GeneratePlugins(ctx, req.(*GeneratePluginsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenerateService_GenerateTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenerateServiceServer).GenerateTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buf.alpha.registry.v1alpha1.GenerateService/GenerateTemplate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenerateServiceServer).GenerateTemplate(ctx, req.(*GenerateTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GenerateService_ServiceDesc is the grpc.ServiceDesc for GenerateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GenerateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "buf.alpha.registry.v1alpha1.GenerateService",
	HandlerType: (*GenerateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GeneratePlugins",
			Handler:    _GenerateService_GeneratePlugins_Handler,
		},
		{
			MethodName: "GenerateTemplate",
			Handler:    _GenerateService_GenerateTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "buf/alpha/registry/v1alpha1/generate.proto",
}
