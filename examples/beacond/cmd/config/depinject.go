package config

import (
	"os"

	clientv2keyring "cosmossdk.io/client/v2/autocli/keyring"
	"cosmossdk.io/core/address"
	"cosmossdk.io/x/auth/tx"
	authtxconfig "cosmossdk.io/x/auth/tx/config"
	"cosmossdk.io/x/auth/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/runtime"
)

func ProvideClientContext(
	appCodec codec.Codec,
	interfaceRegistry codectypes.InterfaceRegistry,
	txConfigOpts tx.ConfigOptions,
	legacyAmino *codec.LegacyAmino,
	addressCodec address.Codec,
	validatorAddressCodec runtime.ValidatorAddressCodec,
	consensusAddressCodec runtime.ConsensusAddressCodec,
) client.Context {
	var err error

	clientCtx := client.Context{}.
		WithCodec(appCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithLegacyAmino(legacyAmino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithAddressCodec(addressCodec).
		WithValidatorAddressCodec(validatorAddressCodec).
		WithConsensusAddressCodec(consensusAddressCodec).
		WithHomeDir(".beacond").
		WithViper("") // uses by default the binary name as prefix

	// Read the config to overwrite the default values with the values from the
	// config file
	customClientTemplate, customClientConfig := initClientConfig()
	clientCtx, err = config.ReadDefaultValuesFromDefaultClientConfig(
		clientCtx,
		customClientTemplate,
		customClientConfig,
	)
	if err != nil {
		panic(err)
	}

	// textual is enabled by default, we need to re-create the tx config grpc
	// instead of bank keeper.
	txConfigOpts.TextualCoinMetadataQueryFn = authtxconfig.NewGRPCCoinMetadataQueryFn(
		clientCtx,
	)
	txConfig, err := tx.NewTxConfigWithOptions(clientCtx.Codec, txConfigOpts)
	if err != nil {
		panic(err)
	}
	clientCtx = clientCtx.WithTxConfig(txConfig)

	return clientCtx
}

func ProvideKeyring(
	clientCtx client.Context,
	addressCodec address.Codec,
) (clientv2keyring.Keyring, error) {
	kb, err := client.NewKeyringFromBackend(
		clientCtx,
		clientCtx.Keyring.Backend(),
	)
	if err != nil {
		return nil, err
	}

	return keyring.NewAutoCLIKeyring(kb)
}
