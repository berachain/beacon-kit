package spec

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/stretchr/testify/require"
)

func TestDomainTypeConversion(t *testing.T) {
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeProposer))
	require.Equal(t, bytes.B4([]byte{0x01, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeAttester))
	require.Equal(t, bytes.B4([]byte{0x02, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeRandao))
	require.Equal(t, bytes.B4([]byte{0x03, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeDeposit))
	require.Equal(t, bytes.B4([]byte{0x04, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeVoluntaryExit))
	require.Equal(t, bytes.B4([]byte{0x05, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeSelectionProof))
	require.Equal(t, bytes.B4([]byte{0x06, 0x00, 0x00, 0x00}), bytes.FromUint32(defaultDomainTypeAggregateAndProof))
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x01}), bytes.FromUint32(defaultDomainTypeApplicationMask))
}
