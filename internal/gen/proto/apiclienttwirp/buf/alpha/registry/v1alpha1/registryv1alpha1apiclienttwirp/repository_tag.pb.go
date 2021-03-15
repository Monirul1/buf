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

// Code generated by protoc-gen-go-apiclienttwirp. DO NOT EDIT.

package registryv1alpha1apiclienttwirp

import (
	context "context"
	v1alpha1 "github.com/bufbuild/buf/internal/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
)

type repositoryTagService struct {
	logger          *zap.Logger
	client          v1alpha1.RepositoryTagService
	contextModifier func(context.Context) context.Context
}

// CreateRepositoryTag creates a new repository tag.
func (s *repositoryTagService) CreateRepositoryTag(
	ctx context.Context,
	repositoryId string,
	name string,
	commitName string,
) (repositoryTag *v1alpha1.RepositoryTag, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.CreateRepositoryTag(
		ctx,
		&v1alpha1.CreateRepositoryTagRequest{
			RepositoryId: repositoryId,
			Name:         name,
			CommitName:   commitName,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.RepositoryTag, nil
}

// ListRepositoryTags lists the repository tags associated with a Repository.
func (s *repositoryTagService) ListRepositoryTags(
	ctx context.Context,
	repositoryId string,
	pageSize uint32,
	pageToken string,
	reverse bool,
) (repositoryTags []*v1alpha1.RepositoryTag, nextPageToken string, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.ListRepositoryTags(
		ctx,
		&v1alpha1.ListRepositoryTagsRequest{
			RepositoryId: repositoryId,
			PageSize:     pageSize,
			PageToken:    pageToken,
			Reverse:      reverse,
		},
	)
	if err != nil {
		return nil, "", err
	}
	return response.RepositoryTags, response.NextPageToken, nil
}
