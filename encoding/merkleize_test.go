package encoding

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/prysmaticlabs/go-bitfield"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/encoding/ssz"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestTransactionsRoot(t *testing.T) {
	tests := []struct {
		name    string
		txs     [][]byte
		want    [32]byte
		wantErr bool
	}{
		{
			name: "nil",
			txs:  nil,
			want: [32]byte{127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123, 250, 34, 222, 53, 209, 249, 190, 215, 171, 6, 29, 148, 1, 253, 71, 227, 74, 84, 251, 237, 225},
		},
		{
			name: "empty",
			txs:  [][]byte{},
			want: [32]byte{127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123, 250, 34, 222, 53, 209, 249, 190, 215, 171, 6, 29, 148, 1, 253, 71, 227, 74, 84, 251, 237, 225},
		},
		{
			name: "one tx",
			txs:  [][]byte{{1, 2, 3}},
			want: [32]byte{102, 209, 140, 87, 217, 28, 68, 12, 133, 42, 77, 136, 191, 18, 234, 105, 166, 228, 216, 235, 230, 95, 200, 73, 85, 33, 134, 254, 219, 97, 82, 209},
		},
		{
			name: "max txs",
			txs: func() [][]byte {
				var txs [][]byte
				for i := 0; i < fieldparams.MaxTxsPerPayloadLength; i++ {
					txs = append(txs, []byte{})
				}
				return txs
			}(),
			want: [32]byte{13, 66, 254, 206, 203, 58, 48, 133, 78, 218, 48, 231, 120, 90, 38, 72, 73, 137, 86, 9, 31, 213, 185, 101, 103, 144, 0, 236, 225, 57, 47, 244},
		},
		{
			name: "exceed max txs",
			txs: func() [][]byte {
				var txs [][]byte
				for i := 0; i < fieldparams.MaxTxsPerPayloadLength+1; i++ {
					txs = append(txs, []byte{})
				}
				return txs
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransactionsRoot(tt.txs)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionsRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ssss, _ := ssz.TransactionsRoot(tt.txs)
			fmt.Println("DIFFFF", got, ssss)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionsRoot() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MerkleizeVectorSSZ(t *testing.T) {
	t.Run("empty vector", func(t *testing.T) {
		attList := make([]*ethpb.Attestation, 0)
		expected := [32]byte{83, 109, 152, 131, 127, 45, 209, 101, 165, 93, 94, 234, 233, 20, 133, 149, 68, 114, 213, 111, 36, 109, 242, 86, 191, 60, 174, 25, 53, 42, 18, 60}
		length := uint64(16)
		root, err := MerkleizeVectorSSZ(attList, length)
		require.NoError(t, err)
		require.Equal(t, expected, root)
	})
	t.Run("non empty vector", func(t *testing.T) {
		sig := make([]byte, 96)
		br := make([]byte, 32)
		attList := make([]*ethpb.Attestation, 1)
		attList[0] = &ethpb.Attestation{
			AggregationBits: bitfield.Bitlist{0x01},
			Data: &ethpb.AttestationData{
				BeaconBlockRoot: br,
				Source: &ethpb.Checkpoint{
					Root: br,
				},
				Target: &ethpb.Checkpoint{
					Root: br,
				},
			},
			Signature: sig,
		}
		expected := [32]byte{199, 186, 55, 142, 200, 75, 219, 191, 66, 153, 100, 181, 200, 15, 143, 160, 25, 133, 105, 26, 183, 107, 10, 198, 232, 231, 107, 162, 243, 243, 56, 20}
		length := uint64(16)
		root, err := MerkleizeVectorSSZ(attList, length)
		require.NoError(t, err)
		require.Equal(t, expected, root)
	})
}

// func Test_MerkleizeListSSZ(t *testing.T) {
// 	t.Run("empty vector", func(t *testing.T) {
// 		attList := make([]*ethpb.Attestation, 0)
// 		expected := [32]byte{121, 41, 48, 187, 213, 186, 172, 67, 188, 199, 152, 238, 73, 170, 129, 133, 239, 118, 187, 59, 68, 186, 98, 185, 29, 134, 174, 86, 158, 75, 181, 53}
// 		length := uint64(16)
// 		root, err := MerkleizeListSSZ(attList, length)
// 		require.NoError(t, err)
// 		require.Equal(t, expected, root)
// 	})
// 	t.Run("non empty vector", func(t *testing.T) {
// 		sig := make([]byte, 96)
// 		br := make([]byte, 32)
// 		attList := make([]*ethpb.Attestation, 1)
// 		attList[0] = &ethpb.Attestation{
// 			AggregationBits: bitfield.Bitlist{0x01},
// 			Data: &ethpb.AttestationData{
// 				BeaconBlockRoot: br,
// 				Source: &ethpb.Checkpoint{
// 					Root: br,
// 				},
// 				Target: &ethpb.Checkpoint{
// 					Root: br,
// 				},
// 			},
// 			Signature: sig,
// 		}
// 		expected := [32]byte{161, 247, 30, 234, 219, 222, 154, 88, 7, 207, 6, 23, 46, 125, 135, 67, 225, 178, 217, 131, 113, 124, 242, 106, 194, 43, 205, 194, 49, 172, 232, 229}
// 		length := uint64(16)
// 		root, err := MerkleizeListSSZ(attList, length)
// 		require.NoError(t, err)
// 		require.Equal(t, expected, root)
// 	})
// }
