package spec_test

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/stretchr/testify/require"
)

func TestDomainTypeConversion(t *testing.T) {
	cs := spec.MainnetChainSpecData()
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x00}), cs.DomainTypeProposer)
	require.Equal(t, bytes.B4([]byte{0x01, 0x00, 0x00, 0x00}), cs.DomainTypeAttester)
	require.Equal(t, bytes.B4([]byte{0x02, 0x00, 0x00, 0x00}), cs.DomainTypeRandao)
	require.Equal(t, bytes.B4([]byte{0x03, 0x00, 0x00, 0x00}), cs.DomainTypeDeposit)
	require.Equal(t, bytes.B4([]byte{0x04, 0x00, 0x00, 0x00}), cs.DomainTypeVoluntaryExit)
	require.Equal(t, bytes.B4([]byte{0x05, 0x00, 0x00, 0x00}), cs.DomainTypeSelectionProof)
	require.Equal(t, bytes.B4([]byte{0x06, 0x00, 0x00, 0x00}), cs.DomainTypeAggregateAndProof)
	require.Equal(t, bytes.B4([]byte{0x00, 0x00, 0x00, 0x01}), cs.DomainTypeApplicationMask)
}
