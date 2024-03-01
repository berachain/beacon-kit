package randao

import "github.com/itsdevbear/bolaris/types/consensus/primitives"

//go:generate sszgen -path . -objs MySSZType --include ../../../types/consensus/primitives

type MySSZType struct {
	MyFirstField  []byte `ssz-size:"96"`
	MySecondField primitives.SSZUint64
}
