package cms

import (
	"cosmossdk.io/store/types"
)

type CustomRootMultistore struct {
	types.CommitMultiStore
}

// LastCommitID implements Committer/CommitStore.
func (rs *CustomRootMultistore) Commit() types.CommitID {
	rs.CommitMultiStore.Commit()
	return types.CommitID{}
}
