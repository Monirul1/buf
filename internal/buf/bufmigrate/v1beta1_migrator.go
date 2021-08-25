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

package bufmigrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bufbuild/buf/internal/buf/bufcheck/bufbreaking"
	"github.com/bufbuild/buf/internal/buf/bufcheck/buflint"
	"github.com/bufbuild/buf/internal/buf/bufconfig"
	"github.com/bufbuild/buf/internal/buf/bufgen"
	"github.com/bufbuild/buf/internal/buf/bufwork"
	"github.com/bufbuild/buf/private/bufpkg/buflock"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmodulebuild"
	"github.com/bufbuild/buf/private/pkg/encoding"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"go.uber.org/multierr"
)

const (
	bufModHeaderWithName = `# Generated by %q. Edit as necessary, and
# remove this comment when you're finished.
#
# This module represents the %q root found in
# the previous configuration file for the
# %q module.
`
	bufModHeaderWithoutName = `# Generated by %q. Edit as necessary, and
# remove this comment when you're finished.
#
# This module represents the %q root found in
# the previous configuration.
`
	bufGenHeader = `# Generated by %q. Edit as necessary, and
# remove this comment when you're finished.
`
	bufWorkHeader = `# Generated by %q. Edit as necessary, and
# remove this comment when you're finished.
#
# This workspace file points to the roots found in your
# previous %q configuration.
`
)

type v1beta1Migrator struct {
	notifier    func(string) error
	commandName string
}

func newV1Beta1Migrator(commandName string, options ...V1Beta1MigrateOption) *v1beta1Migrator {
	migrator := v1beta1Migrator{
		commandName: commandName,
		notifier:    func(string) error { return nil },
	}
	for _, option := range options {
		option(&migrator)
	}
	return &migrator
}

func (m *v1beta1Migrator) Migrate(dirPath string) error {
	migratedConfig, err := m.maybeMigrateConfig(dirPath)
	if err != nil {
		return fmt.Errorf("failed to migrate config: %w", err)
	}
	migratedGenTemplate, err := m.maybeMigrateGenTemplate(dirPath)
	if err != nil {
		return fmt.Errorf("failed to migrate generation template: %w", err)
	}
	migratedLockFile, err := m.maybeMigrateLockFile(dirPath)
	if err != nil {
		return fmt.Errorf("failed to migrate lock file: %w", err)
	}
	if !migratedConfig && !migratedGenTemplate && !migratedLockFile {
		return nil
	}
	var migratedFiles []string
	if migratedConfig {
		migratedFiles = append(migratedFiles, bufconfig.ExternalConfigV1Beta1FilePath)
	}
	if migratedGenTemplate {
		migratedFiles = append(migratedFiles, bufgen.ExternalConfigFilePath)
	}
	if migratedLockFile {
		migratedFiles = append(migratedFiles, buflock.ExternalConfigFilePath)
	}
	if err := m.notifier(
		fmt.Sprintf("Successfully migrated your %s to v1.\n", stringutil.SliceToHumanString(migratedFiles)),
	); err != nil {
		return fmt.Errorf("failed to write success message: %w", err)
	}
	return nil
}

