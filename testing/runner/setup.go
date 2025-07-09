package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	beaconkitgenesiscli "github.com/berachain/beacon-kit/cli/commands/genesis"
	beaconkitgenesisclitypes "github.com/berachain/beacon-kit/cli/commands/genesis/types"
	"github.com/berachain/beacon-kit/config/spec"
	beaconkitconsensustypes "github.com/berachain/beacon-kit/consensus-types/types"
	beaconkitconsensus "github.com/berachain/beacon-kit/consensus/cometbft/service"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/cometbft/cometbft/test/e2e/pkg/infra"
	"github.com/cometbft/cometbft/test/e2e/pkg/infra/docker"
	"github.com/cometbft/cometbft/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	AppAddressTCP  = "tcp://127.0.0.1:30000"
	AppAddressUNIX = "unix:///var/run/app.sock"

	PrivvalAddressTCP     = "tcp://0.0.0.0:27559"
	PrivvalAddressUNIX    = "unix:///var/run/privval.sock"
	PrivvalKeyFile        = "config/priv_validator_key.json"
	PrivvalStateFile      = "data/priv_validator_state.json"
	PrivvalDummyKeyFile   = "config/dummy_validator_key.json"
	PrivvalDummyStateFile = "data/dummy_validator_state.json"

	PrometheusConfigFile = "monitoring/prometheus.yml"

	depositAmount = 32000000000

	// Todo: parameterize this in the manifest
	elGenesisFilePath = "testing/files/eth-genesis.json"
)

