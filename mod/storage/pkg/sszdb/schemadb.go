// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package sszdb

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	fastssz "github.com/ferranbt/fastssz"
)

// bootstrappedKey is a temporary key to check if the database has been
// bootstrapped.  it will be made obsolete by versioning.
const bootstrappedKey = "bootstrapped"
const bytesPerNode = 32

type ObjectPath = merkle.ObjectPath[uint64, [32]byte]

func isByteVector(typ schema.SSZType) bool {
	return typ.ID() == schema.Vector && typ.ElementType("") == schema.U8()
}

type SchemaDB struct {
	*Backend
	schemaRoot schema.SSZType
}

func NewSchemaDB(
	backend *Backend,
	monolith Treeable,
) (*SchemaDB, error) {
	schemaRoot, err := schema.Build(monolith)
	if err != nil {
		return nil, err
	}
	db := &SchemaDB{
		Backend:    backend,
		schemaRoot: schemaRoot,
	}
	return db, db.bootstrap(monolith)
}

func (db *SchemaDB) bootstrap(monolith Treeable) error {
	bootstrapped, err := db.get([]byte(bootstrappedKey))
	if err != nil {
		return err
	}
	if bootstrapped != nil {
		return nil
	}
	err = db.SaveMonolith(monolith)
	if err != nil {
		return err
	}
	return db.Set([]byte(bootstrappedKey), []byte{1})
}

type offsetBytes struct {
	bz  []byte
	idx uint32
}

func (db *SchemaDB) getLeafBytes(
	ctx context.Context,
	path ObjectPath,
) ([]byte, error) {
	typ, gindex, offset, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return nil, err
	}
	size := typ.ItemLength()

	// if the path was to a _byte vector_ unmarshal all of its leaves
	if isByteVector(typ) {
		_, gindex, offset, err = path.Append("0").
			GetGeneralizedIndex(db.schemaRoot)
		if err != nil {
			return nil, err
		}
		// set size to the length of byte vector
		size = typ.Length()
	}

	return db.getNodeBytes(ctx, gindex, size, offset)
}

func (db *SchemaDB) getSSZBytes(
	ctx context.Context,
	root ObjectPath,
) (uint32, *offsetBytes, []byte, error) {
	var (
		offsets  []*offsetBytes
		n        uint32
		sszBytes []byte
		bz       []byte
	)
	typ, _, _, err := root.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return 0, nil, nil, err
	}

	switch typ.ID() {
	case schema.Basic:
		bz, err = db.getLeafBytes(ctx, root)
		if err != nil {
			return 0, nil, nil, err
		}
		sszBytes = append(sszBytes, bz...)
		n += uint32(typ.ItemLength())
		return n, nil, sszBytes, nil
	case schema.Vector, schema.List:
		if isByteVector(typ) {
			bz, err = db.getLeafBytes(ctx, root)
			if err != nil {
				return 0, nil, nil, err
			}
			sszBytes = append(sszBytes, bz...)
			n += uint32(typ.Length())
			return n, nil, sszBytes, nil
		} else if typ.ID().IsList() {
			bz, err = db.getLeafBytes(ctx, root.Append("__len__"))
			if err != nil {
				return 0, nil, nil, err
			}
			length := fastssz.UnmarshallUint64(bz)

			var offsetBz []byte
			for i := range length {
				// TODO: list of dynamic elements not yet supported
				// result not guaranteed to be correct
				_, _, bz, err = db.getSSZBytes(ctx, root.Append(fmt.Sprintf("%d", i)))
				if err != nil {
					return 0, nil, nil, err
				}
				offsetBz = append(offsetBz, bz...)
			}
			// write empty offset address
			sszBytes = append(sszBytes, make([]byte, 4)...)
			n += 4
			return n, &offsetBytes{bz: offsetBz}, sszBytes, nil
		}
	case schema.Container:
		// TODO assumes fixed size container
		paths := make([]ObjectPath, typ.Length())
		for i, p := range schema.ContainerFields(typ) {
			paths[i] = root.Append(p)
		}
		for _, p := range paths {
			size, off, bz, err := db.getSSZBytes(ctx, p)
			if err != nil {
				return 0, nil, nil, err
			}
			sszBytes = append(sszBytes, bz...)
			if off != nil {
				off.idx = n
				offsets = append(offsets, off)
			}
			n += size
		}
	}

	for _, o := range offsets {
		// write offset address
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, n)
		copy(sszBytes[o.idx:], buf)
		sszBytes = append(sszBytes, o.bz...)
	}

	return n, nil, sszBytes, nil
}

