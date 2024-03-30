package comet

import (
	"fmt"
	"regexp"

	"github.com/cometbft/cometbft/crypto/merkle"
	lrpc "github.com/cometbft/cometbft/light/rpc"
)

// DefaultMerkleKeyPathFn creates a function used to generate merkle key paths
// from a path string and a key. This is the default used by the cosmos SDK.
// This merkle key paths are required when verifying /abci_query calls
func DefaultMerkleKeyPathFn() lrpc.KeyPathFunc {
	// regexp for extracting store name from /abci_query path
	storeNameRegexp := regexp.MustCompile(`store\/(.+)\/key`)

	return func(path string, key []byte) (merkle.KeyPath, error) {
		matches := storeNameRegexp.FindStringSubmatch(path)
		if len(matches) != 2 {
			return nil, fmt.Errorf("can't find store name in %s using %s", path, storeNameRegexp)
		}
		storeName := matches[1]

		kp := merkle.KeyPath{}
		kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
		kp = kp.AppendKey(key, merkle.KeyEncodingURL)
		return kp, nil
	}
}