// pregeneratedEthAddresses returns an address defined in testing/files/eth-genesis.json
func pregeneratedEthAddresses(i uint64) string {
	addresses := []string{
		"0x0474f52d25529c4db5f4E72F43303dA71B3541C6",
		"0x0e10cDAd84D788843aF48673C5b260A02ef78742",
		"0x0fb648Cb08e21602AF61AF53fE104E29d46433F7",
		"0x10FdFa4EFc83d6CC42F5ef14c13da8b98E458214",
		"0x12De044207a90709Ef2602D3D9D945d64dAe6147",
		"0x14DA5251a1EB236238969575ccE943e2Fb0f4AA1",
		"0x185F4Eebd01614aE3d12a5E49b184B054C46d37B",
		"0x187bE38A1f448b0F42423151A683dCAea949008B",
		"0x1CF7e940A657eE706718CF180eb21864DE9672C3",
		"0x1a0A57e5e6a66aD732295ddAF0aed286a4e64310",
		"0x1a0c826048DF0E4661E3c53bBd447d497E3f701F",
		"0x1e2e53c2451d0f9ED4B7952991BE0c95165D5c01",
		"0x1f1D0FCa7e19b799c315d4fDf31bA50e6A2AB153",
		"0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4",
		"0x25fc16D8E2314B305dF05C032E617638284801D6",
		"0x28879749Dda99387bdB43295B28bdF251d999F3b",
		"0x2B9935698dc5c19Ab7414AE22f27Da5F4478008a",
		"0x2E5f031578e8FF82199aaF16f42c44D43Fe61819",
		"0x2F4fD8a82A1400E654eeEC59b0e588445ffE0F96",
		"0x2d88ECD4d8F4b0A954886eE8C0802aE14684cd07",
		"0x2f6eB3D9a41157322dE01A6E707F6F118Cb00A7b",
		"0x3124d9885b11B52c56A2aee610AfCf5740d484F0",
		"0x3649839562C8dA64E6215EB0f5371629Ead9729D",
		"0x3DFb4173ec41EB976260fd689E5AB9772C66beaf",
		"0x3bd0E8f1B1E8Ec99a4E1762F4058F9884C93af31",
		"0x3f51B3BB6A18141282Ba002F7709c7E2f337F961",
		"0x4245537d9e3fb36fBBf054247FfFB28b0d931503",
		"0x440C37b22e8D7469128Ea7De6ac2f31419B4A8b1",
		"0x44a5FBfa7d6f3Fd92cca01f6764509f8Fc33dfa5",
		"0x47575DAE85403cD408d4639068D1187C427B9897",
		"0x49cE37B2019bb2d0B8b6a094ef87a6Dd625454A0",
		"0x4Afe0DFDAcc91F0fA2AEe39F9eAd66b64d03EbD6",
		"0x4bD04ABA9fc709835b1EE4789195d10E9e8E53F5",
		"0x4dC3aC871b22F8a98197B0aae976a8dE08e5Bebe",
		"0x4e59b44847b379578588920cA78FbF26c0B4956C",
		"0x5145b1B855bca67A119CB02A42aF4Bdbc66B725C",
		"0x51e15e71c865FE702C9347610667f83658A20e00",
		"0x5227aaebCA3E5e893547A667666E2e4e12Ca20e0",
		"0x54e1F990Dc0B7367F1E8eD96dA63BC4bca0E8061",
		"0x56898d1aFb10cad584961eb96AcD476C6826e41E",
		"0x5DD7bc3BEE395831ce499315ecAFE81DE0556F99",
		"0x611a42A2EF62c2461D123e3F0B64b93938bc4781",
		"0x62cB9bF32EA104f6D5eBf6879e876439f9492E4B",
		"0x67c942Ef50Fc690eA779067a6A0d444a8234baB5",
		"0x6CBcF4198fDA91D00fD469340E6DF6df086159e3",
		"0x6F69542fC88fF84C480FFf510aB7108120447247",
		"0x6a354C708fd248FD778F6adF75E41AA554700F68",
		"0x719Be866A77CeEc1BaC4FD37910c0975eFd52f55",
		"0x7469CeEf99FB67e4990c5F1c085a1B39b2902331",
		"0x7689BE67b205EB5d32811d95D60587Eae4F3036F",
		"0x795B761Db5969B7ba53472d5D37c230C859a472F",
		"0x7c4d7dB81c544B768E1f4782011077202B74B5C0",
		"0x7d7f187C2A05cDDCF700dCF2E02c96E7eF03f9B0",
		"0x7f0E54bc3C1a72405646F5dFbBE0D4565c649fe2",
		"0x800830F031ab1dd5895a5ec5B561427AD18f9ea8",
		"0x868a33C94F91398B6245e1f0E4CF128B2F28714B",
		"0x8724C57fb8f38A1FccA7177543dd1D8FcD49E5aa",
		"0x8a88215ae882dfA519730c40109556c1C235729f",
		"0x8b1e58f651CacaAa40291d2a6E0a6404d7Ed99e6",
		"0x92B3feac5b7816Dcef96a303c1D5112271A70D2c",
		"0x9C75eD1A37ae420b4FC0a1F4c26B673227Fd3AFa",
		"0x9beFa0FB7a1A9E6cC7596204DbB8962E87091D64",
		"0xA1d283f1a11A36D20FF38F29e12CA8F7Cf8709c1",
		"0xA6177defF3b768b1D678EdF7583b8cf210C777c0",
		"0xAC3c80F41C3049A89Aba8072FFbFc38a90fb6D8c",
		"0xAf325Ccc92ae883DEF1634D499d8B093192D7a0c",
		"0xB8865B4B8C56861534CC07ebBD2EA569a9a16323",
		"0xBC3c03b4185A6F10618CC4E7B9f4AdD59AB5FbbA",
		"0xBC9BC89b295a14F3976234Cc37C73e3D286f3a49",
		"0xC4DD08191B4d5173e3698491A11e05b63F9Ee097",
		"0xC4eD09A472B82516daa3A4d8D1E38AE94CF4855C",
		"0xC59D8935c0570E75BA0E55E3C661f535C86e368B",
		"0xD073a84e2ccDF91a9025179330438485E886D206",
		"0xD2a3b89AE8D2c3bD39E2F24612ecFCD8600360C9",
		"0xD3c5dAC705289cD005C402C79C8445a47502d8be",
		"0xD6D4Fb22B91FAa54700852a05698B37d45514166",
		"0xDE8E0E641E2Fb52c22460e6a1533c6BD13A00B37",
		"0xDc6De65f6070b409125217a12Cf576A208Cc1998",
		"0xDe5C7198e2416baB7e7a1EA758858Cd7301740bF",
		"0xE3d2b9191EaBD3636A5dd057D522335cfae8c7CF",
		"0xE5981AA0807eb05611cDb666e32e53b2001bd61d",
		"0xE69ac59e1DF47291AaB8DEc540C796f81De7c892",
		"0xE7F444b5f772281384117674002d540131e533Ca",
		"0xF60fD8632Fc77E19b3A0637d115d0fdd06F36968",
		"0xF99139D2FCc5E25F57B0B91fd382a21B3AFF9cbA",
		"0xF9f58a87C3f0B3A4a0592938c80C41a7c659f855",
		"0xFeb1eafa0154D291e28e393FAF10Bc89e5cCbB22",
		"0xaEf63D7F7e2637c99FeA1B63366b244B4da12D70",
		"0xb86d37333072eFb48cEaa46C67271A27CA5Bda82",
		"0xb87fb371Bd3C2093b608cd0E7a8dDD60Bb05C995",
		"0xbE651bc261b9Da5499a24Bf4214fD494c6e1F5Ac",
		"0xbcC90AD39D377cA0b7b4F36eC463103E2728C33F",
		"0xcB6632daA65e6c921c2963C37320f63f54fC8fE3",
		"0xd0F043dED28773953562f824334C4cbb84210AE7",
		"0xdBfb742BD2e0e6E353cb61E75B9e11257aC8fB1A",
		"0xdb96E9cDD1e457b602f97d33e51736D7a5216496",
		"0xdb9cB94B166DfdC9F337EA63b32B448d993d7008",
		"0xe3024d098953661638d59E06f7FcD0B61c424854",
		"0xea94749deFcc40dC5992687974b1C84B1bB9D6df",
		"0xf11D16e2EE6BefED82Fbca0b005906E09303aB95",
		"0xf22FbA9cBeB75ED353931418E9eca71EF1Ab9921",
		"0xf4b2eb959A4C4b0E148340676999FC0446D446D4",
		"0xf6B6A52aA9BD788837c6682f47ACE009BD84b6fc",
		"0xf97a36c417D33D1fC60a9163A8715e1aecb29102",
	}
	if i >= uint64(len(addresses)) {
		panic(fmt.Sprintf("invalid eth address index %d", i))
	}
	return addresses[i]
}

