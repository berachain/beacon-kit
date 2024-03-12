package beacon

import (
	"cosmossdk.io/collections"
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
)

type validatorsIndex struct {
	Pubkey *indexes.Unique[[]byte, uint64, []byte]
}

func (a validatorsIndex) IndexesList() []collections.Index[uint64, []byte] {
	return []collections.Index[uint64, []byte]{a.Pubkey}
}

func NewValidatorsIndex(sb *collections.SchemaBuilder) validatorsIndex {
	return validatorsIndex{
		Pubkey: indexes.NewUnique(
			sb, sdkcollections.NewPrefix(validatrPubkeyToIndexPrefix), validatrPubkeyToIndexPrefix,
			sdkcollections.BytesKey, sdkcollections.Uint64Key,
			func(_ uint64, pubkey []byte) ([]byte, error) {
				return pubkey, nil
			},
		),
	}
}
