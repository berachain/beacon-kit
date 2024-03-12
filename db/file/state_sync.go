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

const (
	// SnapshotFormat defines the snapshot format of exported blobs.
	// No protobuf envelope, no metadata.
	SnapshotFormat = 1

	// SnapshotName defines the snapshot name of exported blobs.
	SnapshotName = "blob"
)

type Snapshotter[T numeric] struct {
	db             *DB
	snapShotWindow T
}

func NewSnapshotter[T numeric](m *DB, sw T) *Snapshotter[T] {
	return &Snapshotter[T]{db: m, snapShotWindow: sw}
}

func (s *Snapshotter[T]) SnapshotName() string {
	return SnapshotName
}

func (s *Snapshotter[T]) SnapshotFormat() uint32 {
	return SnapshotFormat
}

func (s *Snapshotter[T]) SupportedFormats() []uint32 {
	return []uint32{SnapshotFormat}
}

// SnapshotExtension exports the state
// of the snapshot window given snapshotWriter.
func (s *Snapshotter[T]) SnapshotExtension(height uint64,
	payloadWriter snapshot.ExtensionPayloadWriter) error {
	// export all blobs as a single blob
	exportBlocks := height - uint64(s.snapShotWindow)

	ranger := NewRangeDB[uint64](s.db)
	// TODO: add iteration for the file system storage
	for i := exportBlocks; i < height; i++ {
		// load code and abort on error
		bytes, err := ranger.Get(height, []byte("commitment"))
		if err != nil {
			return err
		}

		err = payloadWriter(bytes)
		if err != nil {
			return err
		}

	}

	return nil
}

/*
loop through all the blob files create a chunk pre file send it out

prune = 100 blocks (100 chunks)
*/

func (s *Snapshotter[T]) RestoreExtension(height uint64, format uint32,
	payloadReader snapshot.ExtensionPayloadReader) error {
	if format == SnapshotFormat {
		return s.restore(height, payloadReader)
	}

	return snapshot.ErrUnknownFormat
}

// restore restores the state at a given height
// using the provided payloadReader.
func (s *Snapshotter[T]) restore(height uint64,
	payloadReader snapshot.ExtensionPayloadReader) error {
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