// Setup sets up the testnet configuration.
func Setup(testnet *e2e.Testnet, infp infra.Provider) error {
	logger.Info("setup", "msg", log.NewLazySprintf("Generating testnet files in %#q", testnet.Dir))

	if err := os.MkdirAll(testnet.Dir, os.ModePerm); err != nil {
		return err
	}
	gethDir, err := filepath.Abs(filepath.Join(testnet.Dir, "geth"))
	if err != nil {
		return fmt.Errorf("error getting absolute path for geth directory: %w", err)
	}
	if err := os.MkdirAll(gethDir, os.ModePerm); err != nil {
		return err
	}

	genesis, err := MakeGenesis(testnet)
	if err != nil {
		return err
	}

	// Create eth-genesis.json and set deposit storage.
	ethGenesisBz, err := MakeEthGenesis(genesis, elGenesisFilePath)
	if err != nil {
		return err
	}
	// EthGenesis.SaveAs
	err = os.WriteFile(filepath.Join(gethDir, "eth-genesis.json"), ethGenesisBz, 0o644)
	if err != nil {
		return err
	}

	// Set execution payload
	genesis, err = UpdateGenesis(genesis, ethGenesisBz)
	if err != nil {
		return err
	}

	for _, node := range testnet.Nodes {
		nodeDir := filepath.Join(testnet.Dir, node.Name)

		dirs := []string{
			filepath.Join(nodeDir, "config"),
			filepath.Join(nodeDir, "data"),
			filepath.Join(nodeDir, "data", "app"),
		}
		for _, dir := range dirs {
			// light clients don't need an app directory
			if node.Mode == e2e.ModeLight && strings.Contains(dir, "app") {
				continue
			}
			err := os.MkdirAll(dir, 0o755)
			if err != nil {
				return err
			}
		}

		cfg, err := MakeConfig(node)
		if err != nil {
			return err
		}
		config.WriteConfigFile(filepath.Join(nodeDir, "config", "config.toml"), cfg) // panics

		appCfg, err := MakeAppConfig(node)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(nodeDir, "config", "app.toml"), appCfg, 0o644) //nolint:gosec
		if err != nil {
			return err
		}

		if node.Mode == e2e.ModeLight {
			// stop early if a light client
			continue
		}

		err = genesis.SaveAs(filepath.Join(nodeDir, "config", "genesis.json"))
		if err != nil {
			return err
		}

		err = (&p2p.NodeKey{PrivKey: node.NodeKey}).SaveAs(filepath.Join(nodeDir, "config", "node_key.json"))
		if err != nil {
			return err
		}

		(privval.NewFilePV(node.PrivvalKey,
			filepath.Join(nodeDir, PrivvalKeyFile),
			filepath.Join(nodeDir, PrivvalStateFile),
		)).Save()

		// Set up a dummy validator. CometBFT requires a file PV even when not used, so we
		// give it a dummy such that it will fail if it actually tries to use it.
		(privval.NewFilePV(ed25519.GenPrivKey(),
			filepath.Join(nodeDir, PrivvalDummyKeyFile),
			filepath.Join(nodeDir, PrivvalDummyStateFile),
		)).Save()

		if testnet.LatencyEmulationEnabled {
			// Generate a shell script file containing tc (traffic control) commands
			// to emulate latency to other nodes.
			tcCmds, err := tcCommands(node, infp)
			if err != nil {
				return err
			}
			latencyPath := filepath.Join(nodeDir, "emulate-latency.sh")
			//nolint: gosec // G306: Expect WriteFile permissions to be 0600 or less
			if err = os.WriteFile(latencyPath, []byte(strings.Join(tcCmds, "\n")), 0o755); err != nil {
				return err
			}
		}
	}

	if testnet.Prometheus {
		if err := WritePrometheusConfig(testnet, PrometheusConfigFile); err != nil {
			return err
		}
		// Make a copy of the Prometheus config file in the testnet directory.
		// This should be temporary to keep it compatible with the qa-infra
		// repository.
		if err := WritePrometheusConfig(testnet, filepath.Join(testnet.Dir, "prometheus.yml")); err != nil {
			return err
		}
	}

	//nolint: revive
	if err := infp.Setup(); err != nil {
		return err
	}

	/////
	// Adding geth node to docker configuration
	/////
	// Todo: Change the template in CometBFT test/e2e/pkg/infra/docker.go and add the geth node statically

	// Todo: Add geth node to digitalocean infrastructure provider in CometBFT.
	if infp.GetInfrastructureData().Provider != "docker" {
		return errors.New("provider must be docker for now")
	}

	// copy jwt.hex
	jwtHexSourcePath := "testing/files/jwt.hex" // Todo: do not hard code the path
	jwtHexPath := filepath.Join(gethDir, "jwt.hex")
	jwt, err := os.ReadFile(jwtHexSourcePath)
	if err != nil {
		return fmt.Errorf("error reading jwt hex file %s: %w", jwtHexSourcePath, err)
	}
	err = os.WriteFile(jwtHexPath, jwt, 0o644)
	if err != nil {
		return fmt.Errorf("error creating jwt hex file %s: %w", jwtHexPath, err)
	}

	// geth init
	err = docker.Exec(context.Background(), "run", "--rm", "-v", fmt.Sprintf("%s:/.tmp", gethDir),
		"ethereum/client-go", "init", "--datadir", "/.tmp", "/.tmp/eth-genesis.json")
	if err != nil {
		return fmt.Errorf("error during geth init: %w", err)
	}

	// geth service compose
	path := filepath.Join(testnet.Dir, "compose.yaml")
	compose, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading compose.yaml file %s: %w", path, err)
	}
	gethService := fmt.Sprintf(`services:
  geth:
    image: ethereum/client-go
    labels:
      e2e: true
    ports:
    - 30303:30303
    - 8545:8545
    - 8551:8551
    command: >
      --syncmode=full
      --http
      --http.addr 0.0.0.0
      --http.api eth,net
      --http.vhosts "*"
      --authrpc.addr 0.0.0.0
      --authrpc.jwtsecret /.tmp/jwt.hex
      --authrpc.vhosts "*"
      --datadir /.tmp
      --ipcpath /.tmp/eth-engine.ipc
    volumes:
    - %s:/.tmp
  load:
    image: polycli:latest
    depends_on:
      - geth
    #restart: on-failure:1
    labels:
      e2e: true
    command: >
      loadtest
      --verbosity 700
      --rpc-url http://geth:8545
      --mode random
      --concurrency 1
      --requests 600
      --rate-limit 100
      --summarize
`, gethDir)
	/* When polycli v0.1.80 is released, these additional parameters should be added:
	   --sending-address-count 10
	   --pre-fund-sending-addresses
	*/
	updated := bytes.Replace(compose, []byte("services:\n"), []byte(gethService), 1)
	err = os.WriteFile(path, updated, 0o644)
	if err != nil {
		return fmt.Errorf("error writing compose.yaml file %s: %w", path, err)
	}
	return nil
}

