// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	snapshot "cosmossdk.io/store/snapshots/types"
	"github.com/spf13/afero"
)

const (
	// SnapshotFormat defines the snapshot format of exported blobs.
	// No protobuf envelope, no metadata.
	SnapshotFormat = 1

	// SnapshotName defines the snapshot name of exported blobs.
	SnapshotName = "blob"
)

type Snapshotter struct {
	db             *DB
	snapShotWindow uint64
}

func NewSnapshotter(m *DB, sw uint64) *Snapshotter {
	return &Snapshotter{db: m, snapShotWindow: sw}
}

func (s *Snapshotter) SnapshotName() string {
	return SnapshotName
}

func (s *Snapshotter) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (s *Snapshotter) SupportedFormats() []uint32 {
	return []uint32{SnapshotFormat}
}

// SnapshotExtension exports the state
// of the snapshot window given snapshotWriter.
//
//nolint:gocognit // todo:reee.
func (s *Snapshotter) SnapshotExtension(height uint64,
	payloadWriter snapshot.ExtensionPayloadWriter) error {
	if err := afero.Walk(s.db.fs, s.db.rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".blob") {
				return nil
			}

			// Extract the filename from the path
			_, filename := filepath.Split(path)

			// Use a regular expression to find numbers in the filename
			re, err := ExtractIndex([]byte(filename))
			if err != nil {
				return err
			}

			if s.snapShotWindow > height {
				value, err1 := afero.ReadFile(s.db.fs, path)
				if err1 != nil {
					return err1
				}

				prefixedData := append([]byte(filename+"\n"), value...)

				if err1 = payloadWriter(prefixedData); err1 != nil {
					return err1
				}
			} else if re >= height-s.snapShotWindow && re <= height {
				value, err1 := afero.ReadFile(s.db.fs, path)
				if err1 != nil {
					return err1
				}

				prefixedData := append([]byte(filename+"\n"), value...)

				if err1 = payloadWriter(prefixedData); err1 != nil {
					return err1
				}
			}

			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (s *Snapshotter) RestoreExtension(height uint64, format uint32,
	payloadReader snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		return s.restore(height, payloadReader)
	}

	return snapshot.ErrUnknownFormat
}

// restore restores the state at a given height
// using the provided payloadReader.
func (s *Snapshotter) restore(_ uint64,
	payloadReader snapshot.ExtensionPayloadReader) error {
	for {
		bz, err := payloadReader()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}

			return err
		}

		split := bytes.SplitN(bz, []byte("\n"), two)
		if len(split) != two {
			return errors.New("invalid blob format")
		}
		receivedFilename := string(split[0])
		receivedData := split[1]

		file, err := s.db.fs.Create(receivedFilename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		_, err = file.Write(receivedData)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
}