func (m *v1beta1Migrator) maybeMigrateConfig(dirPath string) (bool, error) {
	oldConfigPath := filepath.Join(dirPath, bufconfig.ExternalConfigV1Beta1FilePath)
	oldConfigBytes, err := os.ReadFile(oldConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// OK, no old config file
			return false, nil
		}
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	var versionedConfig bufconfig.ExternalConfigVersion
	if err := encoding.UnmarshalYAMLNonStrict(oldConfigBytes, &versionedConfig); err != nil {
		return false, fmt.Errorf(
			"failed to read %s version: %w",
			oldConfigPath,
			err,
		)
	}
	switch versionedConfig.Version {
	case bufconfig.V1Version:
		// OK, file was already v1
		return false, nil
	case bufconfig.V1Beta1Version, "":
		// Continue to migrate
	default:
		return false, fmt.Errorf("unknown config file version: %s", versionedConfig.Version)
	}
	var v1beta1Config bufconfig.ExternalConfigV1Beta1
	if err := encoding.UnmarshalYAMLStrict(oldConfigBytes, &v1beta1Config); err != nil {
		return false, fmt.Errorf(
			"failed to unmarshal %s as %s version v1beta1: %w",
			oldConfigPath,
			bufconfig.ExternalConfigV1Beta1FilePath,
			err,
		)
	}
	buildConfig, err := bufmodulebuild.NewConfigV1Beta1(v1beta1Config.Build, v1beta1Config.Deps...)
	if err != nil {
		return false, err
	}
	if excludes, ok := buildConfig.RootToExcludes["."]; len(buildConfig.RootToExcludes) == 1 && ok {
		// Only "." root present, just recreate file
		v1Config := bufconfig.ExternalConfigV1{
			Version: bufconfig.V1Version,
			Name:    v1beta1Config.Name,
			Deps:    v1beta1Config.Deps,
			Build: bufmodulebuild.ExternalConfigV1{
				Excludes: excludes,
			},
			Breaking: bufbreaking.ExternalConfigV1(v1beta1Config.Breaking),
			Lint:     buflint.ExternalConfigV1(v1beta1Config.Lint),
		}
		newConfigPath := filepath.Join(dirPath, bufconfig.ExternalConfigV1FilePath)
		if err := m.writeV1Config(newConfigPath, v1Config, ".", v1beta1Config.Name); err != nil {
			return false, err
		}
		// Delete the old file once we've created the new one,
		// unless it's the same file as before.
		if newConfigPath != oldConfigPath {
			if err := os.Remove(oldConfigPath); err != nil {
				return false, fmt.Errorf("failed to delete old config file: %w", err)
			}
		}
		return true, nil
	}
	// Check if we have a co-resident lock file, of any version
	oldLockFilePath := filepath.Join(dirPath, buflock.ExternalConfigFilePath)
	externalLockFileV1, hasLockFile, err := maybeReadLockFile(oldLockFilePath)
	if err != nil {
		return false, err
	}
	pathToProcessed := make(map[string]bool)
	for root, excludes := range buildConfig.RootToExcludes {
		// Convert universal settings
		var name string
		if v1beta1Config.Name != "" {
			name = v1beta1Config.Name + "-" + strings.ReplaceAll(root, "/", "-") // Note: roots are normalized, "/" is universal
		}
		v1Config := bufconfig.ExternalConfigV1{
			Version: bufconfig.V1Version,
			Name:    name,
			Deps:    v1beta1Config.Deps,
			Build: bufmodulebuild.ExternalConfigV1{
				Excludes: excludes,
			},
			Breaking: bufbreaking.ExternalConfigV1{
				Use:                    v1beta1Config.Breaking.Use,
				Except:                 v1beta1Config.Breaking.Except,
				IgnoreUnstablePackages: v1beta1Config.Breaking.IgnoreUnstablePackages,
			},
			Lint: buflint.ExternalConfigV1{
				Use:                                  v1beta1Config.Lint.Use,
				Except:                               v1beta1Config.Lint.Except,
				ServiceSuffix:                        v1beta1Config.Lint.ServiceSuffix,
				EnumZeroValueSuffix:                  v1beta1Config.Lint.EnumZeroValueSuffix,
				RPCAllowSameRequestResponse:          v1beta1Config.Lint.RPCAllowSameRequestResponse,
				RPCAllowGoogleProtobufEmptyRequests:  v1beta1Config.Lint.RPCAllowGoogleProtobufEmptyRequests,
				RPCAllowGoogleProtobufEmptyResponses: v1beta1Config.Lint.RPCAllowGoogleProtobufEmptyResponses,
				AllowCommentIgnores:                  v1beta1Config.Lint.AllowCommentIgnores,
			},
		}

		// Process Ignore's for those related to the root
		v1Config.Breaking.Ignore, err = convertIgnoreSlice(v1beta1Config.Breaking.Ignore, dirPath, root, pathToProcessed)
		if err != nil {
			return false, err
		}
		v1Config.Breaking.IgnoreOnly, err = convertIgnoreMap(v1beta1Config.Breaking.IgnoreOnly, dirPath, root, pathToProcessed)
		if err != nil {
			return false, err
		}
		v1Config.Lint.Ignore, err = convertIgnoreSlice(v1beta1Config.Lint.Ignore, dirPath, root, pathToProcessed)
		if err != nil {
			return false, err
		}
		v1Config.Lint.IgnoreOnly, err = convertIgnoreMap(v1beta1Config.Lint.IgnoreOnly, dirPath, root, pathToProcessed)
		if err != nil {
			return false, err
		}
		if err := m.writeV1Config(
			filepath.Join(dirPath, root, bufconfig.ExternalConfigV1FilePath),
			v1Config,
			root,
			v1beta1Config.Name,
		); err != nil {
			return false, err
		}
		if hasLockFile {
			if err := m.writeV1LockFile(
				filepath.Join(dirPath, root, buflock.ExternalConfigFilePath),
				externalLockFileV1,
			); err != nil {
				return false, err
			}
		}
	}
	for path, processed := range pathToProcessed {
		if !processed {
			if err := m.notifier(
				fmt.Sprintf(
					"The ignored file %q was not found in any roots and has been removed.\n",
					path,
				),
			); err != nil {
				return false, fmt.Errorf("failed to warn about ignored file: %w", err)
			}
		}
	}
	workConfig := bufwork.ExternalConfigV1{
		Version:     bufwork.V1Version,
		Directories: v1beta1Config.Build.Roots,
	}
	// Sort directories before marshalling for deterministic output
	sort.Strings(workConfig.Directories)
	workConfigBytes, err := encoding.MarshalYAML(&workConfig)
	if err != nil {
		return false, fmt.Errorf("failed to marshal workspace file: %w", err)
	}
	header := fmt.Sprintf(bufWorkHeader, m.commandName, bufconfig.ExternalConfigV1Beta1FilePath)
	if err := os.WriteFile(
		filepath.Join(dirPath, bufwork.ExternalConfigV1FilePath),
		append([]byte(header), workConfigBytes...),
		0600,
	); err != nil {
		return false, fmt.Errorf("failed to write workspace file: %w", err)
	}
	// Finally, delete the old `buf.yaml` and any `buf.lock`. This is safe to do unconditionally
	// as we know that there can't be a new `buf.yaml` here, since the only case
	// where that would be true is if the only root is ".", which is handled separately.
	if err := os.Remove(oldConfigPath); err != nil {
		return false, fmt.Errorf("failed to clean up old config file: %w", err)
	}
	if hasLockFile {
		if err := os.Remove(oldLockFilePath); err != nil {
			return false, fmt.Errorf("failed to clean up old lock file: %w", err)
		}
	}
	return true, nil
}