// MakeGenesis generates a genesis document.
func MakeGenesis(testnet *e2e.Testnet) (types.GenesisDoc, error) {
	genesis := types.GenesisDoc{
		GenesisTime:     time.Now(),
		ChainID:         testnet.Name,
		ConsensusParams: beaconkitconsensus.DefaultConsensusParams(crypto.CometBLSType),
		InitialHeight:   testnet.InitialHeight,
	}
	// set the app version to 1
	genesis.ConsensusParams.Version.App = 1
	genesis.ConsensusParams.Evidence.MaxAgeNumBlocks = e2e.EvidenceAgeHeight
	genesis.ConsensusParams.Evidence.MaxAgeDuration = e2e.EvidenceAgeTime
	if testnet.BlockMaxBytes != 0 {
		genesis.ConsensusParams.Block.MaxBytes = testnet.BlockMaxBytes
	}
	if testnet.VoteExtensionsUpdateHeight == -1 {
		genesis.ConsensusParams.Feature.VoteExtensionsEnableHeight = testnet.VoteExtensionsEnableHeight
	}
	if testnet.PbtsUpdateHeight == -1 {
		genesis.ConsensusParams.Feature.PbtsEnableHeight = testnet.PbtsEnableHeight
	}
	if len(testnet.InitialState) > 0 { //nolint:nestif
		appState, err := json.Marshal(testnet.InitialState)
		if err != nil {
			return genesis, err
		}
		genesis.AppState = appState
	} else {
		chainSpec, err := spec.DevnetChainSpec()
		if err != nil {
			return genesis, err
		}

		appState := beaconkitconsensustypes.DefaultGenesis(chainSpec.GenesisForkVersion())
		amount := math.Gwei(depositAmount)
		forkData := beaconkitconsensustypes.NewForkData(chainSpec.GenesisForkVersion(), common.Root{})
		domainType := chainSpec.DomainTypeDeposit()

		var i uint64 = 0
		for validator := range testnet.Validators {
			executionAddress := common.NewExecutionAddressFromHex(pregeneratedEthAddresses(i))
			credentials := beaconkitconsensustypes.NewCredentialsFromExecutionAddress(executionAddress)
			blsSigner := signer.BLSSigner{PrivValidator: types.MockPV{PrivKey: validator.PrivvalKey}}

			var signature crypto.BLSSignature
			_, signature, err = beaconkitconsensustypes.CreateAndSignDepositMessage(
				forkData,
				domainType,
				blsSigner,
				credentials,
				amount)
			if err != nil {
				return genesis, err
			}

			deposit := &beaconkitconsensustypes.Deposit{
				Pubkey:      blsSigner.PublicKey(), // compressed public key
				Amount:      amount,
				Signature:   signature,
				Credentials: credentials,
				Index:       i,
			}

			appState.Deposits = append(appState.Deposits, deposit)
			i++
		}
		gen := make(map[string]json.RawMessage)
		gen["beacon"], err = json.Marshal(appState)
		if err != nil {
			return genesis, err
		}
		appStateJson, err := json.Marshal(gen)
		if err != nil {
			return genesis, err
		}
		genesis.AppState = appStateJson
	}

	// Customized genesis fields provided in the manifest
	if len(testnet.Genesis) > 0 {
		v := viper.New()
		v.SetConfigType("json")

		for _, field := range testnet.Genesis {
			key, value, err := e2e.ParseKeyValueField("genesis", field)
			if err != nil {
				return genesis, err
			}
			logger.Debug("Applying 'genesis' field", key, value)
			v.Set(key, value)
		}

		// We use viper because it leaves untouched keys that are not set.
		// The GenesisDoc does not use the original `mapstructure` tag.
		err := v.Unmarshal(&genesis, func(d *mapstructure.DecoderConfig) {
			d.TagName = "json"
			d.ErrorUnused = true
		})
		if err != nil {
			return genesis, fmt.Errorf("failed parsing 'genesis' field: %v", err)
		}
	}

	return genesis, genesis.ValidateAndComplete()
}

