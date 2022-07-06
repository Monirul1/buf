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

package apiclientconnect

import (
	apiclient "github.com/bufbuild/buf/private/gen/proto/apiclient"
	registryv1alpha1apiclient "github.com/bufbuild/buf/private/gen/proto/apiclient/buf/alpha/registry/v1alpha1/registryv1alpha1apiclient"
	webhookv1alpha1apiclient "github.com/bufbuild/buf/private/gen/proto/apiclient/buf/alpha/webhook/v1alpha1/webhookv1alpha1apiclient"
	registryv1alpha1apiclientconnect "github.com/bufbuild/buf/private/gen/proto/apiclientconnect/buf/alpha/registry/v1alpha1/registryv1alpha1apiclientconnect"
	webhookv1alpha1apiclientconnect "github.com/bufbuild/buf/private/gen/proto/apiclientconnect/buf/alpha/webhook/v1alpha1/webhookv1alpha1apiclientconnect"
	connect_go "github.com/bufbuild/connect-go"
	zap "go.uber.org/zap"
)

// NewProvider returns a new provider.
func NewProvider(
	logger *zap.Logger,
	httpClient connect_go.HTTPClient,
	options ...ProviderOption,
) apiclient.Provider {
	providerOptions := &providerOptions{}
	for _, option := range options {
		option(providerOptions)
	}
	return &provider{
		bufAlphaRegistryV1alpha1Provider: registryv1alpha1apiclientconnect.NewProvider(
			logger,
			httpClient,
			providerOptions.bufAlphaRegistryV1alpha1ProviderOptions...,
		),
		bufAlphaWebhookV1alpha1Provider: webhookv1alpha1apiclientconnect.NewProvider(
			logger,
			httpClient,
			providerOptions.bufAlphaWebhookV1alpha1ProviderOptions...,
		),
	}
}

type provider struct {
	bufAlphaRegistryV1alpha1Provider registryv1alpha1apiclient.Provider
	bufAlphaWebhookV1alpha1Provider  webhookv1alpha1apiclient.Provider
}

// ProviderOption is an option for a new Provider.
type ProviderOption func(*providerOptions)

func (p *provider) BufAlphaRegistryV1alpha1() registryv1alpha1apiclient.Provider {
	return p.bufAlphaRegistryV1alpha1Provider
}

func (p *provider) BufAlphaWebhookV1alpha1() webhookv1alpha1apiclient.Provider {
	return p.bufAlphaWebhookV1alpha1Provider
}

type providerOptions struct {
	bufAlphaRegistryV1alpha1ProviderOptions []registryv1alpha1apiclientconnect.ProviderOption
	bufAlphaWebhookV1alpha1ProviderOptions  []webhookv1alpha1apiclientconnect.ProviderOption
}