// writeV1Config atomically replaces the old configuration file by first writing
// the new config to a temporary file and then moving it to the old config file path.
// If we fail to marshal or write, the old config file is not touched.
func (m *v1beta1Migrator) writeV1Config(
	configPath string,
	config bufconfig.ExternalConfigV1,
	originalRootName string,
	originalModuleName string,
) (retErr error) {
	v1ConfigBytes, err := encoding.MarshalYAML(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal new config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		// This happens if the user has a root specified that doesn't have a corresponding
		// directory on the filesystem.
		return fmt.Errorf("failed to create new directories for writing config: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for writing new config: %w", err)
	}
	header := fmt.Sprintf(bufModHeaderWithName, m.commandName, originalRootName, originalModuleName)
	if originalModuleName == "" {
		header = fmt.Sprintf(bufModHeaderWithoutName, m.commandName, originalRootName)
	}
	v1ConfigBytes = append([]byte(header), v1ConfigBytes...)
	if _, err := tmpFile.Write(v1ConfigBytes); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to write new config file: %w", err),
			tmpFile.Close(),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := tmpFile.Close(); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to close new config file: %w", err),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := os.Rename(tmpFile.Name(), configPath); err != nil {
		return fmt.Errorf("failed to overwrite old config: %w", err)
	}
	return nil
}

