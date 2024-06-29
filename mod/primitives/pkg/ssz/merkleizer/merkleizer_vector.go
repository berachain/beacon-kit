package merkleizer

// MerkleizeVecBasic implements the SSZ merkleization algorithm
// for a vector of basic types.
func (m *merkleizer[RootT, T]) MerkleizeVecBasic(
	value []T,
) (RootT, error) {
	packed, err := m.pack(value)
	if err != nil {
		return [32]byte{}, err
	}
	return m.Merkleize(packed)
}

// MerkleizeVecComposite implements the SSZ merkleization algorithm for a vector
// of composite types.
func (m *merkleizer[RootT, T]) MerkleizeVecComposite(
	value []T,
) (RootT, error) {
	var (
		err  error
		htrs = m.bytesBuffer.Get(len(value))
	)

	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}
	return m.Merkleize(htrs)
}
