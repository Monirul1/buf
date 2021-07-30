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

// Code generated by protoc-gen-twirp v8.1.0, DO NOT EDIT.
// source: buf/alpha/registry/v1alpha1/owner.proto

package registryv1alpha1

import context "context"
import fmt "fmt"
import http "net/http"
import ioutil "io/ioutil"
import json "encoding/json"
import strconv "strconv"
import strings "strings"

import protojson "google.golang.org/protobuf/encoding/protojson"
import proto "google.golang.org/protobuf/proto"
import twirp "github.com/twitchtv/twirp"
import ctxsetters "github.com/twitchtv/twirp/ctxsetters"

// Version compatibility assertion.
// If the constant is not defined in the package, that likely means
// the package needs to be updated to work with this generated code.
// See https://twitchtv.github.io/twirp/docs/version_matrix.html
const _ = twirp.TwirpPackageMinVersion_8_1_0

// ======================
// OwnerService Interface
// ======================

// OwnerService is a service that provides RPCs that allow the BSR to query
// for owner information.
type OwnerService interface {
	// GetOwnerByName takes an owner name and returns the owner as
	// either a user or organization.
	GetOwnerByName(context.Context, *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error)
}

// ============================
// OwnerService Protobuf Client
// ============================

type ownerServiceProtobufClient struct {
	client      HTTPClient
	urls        [1]string
	interceptor twirp.Interceptor
	opts        twirp.ClientOptions
}

// NewOwnerServiceProtobufClient creates a Protobuf client that implements the OwnerService interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewOwnerServiceProtobufClient(baseURL string, client HTTPClient, opts ...twirp.ClientOption) OwnerService {
	if c, ok := client.(*http.Client); ok {
		client = withoutRedirects(c)
	}

	clientOpts := twirp.ClientOptions{}
	for _, o := range opts {
		o(&clientOpts)
	}

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	literalURLs := false
	_ = clientOpts.ReadOpt("literalURLs", &literalURLs)
	var pathPrefix string
	if ok := clientOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	// Build method URLs: <baseURL>[<prefix>]/<package>.<Service>/<Method>
	serviceURL := sanitizeBaseURL(baseURL)
	serviceURL += baseServicePath(pathPrefix, "buf.alpha.registry.v1alpha1", "OwnerService")
	urls := [1]string{
		serviceURL + "GetOwnerByName",
	}

	return &ownerServiceProtobufClient{
		client:      client,
		urls:        urls,
		interceptor: twirp.ChainInterceptors(clientOpts.Interceptors...),
		opts:        clientOpts,
	}
}

func (c *ownerServiceProtobufClient) GetOwnerByName(ctx context.Context, in *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "buf.alpha.registry.v1alpha1")
	ctx = ctxsetters.WithServiceName(ctx, "OwnerService")
	ctx = ctxsetters.WithMethodName(ctx, "GetOwnerByName")
	caller := c.callGetOwnerByName
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*GetOwnerByNameRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*GetOwnerByNameRequest) when calling interceptor")
					}
					return c.callGetOwnerByName(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*GetOwnerByNameResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*GetOwnerByNameResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *ownerServiceProtobufClient) callGetOwnerByName(ctx context.Context, in *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
	out := new(GetOwnerByNameResponse)
	ctx, err := doProtobufRequest(ctx, c.client, c.opts.Hooks, c.urls[0], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

// ========================
// OwnerService JSON Client
// ========================

type ownerServiceJSONClient struct {
	client      HTTPClient
	urls        [1]string
	interceptor twirp.Interceptor
	opts        twirp.ClientOptions
}

// NewOwnerServiceJSONClient creates a JSON client that implements the OwnerService interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewOwnerServiceJSONClient(baseURL string, client HTTPClient, opts ...twirp.ClientOption) OwnerService {
	if c, ok := client.(*http.Client); ok {
		client = withoutRedirects(c)
	}

	clientOpts := twirp.ClientOptions{}
	for _, o := range opts {
		o(&clientOpts)
	}

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	literalURLs := false
	_ = clientOpts.ReadOpt("literalURLs", &literalURLs)
	var pathPrefix string
	if ok := clientOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	// Build method URLs: <baseURL>[<prefix>]/<package>.<Service>/<Method>
	serviceURL := sanitizeBaseURL(baseURL)
	serviceURL += baseServicePath(pathPrefix, "buf.alpha.registry.v1alpha1", "OwnerService")
	urls := [1]string{
		serviceURL + "GetOwnerByName",
	}

	return &ownerServiceJSONClient{
		client:      client,
		urls:        urls,
		interceptor: twirp.ChainInterceptors(clientOpts.Interceptors...),
		opts:        clientOpts,
	}
}

func (c *ownerServiceJSONClient) GetOwnerByName(ctx context.Context, in *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
	ctx = ctxsetters.WithPackageName(ctx, "buf.alpha.registry.v1alpha1")
	ctx = ctxsetters.WithServiceName(ctx, "OwnerService")
	ctx = ctxsetters.WithMethodName(ctx, "GetOwnerByName")
	caller := c.callGetOwnerByName
	if c.interceptor != nil {
		caller = func(ctx context.Context, req *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
			resp, err := c.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*GetOwnerByNameRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*GetOwnerByNameRequest) when calling interceptor")
					}
					return c.callGetOwnerByName(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*GetOwnerByNameResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*GetOwnerByNameResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}
	return caller(ctx, in)
}

func (c *ownerServiceJSONClient) callGetOwnerByName(ctx context.Context, in *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
	out := new(GetOwnerByNameResponse)
	ctx, err := doJSONRequest(ctx, c.client, c.opts.Hooks, c.urls[0], in, out)
	if err != nil {
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}
		callClientError(ctx, c.opts.Hooks, twerr)
		return nil, err
	}

	callClientResponseReceived(ctx, c.opts.Hooks)

	return out, nil
}

