// Copyright 2020 Buf Technologies Inc.
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

package fetch

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/bufbuild/buf/internal/pkg/app"
	"github.com/bufbuild/buf/internal/pkg/ioutilextended"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type writer struct {
	logger *zap.Logger
}

func newWriter(
	logger *zap.Logger,
	options ...WriterOption,
) *writer {
	writer := &writer{
		logger: logger,
	}
	for _, option := range options {
		option(writer)
	}
	return writer
}

func (w *writer) PutFile(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
	options ...PutFileOption,
) (io.WriteCloser, error) {
	putFileOptions := newPutFileOptions()
	for _, option := range options {
		option(putFileOptions)
	}
	switch t := fileRef.(type) {
	case SingleRef:
		return w.putSingle(
			ctx,
			container,
			t,
			putFileOptions.noFileCompression,
		)
	case ArchiveRef:
		return w.putArchiveFile(
			ctx,
			container,
			t,
			putFileOptions.noFileCompression,
		)
	default:
		return nil, fmt.Errorf("unknown FileRef type: %T", fileRef)
	}
}

func (w *writer) putSingle(
	ctx context.Context,
	container app.EnvStdoutContainer,
	singleRef SingleRef,
	noFileCompression bool,
) (io.WriteCloser, error) {
	return w.putFileWriteCloser(ctx, container, singleRef, noFileCompression)
}

func (w *writer) putArchiveFile(
	ctx context.Context,
	container app.EnvStdoutContainer,
	archiveRef ArchiveRef,
	noFileCompression bool,
) (io.WriteCloser, error) {
	return w.putFileWriteCloser(ctx, container, archiveRef, noFileCompression)
}

func (w *writer) putFileWriteCloser(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
	noFileCompression bool,
) (_ io.WriteCloser, retErr error) {
	writeCloser, err := w.putFileWriteCloserPotentiallyUncompressed(ctx, container, fileRef)
	if err != nil {
		return nil, err
	}
	defer func() {
		if retErr != nil {
			retErr = multierr.Append(retErr, writeCloser.Close())
		}
	}()
	if noFileCompression {
		return writeCloser, nil
	}
	switch compressionType := fileRef.CompressionType(); compressionType {
	case CompressionTypeNone:
		return writeCloser, nil
	case CompressionTypeGzip:
		gzipWriteCloser := gzip.NewWriter(writeCloser)
		return ioutilextended.CompositeWriteCloser(
			gzipWriteCloser,
			ioutilextended.ChainCloser(
				gzipWriteCloser,
				writeCloser,
			),
		), nil
	default:
		return nil, fmt.Errorf("unknown CompressionType: %v", compressionType)
	}
}

func (w *writer) putFileWriteCloserPotentiallyUncompressed(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
) (io.WriteCloser, error) {
	switch fileScheme := fileRef.FileScheme(); fileScheme {
	case FileSchemeHTTP:
		return nil, fmt.Errorf("http not supported for writes: %v", fileRef.Path())
	case FileSchemeHTTPS:
		return nil, fmt.Errorf("https not supported for writes: %v", fileRef.Path())
	case FileSchemeLocal:
		return os.Create(fileRef.Path())
	case FileSchemeStdio:
		return ioutilextended.NopWriteCloser(container.Stdout()), nil
	case FileSchemeNull:
		return ioutilextended.DiscardWriteCloser, nil
	default:
		return nil, fmt.Errorf("unknown FileScheme: %v", fileScheme)
	}
}

type putFileOptions struct {
	noFileCompression bool
}

func newPutFileOptions() *putFileOptions {
	return &putFileOptions{}
}
