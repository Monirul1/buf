// Copyright 2020-2022 Buf Technologies, Inc.
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

// Code generated by protoc-gen-go-apiclientconnect. DO NOT EDIT.

package registryv1alpha1apiclientconnect

import (
	context "context"
	registryv1alpha1connect "github.com/bufbuild/buf/private/gen/proto/connect/buf/alpha/registry/v1alpha1/registryv1alpha1connect"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	connect_go "github.com/bufbuild/connect-go"
	zap "go.uber.org/zap"
)

type studioRequestServiceClient struct {
	logger *zap.Logger
	client registryv1alpha1connect.StudioRequestServiceClient
}

// CreateStudioRequest registers a favorite Studio Requests to the caller's
// BSR profile.
func (s *studioRequestServiceClient) CreateStudioRequest(
	ctx context.Context,
	repositoryOwner string,
	repositoryName string,
	name string,
	targetBaseUrl string,
	service string,
	method string,
	body string,
	headers map[string]string,
	includeCookies bool,
	protocol v1alpha1.StudioProtocol,
	agentUrl string,
) (createdRequest *v1alpha1.StudioRequest, _ error) {
	response, err := s.client.CreateStudioRequest(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.CreateStudioRequestRequest{
				RepositoryOwner: repositoryOwner,
				RepositoryName:  repositoryName,
				Name:            name,
				TargetBaseUrl:   targetBaseUrl,
				Service:         service,
				Method:          method,
				Body:            body,
				Headers:         headers,
				IncludeCookies:  includeCookies,
				Protocol:        protocol,
				AgentUrl:        agentUrl,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.CreatedRequest, nil
}

// RenameStudioRequest renames an existing Studio Request.
func (s *studioRequestServiceClient) RenameStudioRequest(
	ctx context.Context,
	id string,
	newName string,
) (renamedRequest *v1alpha1.StudioRequest, _ error) {
	response, err := s.client.RenameStudioRequest(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.RenameStudioRequestRequest{
				Id:      id,
				NewName: newName,
			}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg.RenamedRequest, nil
}

// DeleteStudioRequest removes a favorite Studio Request from the caller's BSR
// profile.
func (s *studioRequestServiceClient) DeleteStudioRequest(ctx context.Context, id string) (_ error) {
	_, err := s.client.DeleteStudioRequest(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.DeleteStudioRequestRequest{
				Id: id,
			}),
	)
	if err != nil {
		return err
	}
	return nil
}

// ListStudioRequests shows the caller's favorited Studio Requests.
func (s *studioRequestServiceClient) ListStudioRequests(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
	reverse bool,
) (requests []*v1alpha1.StudioRequest, nextPageToken string, _ error) {
	response, err := s.client.ListStudioRequests(
		ctx,
		connect_go.NewRequest(
			&v1alpha1.ListStudioRequestsRequest{
				PageSize:  pageSize,
				PageToken: pageToken,
				Reverse:   reverse,
			}),
	)
	if err != nil {
		return nil, "", err
	}
	return response.Msg.Requests, response.Msg.NextPageToken, nil
}
