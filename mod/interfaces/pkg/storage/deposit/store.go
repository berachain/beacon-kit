package deposit

// Store is a simple KV store based implementation that assumes
// the deposit indexes are tracked outside of the kv store.
type Store[DepositT any] interface {
	// GetDepositsByIndex returns the first N deposits starting from the given
	// index. If N is greater than the number of deposits, it returns up to the
	// last deposit.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]DepositT, error)
	// EnqueueDeposit pushes a deposit to the queue.
	EnqueueDeposit(deposit DepositT) error
	// EnqueueDeposits pushes multiple deposits to the queue.
	EnqueueDeposits(deposits []DepositT) error
	// Prune removes the [start, end) deposits from the store.
	Prune(start, end uint64) error
}