// MakeEthGenesis generates an eth-genesis document for a geth node.
func MakeEthGenesis(genesis types.GenesisDoc, elGenesisPath string) (json.RawMessage, error) {
	bz, _ := json.RawMessage{}.MarshalJSON()

	// Read existing el genesis file
	elGenesisBz, err := afero.ReadFile(afero.NewOsFs(), elGenesisPath)
	if err != nil {
		return bz, fmt.Errorf("error reading eth-genesis at %s", elGenesisPath)
	}

	// Read app state from AppGenesis
	var gen map[string]json.RawMessage
	err = json.Unmarshal(genesis.AppState, &gen)
	if err != nil {
		return bz, errors.New("error reading app state from genesis")
	}
	appStateJson, exists := gen["beacon"]
	if !exists {
		return bz, errors.New("beacon key not found in genesis app state")
	}

	// Get deposits from app state
	var beaconState struct {
		Deposits beaconkitconsensustypes.Deposits `json:"deposits"`
	}
	err = json.Unmarshal(appStateJson, &beaconState)
	if err != nil {
		return bz, errors.New("deposits not found in genesis app state")
	}
	deposits := beaconState.Deposits

	// Set the storage of the deposit contract with deposits count and root.
	count := big.NewInt(int64(len(deposits)))
	root := deposits.HashTreeRoot()

	// Get allocs key and unmarshal eth-genesis file.
	allocsKey := beaconkitgenesisclitypes.DefaultAllocsKey
	elGenesis := &beaconkitgenesisclitypes.DefaultEthGenesisJSON{}
	if err = json.Unmarshal(elGenesisBz, elGenesis); err != nil {
		return bz, errors.New("error unmarshalling eth-genesis")
	}

	// Generate deposit storage
	chainSpec, err := spec.DevnetChainSpec()
	if err != nil {
		return bz, err
	}
	depositAddr := ethcommon.Address(chainSpec.DepositContractAddress())
	updatedAllocs := beaconkitgenesiscli.WriteDepositStorage(elGenesis, depositAddr, count, root)

	// Unmarshal eth-genesis for update
	var existingGenesis map[string]interface{}
	if err = json.Unmarshal(elGenesisBz, &existingGenesis); err != nil {
		return bz, err
	}

	// Get existing alloc.
	existingAllocs, ok := existingGenesis[allocsKey].(map[string]interface{})
	if !ok {
		return bz, errors.New("invalid alloc format in genesis file")
	}

	// Update only the deposit contract entry
	if account, exists := updatedAllocs[depositAddr]; exists {
		existingAllocs[depositAddr.Hex()] = account
	} else {
		return bz, errors.New("updated allocs could not be set")
	}

	// Marshal updated eth-genesis
	bz, err = json.MarshalIndent(existingGenesis, "", "  ")
	if err != nil {
		return bz, err
	}

	return bz, nil
}

