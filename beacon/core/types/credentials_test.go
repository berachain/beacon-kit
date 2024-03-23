package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/stretchr/testify/require"
)

func TestNewCredentialsFromExecutionAddress(t *testing.T) {
	address := primitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	expectedCredentials := types.DepositCredentials{}
	expectedCredentials[0] = 0x01 // EthSecp256k1CredentialPrefix
	copy(expectedCredentials[12:], address[:])
	for i := 1; i < 12; i++ {
		expectedCredentials[i] = 0x00
	}
	require.Len(t, expectedCredentials, 32, "Expected credentials to be 32 bytes long")
	require.Equal(t, byte(0x01), expectedCredentials[0], "Expected prefix to be 0x01")
	require.Equal(t, address, primitives.ExecutionAddress(expectedCredentials[12:]), "Expected address to be set correctly")
	credentials := types.NewCredentialsFromExecutionAddress(address)
	require.Equal(t, expectedCredentials, credentials, "Generated credentials do not match expected")
}

func TestToExecutionAddress(t *testing.T) {
	expectedAddress := primitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	credentials := types.DepositCredentials{}
	for i := range credentials {
		// First byte should be 0x01
		if i == 0 {
			credentials[i] = 0x01 // EthSecp256k1CredentialPrefix

		} else if i > 0 && i < 12 {
			// then we have 11 bytes of padding
			credentials[i] = 0x00

		} else {
			// then the address
			credentials[i] = expectedAddress[i-12]
		}
	}

	address, err := credentials.ToExecutionAddress()
	require.NoError(t, err, "Conversion to execution address should not error")
	require.Equal(t, expectedAddress, address, "Converted address does not match expected")
}

func TestToExecutionAddress_InvalidPrefix(t *testing.T) {
	credentials := types.DepositCredentials{}
	for i := range credentials {
		if i == 0 {
			credentials[i] = 0x00 // Invalid prefix
		} else {
			credentials[i] = 0x00 // Padding or unused
		}
	}

	_, err := credentials.ToExecutionAddress()

	require.Error(t, err, "Expected an error due to invalid prefix")
}
