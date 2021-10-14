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

// Code generated by protoc-gen-go-apiclientgrpc. DO NOT EDIT.

package registryv1alpha1apiclientgrpc

import (
	context "context"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
)

type recommendationService struct {
	logger          *zap.Logger
	client          v1alpha1.RecommendationServiceClient
	contextModifier func(context.Context) context.Context
}

// RecommendedRepositories returns a list of recommended repositories.
func (s *recommendationService) RecommendedRepositories(ctx context.Context) (repositories []*v1alpha1.RecommendedRepository, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.RecommendedRepositories(
		ctx,
		&v1alpha1.RecommendedRepositoriesRequest{},
	)
	if err != nil {
		return nil, err
	}
	return response.Repositories, nil
}

// RecommendedTemplates returns a list of recommended templates.
func (s *recommendationService) RecommendedTemplates(ctx context.Context) (templates []*v1alpha1.RecommendedTemplate, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.RecommendedTemplates(
		ctx,
		&v1alpha1.RecommendedTemplatesRequest{},
	)
	if err != nil {
		return nil, err
	}
	return response.Templates, nil
}