// UpdateGenesis updates the application genesis execution payload.
func UpdateGenesis(genesis types.GenesisDoc, ethGenesisBz json.RawMessage) (types.GenesisDoc, error) {
	// Unmarshal the eth-genesis file.
	ethGenesis := &gethprimitives.Genesis{}
	if err := ethGenesis.UnmarshalJSON(ethGenesisBz); err != nil {
		return types.GenesisDoc{}, fmt.Errorf("failed to unmarshal eth1 genesis %v", err)
	}
	genesisBlock := ethGenesis.ToBlock()

	// Create the execution payload.
	payload := gethprimitives.BlockToExecutableData(
		genesisBlock,
		nil,
		nil,
		nil,
	).ExecutionPayload

	// Inject the execution payload.
	chainSpec, err := spec.DevnetChainSpec()
	if err != nil {
		return types.GenesisDoc{}, err
	}
	eph, err := beaconkitgenesiscli.ExecutableDataToExecutionPayloadHeader(
		chainSpec.GenesisForkVersion(),
		payload,
		chainSpec.MaxWithdrawalsPerPayload(),
	)
	if err != nil {
		return types.GenesisDoc{}, fmt.Errorf("failed to convert executable data to execution payload header: %v", err)
	}
	if eph == nil {
		return types.GenesisDoc{}, errors.New("failed to get execution payload header")
	}

	// Unmarshal app genesis
	var appState map[string]json.RawMessage
	err = json.Unmarshal(genesis.AppState, &appState)
	if err != nil {
		return types.GenesisDoc{}, err
	}

	// Get beacon state
	genesisInfo := &beaconkitconsensustypes.Genesis{}
	if err = json.Unmarshal(appState["beacon"], genesisInfo); err != nil {
		return types.GenesisDoc{}, fmt.Errorf("failed to unmarshal beacon state: %v", err)
	}

	// Update execution payload header in app genesis
	genesisInfo.ExecutionPayloadHeader = eph

	// Marshal beacon state
	appState["beacon"], err = json.Marshal(genesisInfo)
	if err != nil {
		return types.GenesisDoc{}, fmt.Errorf("failed to marshal beacon state: %v", err)
	}

	// Marshal app genesis
	genesis.AppState, err = json.MarshalIndent(appState, "", "  ")
	if err != nil {
		return types.GenesisDoc{}, err
	}

	return genesis, nil
}

