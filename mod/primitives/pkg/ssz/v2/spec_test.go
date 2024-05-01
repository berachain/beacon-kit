package ssz_test

import (
	"fmt"
	"os"
	"testing"

	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

const TestFileName = "fixtures/beacon_state_bellatrix.ssz" // https://goerli.beaconcha.in/slot/4744352
const debug = false

func debugPrint(s ...any) {
	if debug {
		fmt.Println(s...)
	}
}
func TestParityUint64(t *testing.T) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(t, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(t, err)

	object := sszState.LatestBlockHeader
	slot := object.Slot

	s := sszv2.NewSerializer()
	o2, err := s.MarshalSSZ(sszState.LatestBlockHeader.Slot)
	debugPrint("Local Serializer output:", o2, err)

	res := make([]byte, 0)
	res = ssz.MarshalUint64(res, slot)
	debugPrint("FastSSZ Output:", res)
	require.Equal(t, o2, res, "local output and fastssz output doesnt match")
}

func BenchmarkNativeUint64(b *testing.B) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	s := sszv2.NewSerializer()
	for i := 0; i < b.N; i++ {
		// Native impl
		o2, err := s.MarshalSSZ(sszState.LatestBlockHeader.Slot)
		debugPrint("Local Serializer output:", o2, err)
	}
}

func BenchmarkFastSSZUint64(b *testing.B) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		res := make([]byte, 0)
		res = ssz.MarshalUint64(res, sszState.LatestBlockHeader.Slot)
		debugPrint("FastSSZ Output:", res)
	}
}
