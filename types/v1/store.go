package v1

type ForkChoiceStore interface {
	SetSafeBlockHash(safeBlockHash [32]byte) error
	GetSafeBlockHash() ([32]byte, error)
	SetFinalizedBlockHash(finalizedBlockHash [32]byte) error
	GetFinalizedBlockHash() ([32]byte, error)
	SetLastValidHead(lastValidHead [32]byte) error
	GetLastValidHead() [32]byte
}