// MakeConfig generates a CometBFT config for a node.
func MakeConfig(node *e2e.Node) (*config.Config, error) {
	cfg := beaconkitconsensus.DefaultConfig()
	cfg.Moniker = node.Name
	cfg.ProxyApp = AppAddressTCP

	cfg.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	cfg.RPC.PprofListenAddress = ":6060"

	cfg.GRPC.ListenAddress = "tcp://0.0.0.0:26670"
	cfg.GRPC.VersionService.Enabled = true
	cfg.GRPC.BlockService.Enabled = true
	cfg.GRPC.BlockResultsService.Enabled = true

	cfg.P2P.ExternalAddress = fmt.Sprintf("tcp://%v", node.AddressP2P(false))
	cfg.P2P.AddrBookStrict = false

	cfg.DBBackend = node.Database
	cfg.StateSync.DiscoveryTime = 5 * time.Second
	cfg.BlockSync.Version = node.BlockSyncVersion
	cfg.Consensus.PeerGossipIntraloopSleepDuration = node.Testnet.PeerGossipIntraloopSleepDuration
	cfg.Mempool.ExperimentalMaxGossipConnectionsToNonPersistentPeers = int(node.Testnet.ExperimentalMaxGossipConnectionsToNonPersistentPeers)
	cfg.Mempool.ExperimentalMaxGossipConnectionsToPersistentPeers = int(node.Testnet.ExperimentalMaxGossipConnectionsToPersistentPeers)

	// Assume that full nodes and validators will have a data companion
	// attached, which will need access to the privileged gRPC endpoint.
	if (node.Mode == e2e.ModeValidator || node.Mode == e2e.ModeFull) && node.EnableCompanionPruning {
		cfg.Storage.Pruning.DataCompanion.Enabled = true
		cfg.Storage.Pruning.DataCompanion.InitialBlockRetainHeight = 0
		cfg.Storage.Pruning.DataCompanion.InitialBlockResultsRetainHeight = 0
		cfg.GRPC.Privileged.ListenAddress = "tcp://0.0.0.0:26671"
		cfg.GRPC.Privileged.PruningService.Enabled = true
	}

	switch node.ABCIProtocol {
	case e2e.ProtocolUNIX:
		cfg.ProxyApp = AppAddressUNIX
	case e2e.ProtocolTCP:
		cfg.ProxyApp = AppAddressTCP
	case e2e.ProtocolGRPC:
		cfg.ProxyApp = AppAddressTCP
		cfg.ABCI = "grpc"
	case e2e.ProtocolBuiltin, e2e.ProtocolBuiltinConnSync:
		cfg.ProxyApp = ""
		cfg.ABCI = ""
	default:
		return nil, fmt.Errorf("unexpected ABCI protocol setting %q", node.ABCIProtocol)
	}

	// CometBFT errors if it does not have a privval key set up, regardless of whether
	// it's actually needed (e.g. for remote KMS or non-validators). We set up a dummy
	// key here by default, and use the real key for actual validators that should use
	// the file privval.
	cfg.PrivValidatorListenAddr = ""
	cfg.PrivValidatorKey = PrivvalDummyKeyFile
	cfg.PrivValidatorState = PrivvalDummyStateFile

	switch node.Mode {
	case e2e.ModeValidator:
		switch node.PrivvalProtocol {
		case e2e.ProtocolFile:
			cfg.PrivValidatorKey = PrivvalKeyFile
			cfg.PrivValidatorState = PrivvalStateFile
		case e2e.ProtocolUNIX:
			cfg.PrivValidatorListenAddr = PrivvalAddressUNIX
		case e2e.ProtocolTCP:
			cfg.PrivValidatorListenAddr = PrivvalAddressTCP
		default:
			return nil, fmt.Errorf("invalid privval protocol setting %q", node.PrivvalProtocol)
		}
	case e2e.ModeSeed:
		cfg.P2P.SeedMode = true
		cfg.P2P.PexReactor = true
	case e2e.ModeFull, e2e.ModeLight:
		// Don't need to do anything, since we're using a dummy privval key by default.
	default:
		return nil, fmt.Errorf("unexpected mode %q", node.Mode)
	}

	if node.StateSync {
		cfg.StateSync.Enable = true
		cfg.StateSync.RPCServers = []string{}
		for _, peer := range node.Testnet.ArchiveNodes() {
			if peer.Name == node.Name {
				continue
			}
			cfg.StateSync.RPCServers = append(cfg.StateSync.RPCServers, peer.AddressRPC())
		}
		if len(cfg.StateSync.RPCServers) < 2 {
			return nil, errors.New("unable to find 2 suitable state sync RPC servers")
		}
	}

	cfg.P2P.Seeds = ""
	for _, seed := range node.Seeds {
		if len(cfg.P2P.Seeds) > 0 {
			cfg.P2P.Seeds += ","
		}
		cfg.P2P.Seeds += seed.AddressP2P(true)
	}
	cfg.P2P.PersistentPeers = ""
	for _, peer := range node.PersistentPeers {
		if len(cfg.P2P.PersistentPeers) > 0 {
			cfg.P2P.PersistentPeers += ","
		}
		cfg.P2P.PersistentPeers += peer.AddressP2P(true)
	}
	if node.Testnet.DisablePexReactor {
		cfg.P2P.PexReactor = false
	}

	if node.Testnet.LogLevel != "" {
		cfg.LogLevel = node.Testnet.LogLevel
	}

	if node.Testnet.LogFormat != "" {
		cfg.LogFormat = node.Testnet.LogFormat
	}

	if node.Prometheus {
		cfg.Instrumentation.Prometheus = true
	}

	if node.ExperimentalKeyLayout != "" {
		cfg.Storage.ExperimentalKeyLayout = node.ExperimentalKeyLayout
	}

	if node.Compact {
		cfg.Storage.Compact = node.Compact
	}

	if node.DiscardABCIResponses {
		cfg.Storage.DiscardABCIResponses = node.DiscardABCIResponses
	}

	if node.Indexer != "" {
		cfg.TxIndex.Indexer = node.Indexer
	}

	if node.CompactionInterval != 0 && node.Compact {
		cfg.Storage.CompactionInterval = node.CompactionInterval
	}

	// We currently need viper in order to parse config files.
	if len(node.Config) > 0 {
		v := viper.New()
		for _, field := range node.Config {
			key, value, err := e2e.ParseKeyValueField("config", field)
			if err != nil {
				return nil, err
			}
			logger.Debug("Applying 'config' field", "node", node.Name, key, value)
			v.Set(key, value)
		}
		err := v.Unmarshal(cfg, func(d *mapstructure.DecoderConfig) {
			d.ErrorUnused = true
		})
		if err != nil {
			return nil, fmt.Errorf("failed parsing 'config' field of node %v: %v", node.Name, err)
		}
	}

	return cfg, nil
}

