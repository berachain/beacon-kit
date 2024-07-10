package runtime

import "time"

// ABCIRequest represents the interface for an ABCI request.
type ABCIRequest interface {
	// GetHeight returns the height of the request.
	GetHeight() int64
	// GetTime returns the time of the request.
	GetTime() time.Time
	// GetTxs returns the transactions included in the request.
	GetTxs() [][]byte
}
