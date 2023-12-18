package store

import (
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"

	"cosmossdk.io/store"
)

type Forkchoice struct {
	store store.KVStore
}

func NewForkchoice(store store.KVStore) *Forkchoice {
	return &Forkchoice{
		store: store,
	}
}

func (f *Forkchoice) Store(fcs *enginev1.ForkchoiceState) error {
	bz, err := proto.Marshal(fcs)
	if err != nil {
		return err
	}
	f.store.Set([]byte("forkchoice"), bz)
	return nil
}

func (f *Forkchoice) Retrieve() (*enginev1.ForkchoiceState, error) {
	bz := f.store.Get([]byte("forkchoice"))
	fcs := &enginev1.ForkchoiceState{}
	if err := proto.Unmarshal(bz, fcs); err != nil {
		return nil, err
	}
	return fcs, nil
}
