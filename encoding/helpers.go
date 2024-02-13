package encoding

// chunkSize is the size of each chunk in bytes
const chunkSize = 32

// ByteSlice is a 32 byte array
type ByteSlice [chunkSize]byte

// ByteSliceRoot is a helper func to merkleize an arbitrary List[Byte, N]
// this func runs Chunkify + MerkleizeVector
// max length is dividable by 32 ( root length )
// func ByteSliceRoot(slice []byte, maxLength uint64) (ByteSlice, error) {
// 	chunkedRoots, err := PackByChunk([][]byte{slice})
// 	if err != nil {
// 		return ByteSlice{}, err
// 	}
// 	// nearest number divisible by root length (32)
// 	maxRootLength := (maxLength + chunkSize - 1) / chunkSize

// 	// TODO: Uncomment this
// 	// bytesRoot, err := BitwiseMerkleize(chunkedRoots, uint64(len(chunkedRoots)), maxRootLength)
// 	bytesRoot := ByteSlice{}

// 	if err != nil {
// 		return ByteSlice{}, errors.Wrap(err, "could not compute merkleization")
// 	}
// 	bytesRootBuf := new(bytes.Buffer)
// 	if err := binary.Write(bytesRootBuf, binary.LittleEndian, uint64(len(slice))); err != nil {
// 		return ByteSlice{}, errors.Wrap(err, "could not marshal length")
// 	}
// 	bytesRootBufRoot := make([]byte, chunkSize)
// 	copy(bytesRootBufRoot, bytesRootBuf.Bytes())
// 	return MixInLength(bytesRoot, bytesRootBufRoot), nil
// }

// PackByChunk a given byte array's final chunk with zeroes if needed.
// func PackByChunk(serializedItems [][]byte) ([]ByteSlice, error) {
// 	// Return early if no items are provided.
// 	if len(serializedItems) == 0 {
// 		return []ByteSlice{}, nil
// 		// If each item has the same chunk length, return the serialized items.
// 	} else if len(serializedItems[0]) == chunkSize {
// 		chunks := make([]ByteSlice, len(serializedItems))
// 		for i, c := range serializedItems {
// 			chunks[i] = bytesutil.ToBytes32(c)
// 		}
// 		return chunks, nil
// 	}

// 	// We flatten the list in order to pack its items into byte chunks correctly.
// 	orderedItems := make([]byte, 0, len(serializedItems)*len(serializedItems[0]))
// 	for _, item := range serializedItems {
// 		orderedItems = append(orderedItems, item...)
// 	}

// 	return chunkItems(orderedItems), nil
// }

// Assuming orderedItems is a slice of some type that needs to be chunked
// and ToBytes32 is a function that takes a slice of this type and returns
// a [32]byte array, possibly right-padding with zeros.
//
// chunkItems slices the orderedItems into chunks of chunkSize,
// right-padding the last chunk with zero bytes if necessary.
// func chunkItems(orderedItems []byte) []ByteSlice {
// 	var chunks []ByteSlice
// 	numItems := len(orderedItems)

// 	// If all our serialized item slices are length zero, we exit early.
// 	if numItems == 0 {
// 		return []ByteSlice{}
// 	}

// 	for i := 0; i < numItems; i += chunkSize {
// 		end := i + chunkSize
// 		if end > numItems {
// 			end = numItems
// 		}

// 		// Assuming ToBytes32 can handle slices smaller than 32 by right-padding them with zeros.
// 		chunk := bytesutil.ToBytes32(orderedItems[i:end])
// 		chunks = append(chunks, chunk)
// 	}

// 	return chunks
// }

// MixInLength appends hash length to root
// func MixInLength(root ByteSlice, length []byte) ByteSlice {
// 	var hash ByteSlice
// 	h := sha256.New()
// 	h.Write(root[:])
// 	h.Write(length)
// 	// The hash interface never returns an error, for that reason
// 	// we are not handling the error below. For reference, it is
// 	// stated here https://golang.org/pkg/hash/#Hash
// 	// #nosec G104
// 	h.Sum(hash[:0])
// 	return hash
// }