// ===========================
// OwnerService Server Handler
// ===========================

type ownerServiceServer struct {
	OwnerService
	interceptor      twirp.Interceptor
	hooks            *twirp.ServerHooks
	pathPrefix       string // prefix for routing
	jsonSkipDefaults bool   // do not include unpopulated fields (default values) in the response
	jsonCamelCase    bool   // JSON fields are serialized as lowerCamelCase rather than keeping the original proto names
}

// NewOwnerServiceServer builds a TwirpServer that can be used as an http.Handler to handle
// HTTP requests that are routed to the right method in the provided svc implementation.
// The opts are twirp.ServerOption modifiers, for example twirp.WithServerHooks(hooks).
func NewOwnerServiceServer(svc OwnerService, opts ...interface{}) TwirpServer {
	serverOpts := newServerOpts(opts)

	// Using ReadOpt allows backwards and forwads compatibility with new options in the future
	jsonSkipDefaults := false
	_ = serverOpts.ReadOpt("jsonSkipDefaults", &jsonSkipDefaults)
	jsonCamelCase := false
	_ = serverOpts.ReadOpt("jsonCamelCase", &jsonCamelCase)
	var pathPrefix string
	if ok := serverOpts.ReadOpt("pathPrefix", &pathPrefix); !ok {
		pathPrefix = "/twirp" // default prefix
	}

	return &ownerServiceServer{
		OwnerService:     svc,
		hooks:            serverOpts.Hooks,
		interceptor:      twirp.ChainInterceptors(serverOpts.Interceptors...),
		pathPrefix:       pathPrefix,
		jsonSkipDefaults: jsonSkipDefaults,
		jsonCamelCase:    jsonCamelCase,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *ownerServiceServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// handleRequestBodyError is used to handle error when the twirp server cannot read request
func (s *ownerServiceServer) handleRequestBodyError(ctx context.Context, resp http.ResponseWriter, msg string, err error) {
	if context.Canceled == ctx.Err() {
		s.writeError(ctx, resp, twirp.NewError(twirp.Canceled, "failed to read request: context canceled"))
		return
	}
	if context.DeadlineExceeded == ctx.Err() {
		s.writeError(ctx, resp, twirp.NewError(twirp.DeadlineExceeded, "failed to read request: deadline exceeded"))
		return
	}
	s.writeError(ctx, resp, twirp.WrapError(malformedRequestError(msg), err))
}

// OwnerServicePathPrefix is a convenience constant that may identify URL paths.
// Should be used with caution, it only matches routes generated by Twirp Go clients,
// with the default "/twirp" prefix and default CamelCase service and method names.
// More info: https://twitchtv.github.io/twirp/docs/routing.html
const OwnerServicePathPrefix = "/twirp/buf.alpha.registry.v1alpha1.OwnerService/"

func (s *ownerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "buf.alpha.registry.v1alpha1")
	ctx = ctxsetters.WithServiceName(ctx, "OwnerService")
	ctx = ctxsetters.WithResponseWriter(ctx, resp)

	var err error
	ctx, err = callRequestReceived(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}

	// Verify path format: [<prefix>]/<package>.<Service>/<Method>
	prefix, pkgService, method := parseTwirpPath(req.URL.Path)
	if pkgService != "buf.alpha.registry.v1alpha1.OwnerService" {
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}
	if prefix != s.pathPrefix {
		msg := fmt.Sprintf("invalid path prefix %q, expected %q, on path %q", prefix, s.pathPrefix, req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}

	switch method {
	case "GetOwnerByName":
		s.serveGetOwnerByName(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		s.writeError(ctx, resp, badRouteError(msg, req.Method, req.URL.Path))
		return
	}
}

func (s *ownerServiceServer) serveGetOwnerByName(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetOwnerByNameJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveGetOwnerByNameProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *ownerServiceServer) serveGetOwnerByNameJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetOwnerByName")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	d := json.NewDecoder(req.Body)
	rawReqBody := json.RawMessage{}
	if err := d.Decode(&rawReqBody); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}
	reqContent := new(GetOwnerByNameRequest)
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = unmarshaler.Unmarshal(rawReqBody, reqContent); err != nil {
		s.handleRequestBodyError(ctx, resp, "the json request could not be decoded", err)
		return
	}

	handler := s.OwnerService.GetOwnerByName
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*GetOwnerByNameRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*GetOwnerByNameRequest) when calling interceptor")
					}
					return s.OwnerService.GetOwnerByName(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*GetOwnerByNameResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*GetOwnerByNameResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *GetOwnerByNameResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetOwnerByNameResponse and nil error while calling GetOwnerByName. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	marshaler := &protojson.MarshalOptions{UseProtoNames: !s.jsonCamelCase, EmitUnpopulated: !s.jsonSkipDefaults}
	respBytes, err := marshaler.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal json response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)

	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *ownerServiceServer) serveGetOwnerByNameProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "GetOwnerByName")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.handleRequestBodyError(ctx, resp, "failed to read request body", err)
		return
	}
	reqContent := new(GetOwnerByNameRequest)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		s.writeError(ctx, resp, malformedRequestError("the protobuf request could not be decoded"))
		return
	}

	handler := s.OwnerService.GetOwnerByName
	if s.interceptor != nil {
		handler = func(ctx context.Context, req *GetOwnerByNameRequest) (*GetOwnerByNameResponse, error) {
			resp, err := s.interceptor(
				func(ctx context.Context, req interface{}) (interface{}, error) {
					typedReq, ok := req.(*GetOwnerByNameRequest)
					if !ok {
						return nil, twirp.InternalError("failed type assertion req.(*GetOwnerByNameRequest) when calling interceptor")
					}
					return s.OwnerService.GetOwnerByName(ctx, typedReq)
				},
			)(ctx, req)
			if resp != nil {
				typedResp, ok := resp.(*GetOwnerByNameResponse)
				if !ok {
					return nil, twirp.InternalError("failed type assertion resp.(*GetOwnerByNameResponse) when calling interceptor")
				}
				return typedResp, err
			}
			return nil, err
		}
	}

	// Call service method
	var respContent *GetOwnerByNameResponse
	func() {
		defer ensurePanicResponses(ctx, resp, s.hooks)
		respContent, err = handler(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *GetOwnerByNameResponse and nil error while calling GetOwnerByName. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		s.writeError(ctx, resp, wrapInternal(err, "failed to marshal proto response"))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		ctx = callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *ownerServiceServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor3, 0
}

func (s *ownerServiceServer) ProtocGenTwirpVersion() string {
	return "v8.1.0"
}

// PathPrefix returns the base service path, in the form: "/<prefix>/<package>.<Service>/"
// that is everything in a Twirp route except for the <Method>. This can be used for routing,
// for example to identify the requests that are targeted to this service in a mux.
func (s *ownerServiceServer) PathPrefix() string {
	return baseServicePath(s.pathPrefix, "buf.alpha.registry.v1alpha1", "OwnerService")
}

var twirpFileDescriptor3 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4d, 0x6a, 0xe3, 0x30,
	0x18, 0x86, 0x63, 0x4f, 0x32, 0xc3, 0x68, 0xc2, 0x2c, 0x04, 0x33, 0x84, 0x04, 0xe6, 0xc7, 0x8b,
	0xfe, 0x50, 0x90, 0x48, 0xb2, 0x68, 0x71, 0x57, 0xf1, 0xa6, 0x59, 0x35, 0x41, 0xa5, 0x5d, 0x04,
	0x6f, 0xec, 0x54, 0x76, 0x0c, 0xb1, 0x94, 0xca, 0x56, 0x4a, 0x7a, 0x83, 0x5e, 0xa1, 0xcb, 0x2e,
	0x7b, 0x94, 0x1e, 0xa5, 0xbb, 0xde, 0xa0, 0xf8, 0x4b, 0x04, 0x49, 0x29, 0x82, 0xee, 0xc4, 0xa7,
	0xe7, 0x7d, 0xf5, 0xc8, 0x16, 0xda, 0x8f, 0x75, 0x42, 0xa3, 0xf9, 0x62, 0x16, 0x51, 0xc5, 0xd3,
	0xac, 0x28, 0xd5, 0x8a, 0x2e, 0xbb, 0x30, 0xe8, 0x52, 0x79, 0x2b, 0xb8, 0x22, 0x0b, 0x25, 0x4b,
	0x89, 0x3b, 0xb1, 0x4e, 0x08, 0xcc, 0x89, 0x01, 0x89, 0x01, 0xdb, 0x7b, 0xb6, 0x16, 0x5d, 0x98,
	0x92, 0x36, 0xb1, 0x9e, 0xa6, 0xd2, 0x48, 0x64, 0x77, 0x51, 0x99, 0x49, 0xb1, 0xe6, 0xbd, 0x07,
	0x07, 0x35, 0x46, 0x95, 0x04, 0x3e, 0x46, 0xf5, 0xaa, 0xa7, 0xe5, 0xfc, 0x73, 0x0e, 0x7e, 0xf4,
	0xfe, 0x13, 0x8b, 0x0d, 0xb9, 0x2c, 0xb8, 0x1a, 0xd6, 0x18, 0x04, 0xf0, 0x08, 0x35, 0xb7, 0x8b,
	0x5b, 0x2e, 0x14, 0x1c, 0x5a, 0x0b, 0x46, 0x5b, 0x81, 0x61, 0x8d, 0xed, 0x14, 0x04, 0xdf, 0x50,
	0x03, 0xbe, 0x8b, 0x77, 0x84, 0x7e, 0x9d, 0xf1, 0x12, 0xf4, 0x82, 0xd5, 0x79, 0x94, 0x73, 0xc6,
	0x6f, 0x34, 0x2f, 0x4a, 0x8c, 0x51, 0x5d, 0x44, 0x39, 0x07, 0xd7, 0xef, 0x0c, 0xd6, 0x1e, 0x43,
	0xbf, 0xdf, 0xc3, 0xc5, 0x42, 0x8a, 0x82, 0xe3, 0x93, 0x4d, 0xdf, 0xe6, 0x6a, 0x9e, 0xdd, 0xac,
	0x22, 0xd9, 0x3a, 0xd0, 0xbb, 0x77, 0x50, 0x13, 0x06, 0x17, 0x5c, 0x2d, 0xb3, 0x29, 0xc7, 0x2b,
	0xf4, 0x73, 0xf7, 0x10, 0xdc, 0xb3, 0xb6, 0x7d, 0xa8, 0xdf, 0xee, 0x7f, 0x2a, 0xb3, 0xbe, 0x45,
	0xf0, 0xea, 0xa0, 0xbf, 0x53, 0x99, 0xdb, 0xa2, 0x01, 0x82, 0xe0, 0xb8, 0xfa, 0xb3, 0x63, 0x67,
	0x32, 0x49, 0xb3, 0x72, 0xa6, 0x63, 0x32, 0x95, 0x39, 0x8d, 0x75, 0x12, 0xeb, 0x6c, 0x7e, 0x5d,
	0x2d, 0x68, 0x26, 0x4a, 0xae, 0x44, 0x34, 0xa7, 0x29, 0x17, 0x14, 0xde, 0x01, 0x4d, 0x25, 0xb5,
	0xbc, 0x9c, 0x53, 0x33, 0x31, 0x83, 0x47, 0xf7, 0x4b, 0x30, 0x60, 0x4f, 0x6e, 0x27, 0xd0, 0x09,
	0x19, 0x80, 0x0d, 0x33, 0x36, 0x57, 0x1b, 0xe6, 0xd9, 0xfd, 0x13, 0xe8, 0x24, 0x0c, 0x61, 0x3b,
	0x0c, 0xcd, 0x7e, 0x18, 0x1a, 0xe0, 0x05, 0x00, 0xdf, 0x07, 0xc0, 0xf7, 0x0d, 0xe0, 0xfb, 0x06,
	0x88, 0xbf, 0x82, 0x5c, 0xff, 0x2d, 0x00, 0x00, 0xff, 0xff, 0xee, 0x8b, 0x7b, 0x18, 0x44, 0x03,
	0x00, 0x00,
}
