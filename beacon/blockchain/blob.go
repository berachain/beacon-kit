package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"sync"

	"github.com/berachain/beacon-kit/db"
)

// Create a pool of bytes.Buffers.
var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Store the blobs in the blobstore.
func PrepareBlobsHandler(storage db.BeaconKitDB,
	height int64, blobs [][48]byte) ([]byte, error) {

	// TODO: before storage handle validation
	// Encode blobs to bytes
	buf, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		buf = new(bytes.Buffer)
	}
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	enc := gob.NewEncoder(buf)
	err := enc.Encode(blobs)
	if err != nil {
		return nil, err
	}
	encodedData := buf.Bytes()

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))
	// store the blobs under a single height.
	if err := storage.Set(heightBytes, encodedData); err != nil {
		return nil, err
	}

	return encodedData, nil
}

// Store the blobs in the blobstore.
func ProcessBlobsHandler(storage db.BeaconKitDB,
	height int64, blobTx []byte) error {

	// Decode the blobs from bytes to []byte
	var blobs [][]byte
	buf, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		buf = new(bytes.Buffer)
	}
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.Write(blobTx)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&blobs)
	if err != nil {
		return err
	}

	// TODO: before storage handle validation

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))
	// Store the blobs under a single height.
	if err := storage.Set(heightBytes, blobTx); err != nil {
		return err
	}

	return nil
}