func (m *v1beta1Migrator) maybeMigrateGenTemplate(dirPath string) (bool, error) {
	oldConfigPath := filepath.Join(dirPath, bufgen.ExternalConfigFilePath)
	oldConfigBytes, err := os.ReadFile(oldConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// OK, no old config file
			return false, nil
		}
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	var versionedConfig bufgen.ExternalConfigVersion
	if err := encoding.UnmarshalYAMLNonStrict(oldConfigBytes, &versionedConfig); err != nil {
		return false, fmt.Errorf(
			"failed to read %s version: %w",
			oldConfigPath,
			err,
		)
	}
	switch versionedConfig.Version {
	case bufgen.V1Version:
		// OK, file was already v1
		return false, nil
	case bufgen.V1Beta1Version, "":
		// Continue to migrate
	default:
		return false, fmt.Errorf("unknown config file version: %s", versionedConfig.Version)
	}
	var v1beta1GenTemplate bufgen.ExternalConfigV1Beta1
	if err := encoding.UnmarshalYAMLStrict(oldConfigBytes, &v1beta1GenTemplate); err != nil {
		return false, fmt.Errorf(
			"failed to unmarshal %s as %s version v1beta1: %w",
			oldConfigPath,
			bufgen.ExternalConfigFilePath,
			err,
		)
	}
	v1GenTemplate := bufgen.ExternalConfigV1{
		Version: bufgen.V1Version,
		Managed: bufgen.ExternalManagedConfigV1{
			Enabled:           v1beta1GenTemplate.Managed,
			CcEnableArenas:    v1beta1GenTemplate.Options.CcEnableArenas,
			JavaMultipleFiles: v1beta1GenTemplate.Options.JavaMultipleFiles,
			OptimizeFor:       v1beta1GenTemplate.Options.OptimizeFor,
		},
	}
	for _, plugin := range v1beta1GenTemplate.Plugins {
		v1GenTemplate.Plugins = append(v1GenTemplate.Plugins, bufgen.ExternalPluginConfigV1(plugin))
	}
	newConfigPath := filepath.Join(dirPath, bufgen.ExternalConfigFilePath)
	if err := m.writeV1GenTemplate(newConfigPath, v1GenTemplate); err != nil {
		return false, err
	}
	return true, nil
}

// writeV1GenTemplate atomically replaces the old configuration file by first writing
// the new config to a temporary file and then moving it to the old config file path.
// If we fail to marshal or write, the old config file is not touched.
func (m *v1beta1Migrator) writeV1GenTemplate(
	configPath string,
	config bufgen.ExternalConfigV1,
) (retErr error) {
	v1ConfigBytes, err := encoding.MarshalYAML(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal new config: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for writing new config: %w", err)
	}
	header := fmt.Sprintf(bufGenHeader, m.commandName)
	v1ConfigBytes = append([]byte(header), v1ConfigBytes...)
	if _, err := tmpFile.Write(v1ConfigBytes); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to write new config file: %w", err),
			tmpFile.Close(),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := tmpFile.Close(); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to close new config file: %w", err),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := os.Rename(tmpFile.Name(), configPath); err != nil {
		return fmt.Errorf("failed to overwrite old config: %w", err)
	}
	return nil
}

func (m *v1beta1Migrator) maybeMigrateLockFile(dirPath string) (bool, error) {
	oldConfigPath := filepath.Join(dirPath, buflock.ExternalConfigFilePath)
	oldConfigBytes, err := os.ReadFile(oldConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// OK, no old config file
			return false, nil
		}
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	var versionedConfig buflock.ExternalConfigVersion
	if err := encoding.UnmarshalYAMLNonStrict(oldConfigBytes, &versionedConfig); err != nil {
		return false, fmt.Errorf(
			"failed to read %s version: %w",
			oldConfigPath,
			err,
		)
	}
	switch versionedConfig.Version {
	case buflock.V1Version:
		// OK, file was already v1
		return false, nil
	case buflock.V1Beta1Version, "":
		// Continue to migrate
	default:
		return false, fmt.Errorf("unknown lock file version: %s", versionedConfig.Version)
	}
	var v1beta1LockFile buflock.ExternalConfigV1Beta1
	if err := encoding.UnmarshalYAMLStrict(oldConfigBytes, &v1beta1LockFile); err != nil {
		return false, fmt.Errorf(
			"failed to unmarshal %s as %s version v1beta1: %w",
			oldConfigPath,
			buflock.ExternalConfigFilePath,
			err,
		)
	}
	v1LockFile := buflock.ExternalConfigV1{
		Version: buflock.V1Version,
	}
	for _, dependency := range v1beta1LockFile.Deps {
		v1LockFile.Deps = append(v1LockFile.Deps, buflock.ExternalConfigDependencyV1(dependency))
	}
	newConfigPath := filepath.Join(dirPath, buflock.ExternalConfigFilePath)
	if err := m.writeV1LockFile(newConfigPath, v1LockFile); err != nil {
		return false, err
	}
	return true, nil
}

