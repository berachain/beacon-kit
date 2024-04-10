package genesis_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/berachain/beacon-kit/mod/node-builder/commands/genesis"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const (
	// testGenesisFile is the path to the test genesis file.
	testGenesisFile = "./beacond/eth-genesis.json"
)

func TestAddExecutionPayloadCmd(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(cmd *cobra.Command) error
		testFunc      func(t *testing.T, cmd *cobra.Command)
		expectedError bool
	}{
		{
			name: "AddExecutionPayloadCmd",
			testFunc: func(t *testing.T, cmd *cobra.Command) {
				require.NotNil(t, cmd)
			},
		},
		{
			name: "RunE",
			setupFunc: func(cmd *cobra.Command) error {
				fmt.Println("getTestGenesisFilePath(t): ", getTestGenesisFilePath(t))
				return cmd.RunE(cmd, []string{getTestGenesisFilePath(t)})
			},
			testFunc: func(t *testing.T, cmd *cobra.Command) {
				// Check the output
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := genesis.AddExecutionPayloadCmd()
			if tt.setupFunc != nil {
				err := tt.setupFunc(cmd)
				if tt.expectedError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
			}
			tt.testFunc(t, cmd)
		})
	}
}

func getTestGenesisFilePath(t *testing.T) string {
	_, file, _, ok := runtime.Caller(2)
	require.True(t, ok)
	dir := filepath.Dir(file)
	return filepath.Join(dir, "../../../../", testGenesisFile)
}