func (db *SchemaDB) GetPath(
	ctx context.Context,
	path ObjectPath,
) ([]byte, error) {
	_, offsetBz, bz, err := db.getSSZBytes(ctx, path)
	if offsetBz != nil {
		return offsetBz.bz, err
	}
	return bz, err
}

func (db *SchemaDB) GetObject(
	ctx context.Context,
	path ObjectPath,
	obj constraints.SSZUnmarshaler,
) error {
	bz, err := db.GetPath(ctx, path)
	if err != nil {
		return err
	}
	return obj.UnmarshalSSZ(bz)
}

func (db *SchemaDB) SetObject(
	ctx context.Context,
	path ObjectPath,
	obj Treeable,
) error {
	treeNode, err := NewTreeFromFastSSZ(obj)
	if err != nil {
		return err
	}
	_, gidx, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	return db.stage(ctx, treeNode, gidx)
}

func (db *SchemaDB) SetRaw(
	ctx context.Context,
	path ObjectPath,
	bz []byte,
) error {
	if len(bz) > bytesPerNode {
		return fmt.Errorf(
			"expected max %d bytes, got %d",
			bytesPerNode,
			len(bz),
		)
	}
	if len(bz) < bytesPerNode {
		// pad with zeros
		bz = append(bz, make([]byte, bytesPerNode-len(bz))...)
	}
	_, gidx, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	return db.stage(ctx, &Node{Value: bz}, gidx)
}

func (db *SchemaDB) GetListLength(
	ctx context.Context,
	path ObjectPath,
) (uint64, error) {
	op := path.Append("__len__")
	bz, err := db.GetPath(ctx, op)
	if err != nil {
		return 0, err
	}
	return fastssz.UnmarshallUint64(bz), nil
}

func (db *SchemaDB) setListLength(
	ctx context.Context,
	path ObjectPath,
	length uint64,
) error {
	op := path.Append("__len__")
	_, gindex, _, err := op.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	val := make([]byte, 32)
	binary.LittleEndian.PutUint64(val, length)
	return db.stage(
		ctx,
		&Node{Value: val},
		gindex,
	)
}

func (db *SchemaDB) SetListElementRaw(
	ctx context.Context,
	path ObjectPath,
	index uint64,
	bz []byte,
) error {
	length := len(bz)
	if length > 32 {
		return fmt.Errorf("expected max 32 bytes, got %d", len(bz))
	}
	if length == 32 {
		return db.setListElement(ctx, path, index, &Node{Value: bz})
	}

	objPath := ObjectPath(fmt.Sprintf("%s/%d", path, index))
	_, gindex, offset, err := objPath.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	nodeBz, ok, err := db.getNode(ctx, gindex)
	if err != nil {
		return err
	}
	if !ok {
		if index > 0 {
			return fmt.Errorf(
				"attempted to set list element %s/%d but node not found gindex=%d",
				objPath,
				index,
				gindex,
			)
		}
		nodeBz = make([]byte, 32)
	}
	copy(nodeBz[offset:], bz)
	return db.setListElement(ctx, path, index, &Node{Value: nodeBz})
}

func (db *SchemaDB) SetListElementObject(
	ctx context.Context,
	path ObjectPath,
	index uint64,
	obj Treeable,
) error {
	treeNode, err := NewTreeFromFastSSZ(obj)
	if err != nil {
		return err
	}
	return db.setListElement(ctx, path, index, treeNode)
}

func (db *SchemaDB) setListElement(
	ctx context.Context,
	path ObjectPath,
	index uint64,
	node *Node,
) error {
	length, err := db.GetListLength(ctx, path)
	if err != nil {
		return err
	}
	if index > length {
		return fmt.Errorf("index %d out of bounds; len=%d", index, length)
	}
	objPath := ObjectPath(fmt.Sprintf("%s/%d", path, index))
	_, gidx, _, err := objPath.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	err = db.stage(ctx, node, gidx)
	if err != nil {
		return err
	}

	if index != length {
		return nil
	}

	// when the index is at the end of the list, we need to update the length
	// and potentially add some zero hashes
	gindex := gidx
	depth := 0
	branchID := db.stageID(ctx)
	for gindex > 1 {
		if gindex%2 == 0 {
			var ok bool
			_, ok, err = db.getNode(ctx, gindex+1)
			if err != nil {
				return err
			}
			if ok {
				// exit condition: once pre-existing sibling is found
				// upward traversal can be stopped
				break
			}
			db.stages[branchID][gindex+1] = db.zeroHashes[depth]
		}
		depth++
		gindex /= 2
	}
	return db.setListLength(ctx, path, index+1)
}