// writeV1LockFile atomically replaces the old lock file by first writing
// the new lock file to a temporary file and then moving it to the old lock file path.
// If we fail to marshal or write, the old lock file is not touched.
func (m *v1beta1Migrator) writeV1LockFile(
	configPath string,
	config buflock.ExternalConfigV1,
) (retErr error) {
	v1ConfigBytes, err := encoding.MarshalYAML(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal new lock file: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for writing new lock file: %w", err)
	}
	v1ConfigBytes = append([]byte(buflock.Header), v1ConfigBytes...)
	if _, err := tmpFile.Write(v1ConfigBytes); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to write new lock file: %w", err),
			tmpFile.Close(),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := tmpFile.Close(); err != nil {
		return multierr.Combine(
			fmt.Errorf("failed to close new lock file: %w", err),
			os.Remove(tmpFile.Name()),
		)
	}
	if err := os.Rename(tmpFile.Name(), configPath); err != nil {
		return fmt.Errorf("failed to overwrite old lock file: %w", err)
	}
	return nil
}

func maybeReadLockFile(oldLockFilePath string) (buflock.ExternalConfigV1, bool, error) {
	lockFileBytes, err := os.ReadFile(oldLockFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// OK, no old lock file
			return buflock.ExternalConfigV1{}, false, nil
		}

		return buflock.ExternalConfigV1{}, false, fmt.Errorf("failed to read lock file path: %w", err)
	}
	var versionedConfig buflock.ExternalConfigVersion
	if err := encoding.UnmarshalYAMLNonStrict(lockFileBytes, &versionedConfig); err != nil {
		return buflock.ExternalConfigV1{}, false, fmt.Errorf(
			"failed to read %s version: %w",
			oldLockFilePath,
			err,
		)
	}
	switch versionedConfig.Version {
	case "", buflock.V1Beta1Version:
		var externalConfig buflock.ExternalConfigV1Beta1
		if err := encoding.UnmarshalYAMLStrict(lockFileBytes, &externalConfig); err != nil {
			return buflock.ExternalConfigV1{}, false, fmt.Errorf(
				"failed to unmarshal lock file at %s: %w",
				buflock.V1Beta1Version,
				err,
			)
		}
		externalLockFileV1 := buflock.ExternalConfigV1{
			Version: buflock.V1Version,
		}
		for _, dependency := range externalConfig.Deps {
			externalLockFileV1.Deps = append(externalLockFileV1.Deps, buflock.ExternalConfigDependencyV1(dependency))
		}
		return externalLockFileV1, true, nil
	case buflock.V1Version:
		externalLockFileV1 := buflock.ExternalConfigV1{}
		if err := encoding.UnmarshalYAMLStrict(lockFileBytes, &externalLockFileV1); err != nil {
			return buflock.ExternalConfigV1{}, false, fmt.Errorf("failed to unmarshal lock file at %s: %w", buflock.V1Version, err)
		}
		return externalLockFileV1, true, nil
	default:
		return buflock.ExternalConfigV1{}, false, fmt.Errorf("unknown lock file version: %s", versionedConfig.Version)
	}
}

func convertIgnoreSlice(paths []string, dirPath string, root string, pathToProcessed map[string]bool) ([]string, error) {
	var ignoresForRoot []string
	for _, ignoredFile := range paths {
		if _, ok := pathToProcessed[ignoredFile]; !ok {
			pathToProcessed[ignoredFile] = false
		}
		filePath := filepath.Join(dirPath, root, ignoredFile)
		if _, err := os.Stat(filePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("failed to check for presence of file %s: %w", filePath, err)
		}
		pathToProcessed[ignoredFile] = true
		ignoresForRoot = append(ignoresForRoot, ignoredFile)
	}
	sort.Strings(ignoresForRoot)
	return ignoresForRoot, nil
}

func convertIgnoreMap(ruleToIgnores map[string][]string, dirPath string, root string, pathToProcessed map[string]bool) (map[string][]string, error) {
	var ruleToIgnoresForRoot map[string][]string
	for rule, ignores := range ruleToIgnores {
		for _, ignoredFile := range ignores {
			if _, ok := pathToProcessed[ignoredFile]; !ok {
				pathToProcessed[ignoredFile] = false
			}
			filePath := filepath.Join(dirPath, root, ignoredFile)
			if _, err := os.Stat(filePath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return nil, fmt.Errorf("failed to check for presence of file %s: %w", filePath, err)
			}
			if ruleToIgnoresForRoot == nil {
				ruleToIgnoresForRoot = make(map[string][]string)
			}
			pathToProcessed[ignoredFile] = true
			ruleToIgnoresForRoot[rule] = append(
				ruleToIgnoresForRoot[rule],
				ignoredFile,
			)
		}
		sort.Strings(ruleToIgnoresForRoot[rule])
	}
	return ruleToIgnoresForRoot, nil
}