// MakeAppConfig generates an ABCI application config for a node.
func MakeAppConfig(node *e2e.Node) ([]byte, error) {
	cfg := map[string]any{
		"chain_id":                      node.Testnet.Name,
		"dir":                           "data/app",
		"listen":                        AppAddressUNIX,
		"mode":                          node.Mode,
		"protocol":                      "socket",
		"persist_interval":              node.PersistInterval,
		"snapshot_interval":             node.SnapshotInterval,
		"retain_blocks":                 node.RetainBlocks,
		"key_type":                      node.PrivvalKey.Type(),
		"prepare_proposal_delay":        node.Testnet.PrepareProposalDelay,
		"process_proposal_delay":        node.Testnet.ProcessProposalDelay,
		"check_tx_delay":                node.Testnet.CheckTxDelay,
		"vote_extension_delay":          node.Testnet.VoteExtensionDelay,
		"finalize_block_delay":          node.Testnet.FinalizeBlockDelay,
		"vote_extension_size":           node.Testnet.VoteExtensionSize,
		"vote_extensions_enable_height": node.Testnet.VoteExtensionsEnableHeight,
		"vote_extensions_update_height": node.Testnet.VoteExtensionsUpdateHeight,
		"abci_requests_logging_enabled": node.Testnet.ABCITestsEnabled,
		"pbts_enable_height":            node.Testnet.PbtsEnableHeight,
		"pbts_update_height":            node.Testnet.PbtsUpdateHeight,
	}
	switch node.ABCIProtocol {
	case e2e.ProtocolUNIX:
		cfg["listen"] = AppAddressUNIX
	case e2e.ProtocolTCP:
		cfg["listen"] = AppAddressTCP
	case e2e.ProtocolGRPC:
		cfg["listen"] = AppAddressTCP
		cfg["protocol"] = "grpc"
	case e2e.ProtocolBuiltin, e2e.ProtocolBuiltinConnSync:
		delete(cfg, "listen")
		cfg["protocol"] = string(node.ABCIProtocol)
	default:
		return nil, fmt.Errorf("unexpected ABCI protocol setting %q", node.ABCIProtocol)
	}
	if node.Mode == e2e.ModeValidator {
		switch node.PrivvalProtocol {
		case e2e.ProtocolFile:
		case e2e.ProtocolTCP:
			cfg["privval_server"] = PrivvalAddressTCP
			cfg["privval_key"] = PrivvalKeyFile
			cfg["privval_state"] = PrivvalStateFile
		case e2e.ProtocolUNIX:
			cfg["privval_server"] = PrivvalAddressUNIX
			cfg["privval_key"] = PrivvalKeyFile
			cfg["privval_state"] = PrivvalStateFile
		default:
			return nil, fmt.Errorf("unexpected privval protocol setting %q", node.PrivvalProtocol)
		}
	}

	if len(node.Testnet.ValidatorUpdates) > 0 {
		validatorUpdates := map[string]map[string]int64{}
		for height, validators := range node.Testnet.ValidatorUpdates {
			updateVals := map[string]int64{}
			for node, power := range validators {
				updateVals[base64.StdEncoding.EncodeToString(node.PrivvalKey.PubKey().Bytes())] = power
			}
			validatorUpdates[strconv.FormatInt(height, 10)] = updateVals
		}
		cfg["validator_update"] = validatorUpdates
	}

	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate app config: %w", err)
	}
	return buf.Bytes(), nil
}

//go:embed templates/prometheus-yml.tmpl
var prometheusYamlTemplate string

func WritePrometheusConfig(testnet *e2e.Testnet, path string) error {
	tmpl, err := template.New("prometheus-yaml").Parse(prometheusYamlTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, testnet)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, buf.Bytes(), 0o644) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

// UpdateConfigStateSync updates the state sync config for a node.
func UpdateConfigStateSync(node *e2e.Node, height int64, hash []byte) error {
	cfgPath := filepath.Join(node.Testnet.Dir, node.Name, "config", "config.toml")

	// FIXME Apparently there's no function to simply load a config file without
	// involving the entire Viper apparatus, so we'll just resort to regexps.
	bz, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	bz = regexp.MustCompile(`(?m)^trust_height =.*`).ReplaceAll(bz, []byte(fmt.Sprintf(`trust_height = %v`, height)))
	bz = regexp.MustCompile(`(?m)^trust_hash =.*`).ReplaceAll(bz, []byte(fmt.Sprintf(`trust_hash = "%X"`, hash)))
	return os.WriteFile(cfgPath, bz, 0o644) //nolint:gosec
}
