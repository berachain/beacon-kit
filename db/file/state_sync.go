package file

import (
	"errors"
	"io"

	snapshot "cosmossdk.io/store/snapshots/types"
)

const (
	txHashSize = 32
	ttlSize    = 8
	chunkSize  = txHashSize + ttlSize
)

var _ snapshot.ExtensionSnapshotter = &Snapshotter{}

const (
	// SnapshotFormat defines the snapshot format of exported blobs.
	// No protobuf envelope, no metadata.
	SnapshotFormat = 1

	// SnapshotName defines the snapshot name of exported blobs.
	SnapshotName = "blob"
)

type Snapshotter struct {
	m *DB
}

func NewSnapshotter(m *DB) *Snapshotter {
	return &Snapshotter{m: m}
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

// SnapshotExtension exports the state at a given height using the provided snapshotWriter.
func (s *Snapshotter) SnapshotExtension(height uint64, payloadWriter snapshot.ExtensionPayloadWriter) error {
	// export all blobs as a single blob
	return s.m.exportSnapshot(height, payloadWriter)
}

func (s *Snapshotter) RestoreExtension(height uint64, format uint32, payloadReader snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		return s.restore(height, payloadReader)
	}

	return snapshot.ErrUnknownFormat
}

// restore restores the state at a given height using the provided payloadReader.
func (s *Snapshotter) restore(height uint64, payloadReader snapshot.ExtensionPayloadReader) error {
	_, err := payloadReader()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return io.ErrUnexpectedEOF
		}

		return err
	}

	// TODO: restore the blob

	return nil
}
