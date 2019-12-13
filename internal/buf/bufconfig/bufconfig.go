// Package bufconfig contains the configuration functionality.
package bufconfig

import (
	"bytes"
	"context"
	"io"
	"sort"

	"github.com/bufbuild/buf/internal/buf/bufcheck/bufbreaking"
	"github.com/bufbuild/buf/internal/buf/bufcheck/buflint"
	"github.com/bufbuild/buf/internal/pkg/analysis"
	"github.com/bufbuild/buf/internal/pkg/storage"
	"go.uber.org/zap"
)

// ConfigFilePath is the default config file path within a bucket.
//
// TODO: make sure copied for git
const ConfigFilePath = "buf.yaml"

// Config is the user config.
//
// Configs must not be linked to a specific Bucket object, that is if a Config
// is Generated from bucket1, and bucket1 is copied to bucket2, the Config must
// be valid for bucket2.
//
// TODO: remove individual configs as part of refactor.
type Config struct {
	Build    ExternalBuildConfig
	Breaking *bufbreaking.Config
	Lint     *buflint.Config
}

// Provider is a provider.
type Provider interface {
	// GetConfigForBucket gets the Config for the given JSON or YAML data.
	//
	// If the data is of length 0, returns the default config.
	GetConfigForBucket(ctx context.Context, bucket storage.ReadBucket) (*Config, error)
	// GetConfig gets the Config for the given JSON or YAML data.
	//
	// If the data is of length 0, returns the default config.
	GetConfigForData(data []byte) (*Config, error)
}

// ProviderOption is an option for a new Provider.
type ProviderOption func(*provider)

// ProviderWithExternalConfigModifier returns a new ProviderOption that applies the following
// external config modifier before processing an ExternalConfig.
//
// Useful for testing.
func ProviderWithExternalConfigModifier(externalConfigModifier func(*ExternalConfig) error) ProviderOption {
	return func(provider *provider) {
		provider.externalConfigModifier = externalConfigModifier
	}
}

// NewProvider returns a new Provider.
func NewProvider(logger *zap.Logger, options ...ProviderOption) Provider {
	return newProvider(logger, options...)
}

// ExternalBuildConfig is an external config.
type ExternalBuildConfig struct {
	Roots    []string `json:"roots,omitempty" yaml:"roots,omitempty"`
	Excludes []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`
}

// ExternalConfig is an external config.
//
// Should only be used outside this package for testing.
type ExternalConfig struct {
	Build    ExternalBuildConfig    `json:"build,omitempty" yaml:"build,omitempty"`
	Breaking ExternalBreakingConfig `json:"breaking,omitempty" yaml:"breaking,omitempty"`
	Lint     ExternalLintConfig     `json:"lint,omitempty" yaml:"lint,omitempty"`
}

// ExternalBreakingConfig is an external config.
//
// Should only be used outside this package for testing.
type ExternalBreakingConfig struct {
	Use        []string            `json:"use,omitempty" yaml:"use,omitempty"`
	Except     []string            `json:"except,omitempty" yaml:"except,omitempty"`
	Ignore     []string            `json:"ignore,omitempty" yaml:"ignore,omitempty"`
	IgnoreOnly map[string][]string `json:"ignore_only,omitempty" yaml:"ignore_only,omitempty"`
}

// ExternalLintConfig is an external config.
//
// Should only be used outside this package for testing.
type ExternalLintConfig struct {
	Use                                  []string            `json:"use,omitempty" yaml:"use,omitempty"`
	Except                               []string            `json:"except,omitempty" yaml:"except,omitempty"`
	Ignore                               []string            `json:"ignore,omitempty" yaml:"ignore,omitempty"`
	IgnoreOnly                           map[string][]string `json:"ignore_only,omitempty" yaml:"ignore_only,omitempty"`
	EnumZeroValueSuffix                  string              `json:"enum_zero_value_suffix,omitempty" yaml:"enum_zero_value_suffix,omitempty"`
	RPCAllowSameRequestResponse          bool                `json:"rpc_allow_same_request_response,omitempty" yaml:"rpc_allow_same_request_response,omitempty"`
	RPCAllowGoogleProtobufEmptyRequests  bool                `json:"rpc_allow_google_protobuf_empty_requests,omitempty" yaml:"rpc_allow_google_protobuf_empty_requests,omitempty"`
	RPCAllowGoogleProtobufEmptyResponses bool                `json:"rpc_allow_google_protobuf_empty_responses,omitempty" yaml:"rpc_allow_google_protobuf_empty_responses,omitempty"`
	ServiceSuffix                        string              `json:"service_suffix,omitempty" yaml:"service_suffix,omitempty"`
}

// PrintAnnotationsLintConfigIgnoreYAML prints the annotations to the Writer as config-ignore-yaml.
//
// TODO: this probably belongs in buflint, but since ExternalConfig is not supposed to be used
// outside of this package, we put it here for now.
func PrintAnnotationsLintConfigIgnoreYAML(writer io.Writer, annotations []*analysis.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}
	ignoreIDToRootPathMap := make(map[string]map[string]struct{})
	for _, annotation := range annotations {
		if annotation.Filename == "" || annotation.Type == "" {
			continue
		}
		rootPathMap, ok := ignoreIDToRootPathMap[annotation.Type]
		if !ok {
			rootPathMap = make(map[string]struct{})
			ignoreIDToRootPathMap[annotation.Type] = rootPathMap
		}
		rootPathMap[annotation.Filename] = struct{}{}
	}
	if len(ignoreIDToRootPathMap) == 0 {
		return nil
	}

	sortedIgnoreIDs := make([]string, 0, len(ignoreIDToRootPathMap))
	ignoreIDToSortedRootPaths := make(map[string][]string, len(ignoreIDToRootPathMap))
	for id, rootPathMap := range ignoreIDToRootPathMap {
		sortedIgnoreIDs = append(sortedIgnoreIDs, id)
		rootPaths := make([]string, 0, len(rootPathMap))
		for rootPath := range rootPathMap {
			rootPaths = append(rootPaths, rootPath)
		}
		sort.Strings(rootPaths)
		ignoreIDToSortedRootPaths[id] = rootPaths
	}
	sort.Strings(sortedIgnoreIDs)

	buffer := bytes.NewBuffer(nil)
	_, _ = buffer.WriteString(`lint:
  ignore_only:
`)
	for _, id := range sortedIgnoreIDs {
		_, _ = buffer.WriteString("    ")
		_, _ = buffer.WriteString(id)
		_, _ = buffer.WriteString(":\n")
		for _, rootPath := range ignoreIDToSortedRootPaths[id] {
			_, _ = buffer.WriteString("      - ")
			_, _ = buffer.WriteString(rootPath)
			_, _ = buffer.WriteString("\n")
		}
	}
	_, err := writer.Write(buffer.Bytes())
	return err
}
