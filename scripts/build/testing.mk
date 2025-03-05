#!/usr/bin/make -f

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

# ask_reset_dir_func checks if the directory passed in exists, and if so asks the user whether it
# should delete it. Note that on linux, docker may have created the directory with root
# permissions, so we may need to ask the user to delete it with sudo
define ask_reset_dir_func
	@abs_path=$(abspath $(1)); \
	if test -d "$$abs_path"; then \
		read -p "Directory '$$abs_path' exists. Do you want to delete it? (y/n): " confirm && \
		if [ "$$confirm" = "y" ]; then \
			echo "Deleting directory '$$abs_path'..."; \
			rm -rf "$$abs_path" 2>/dev/null || sudo rm -rf "$$abs_path"; \
			if test -d "$$abs_path"; then \
				echo "Failed to delete directory '$$abs_path'."; \
				exit 1; \
			fi; \
		fi \
	else \
		echo "Directory '$$abs_path' does not exist."; \
	fi
endef

#################
#     Local     #
#################

# Use the genesis file from the beacond folder as it has been modified by 
# beacond genesis set-deposit-storage.
ETH_GENESIS_PATH = ${HOMEDIR}/eth-genesis.json
NETHER_ETH_GENESIS_PATH = ${HOMEDIR}/eth-nether-genesis.json

HOMEDIR = .tmp/beacond
JWT_PATH = ${TESTAPP_FILES_DIR}/jwt.hex
ETH_DATA_DIR = .tmp/eth-home
DEVNET_CHAIN_SPEC = devnet

## Start an ephemeral `beacond` node. Must be run before running the EL to
## configure the deposit contract storage slots pre-genesis.
start: 
	@JWT_SECRET_PATH=$(JWT_PATH) \
	CHAIN_SPEC=$(DEVNET_CHAIN_SPEC) \
	${TESTAPP_FILES_DIR}/entrypoint.sh

# URLs used for dialing the eth client
IPC_PATH = .tmp/eth-home/eth-engine.ipc
IPC_PREFIX = ipc://

# start-ipc is currently only supported while running eth client the host machine
# Only works with geth-host rn
start-ipc: ## start a local ephemeral `beacond` node with IPC
	@JWT_SECRET_PATH=$(JWT_PATH) \
	RPC_DIAL_URL=${IPC_PATH} \
	RPC_PREFIX=${IPC_PREFIX} \
	${TESTAPP_FILES_DIR}/entrypoint.sh

## Start an ephemeral `reth` node
start-reth: 
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	@docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ghcr.io/paradigmxyz/reth node \
	--chain ${ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

## Start a local ephemeral `reth` node on host machine
start-reth-host: 
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	reth init --datadir ${ETH_DATA_DIR} --chain ${ETH_GENESIS_PATH}
	reth node \
	--chain ${ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

## Start an ephemeral `geth` node with docker
start-geth: 
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go init \
	--datadir ${ETH_DATA_DIR} \
	${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go \
	--syncmode=full \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

## Start a local ephemeral `geth` node on host machine
start-geth-host: 
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	geth init --datadir ${ETH_DATA_DIR} ${ETH_GENESIS_PATH}
	geth \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*"

## Start an ephemeral `nethermind` node
start-nethermind: 
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/${HOMEDIR}:/${HOMEDIR} \
	nethermind/nethermind \
	--JsonRpc.Port 8545 \
	--JsonRpc.EngineEnabledModules "eth,net,engine" \
	--JsonRpc.EnginePort 8551 \
	--JsonRpc.EngineHost 0.0.0.0 \
	--JsonRpc.Host 0.0.0.0 \
	--JsonRpc.JwtSecretFile ../$(JWT_PATH) \
	--Sync.PivotNumber 0 \
	--Init.ChainSpecPath ../$(NETHER_ETH_GENESIS_PATH)

## Start an ephemeral `besu` node
start-besu: 
	$(call ask_reset_dir_func, .tmp/besu)
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/.tmp:/.tmp \
	-v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	hyperledger/besu:latest \
	--data-path=/.tmp/besu \
	--genesis-file=../../${ETH_GENESIS_PATH} \
	--rpc-http-enabled \
	--rpc-http-api=ETH,NET,ENGINE,DEBUG,NET,WEB3 \
	--host-allowlist="*" \
	--rpc-http-cors-origins="all" \
	--engine-rpc-port=8551 \
	--engine-rpc-enabled \
	--engine-host-allowlist="*" \
	--engine-jwt-secret=../../${JWT_PATH}

## Start an ephemeral `erigon` node
start-erigon: 
	$(call ask_reset_dir_func, .tmp/erigon)
	docker run \
	--user 1000:1000 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	erigontech/erigon:latest init \
	--datadir /.tmp/erigon \
	/${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--user 1000:1000 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	erigontech/erigon:latest \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,erigon,engine,web3,net,debug,trace,txpool,admin,ots \
	--http.vhosts "*" \
	--port 30303 \
	--http.corsdomain "*" \
	--http.port 8545 \
	--authrpc.addr	0.0.0.0 \
	--authrpc.jwtsecret /$(JWT_PATH) \
	--authrpc.vhosts "*" \
	--networkid 80087 \
	--db.size.limit	3000MB \
	--datadir /.tmp/erigon

start-ethereumjs:
	$(call ask_reset_dir_func, .tmp/ethereumjs)
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	ethpandaops/ethereumjs:master \
	--gethGenesis ../../${ETH_GENESIS_PATH} \
	--rpcEngine \
	--jwtSecret ../../$(JWT_PATH) \
	--rpcEngineAddr 0.0.0.0 \
	--dataDir ../../.tmp/ethereumjs \
	--isSingleNode \
	--rpc \
	--rpcAddr 0.0.0.0

#################
#    Bepolia    #
#################

TESTNET_CHAIN_SPEC = testnet
BEPOLIA_NETWORK_FILES_DIR = ${TESTAPP_FILES_DIR}/../networks/80069
BEPOLIA_ETH_GENESIS_PATH = ${BEPOLIA_NETWORK_FILES_DIR}/eth-genesis.json

start-bepolia:
	@JWT_SECRET_PATH=$(JWT_PATH) \
	CHAIN_SPEC=$(TESTNET_CHAIN_SPEC) \
	${TESTAPP_FILES_DIR}/entrypoint.sh

# NOTE: Peers are used as bootnodes for the Geth node.
start-geth-bepolia:
	# TODO: Update to use latest Geth once ready
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BEPOLIA_NETWORK_FILES_DIR}:/${BEPOLIA_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go:v1.14.13 init \
	--datadir ${ETH_DATA_DIR} \
	${BEPOLIA_ETH_GENESIS_PATH}

	@# Read bootnodes from the file; the file is mounted into the container.
	@bootnodes=`cat $(PWD)/$(BEPOLIA_NETWORK_FILES_DIR)/el-peers.txt`; \
	echo "Using bootnodes: $$bootnodes"; \
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BEPOLIA_NETWORK_FILES_DIR}:/${BEPOLIA_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go:v1.14.13 \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--syncmode=full \
	--bootnodes $$bootnodes

start-reth-bepolia:
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	@trustedpeers=`cat $(PWD)/$(BEPOLIA_NETWORK_FILES_DIR)/el-peers.txt`; \
	echo "Using truted peers: $$trustedpeers"; \
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BEPOLIA_NETWORK_FILES_DIR}:/${BEPOLIA_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ghcr.io/paradigmxyz/reth node \
	--chain ${BEPOLIA_ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--trusted-peers $$trustedpeers

#################
#    Mainnet    #
#################

MAINNET_CHAIN_SPEC = mainnet
MAINNET_NETWORK_FILES_DIR = ${TESTAPP_FILES_DIR}/../networks/80094
MAINNET_ETH_GENESIS_PATH = ${MAINNET_NETWORK_FILES_DIR}/eth-genesis.json

start-mainnet:
	@JWT_SECRET_PATH=$(JWT_PATH) \
	CHAIN_SPEC=$(MAINNET_CHAIN_SPEC) \
	${TESTAPP_FILES_DIR}/entrypoint.sh

# NOTE: By default this will use the EL peers as your bootnodes. If you want specific 
# discovery bootnodes by region, refer to testing/networks/80094/el-bootnodes.txt
start-geth-mainnet:
	# TODO: Update to use latest Geth once ready
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${MAINNET_NETWORK_FILES_DIR}:/${MAINNET_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go:v1.14.13 init \
	--datadir ${ETH_DATA_DIR} \
	${MAINNET_ETH_GENESIS_PATH}

	@# Read bootnodes from the file; the file is mounted into the container.
	@bootnodes=`cat $(PWD)/$(MAINNET_NETWORK_FILES_DIR)/el-peers.txt`; \
	echo "Using bootnodes: $$bootnodes"; \
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BEPOLIA_NETWORK_FILES_DIR}:/${BEPOLIA_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go:v1.14.13 \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--syncmode=full \
	--bootnodes $$bootnodes

start-reth-mainnet:
	$(call ask_reset_dir_func, $(ETH_DATA_DIR))
	@trustedpeers=`cat $(PWD)/$(MAINNET_NETWORK_FILES_DIR)/el-peers.txt`; \
	echo "Using truted peers: $$trustedpeers"; \
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${MAINNET_NETWORK_FILES_DIR}:/${MAINNET_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ghcr.io/paradigmxyz/reth node \
	--chain ${MAINNET_ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--trusted-peers $$trustedpeers

#################
#    Testing    #
#################

SHORT_FUZZ_TIME=10s
MEDIUM_FUZZ_TIME=30s
LONG_FUZZ_TIME=3m

# Define a function to filter out lines with "/testing/", "/mock/", "/mocks/", or ".mock.go"
define FILTER_COVERAGE
	grep -Ev '(/testing/|/mock/|/mocks/|\.mock\.go)' $(1) > $(2)
endef

test:
	@$(MAKE) test-unit test-forge-fuzz

test-unit-no-coverage: ## run golang unit tests
	@echo "Running unit tests..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -race -tags bls12381,test

coverage-summary: test-unit test-simulated
	@echo "Merging coverage reports..."
	@go install github.com/wadey/gocovmerge@latest
	@gocovmerge test-unit-cover.txt test-simulated.txt > coverage-merged.txt
	@echo "Coverage Summary:"
	@go tool cover -html=coverage-merged.txt

test-unit-cover: test-unit test-simulated test-unit-quick ## run golang unit tests with coverage

test-unit:
	@echo "Running unit tests with coverage and race checks..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -race -covermode=atomic -coverpkg=github.com/berachain/beacon-kit/... -coverprofile=temp-test-unit-cover.txt -tags bls12381,test
	# Filter out any coverage lines from the testing directory
	$(call FILTER_COVERAGE, temp-test-unit-cover.txt, test-unit-cover.txt)
	@rm temp-test-unit-cover.txt

test-unit-quick: ## run quick tests. We run these without coverage as covermode=atomic is too slow and coverage here provides little value
	@echo "Running 'quick' tests..."
	@go list -f '{{.Dir}}/testing/quick' -m | xargs \
		go test -v -tags quick

test-simulated: ## run simulation tests
	@echo "Running simulation tests with coverage"
	@go list -f '{{.Dir}}/testing/simulated' -m | xargs \
		go test -cover -covermode=atomic -coverpkg=github.com/berachain/beacon-kit/... -coverprofile=temp-test-simulated.txt -tags simulated -v
	# Filter out any coverage lines from the testing directory
	$(call FILTER_COVERAGE, temp-test-simulated.txt, test-simulated.txt)
	@rm temp-test-simulated.txt

test-unit-bench: ## run golang unit benchmarks
	@echo "Running unit tests with benchmarks..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -bench=. -run=^$ -benchmem -tags bls12381,test

# On MacOS, if there is a linking issue on the fuzz tests,
# use the old linker with flags -ldflags=-extldflags=-Wl,-ld_classic
test-unit-fuzz: ## run fuzz tests
	@echo "Running fuzz tests with coverage..."
	go test -run ^FuzzPayloadIDCacheBasic -fuzztime=${SHORT_FUZZ_TIME} github.com/berachain/beacon-kit/payload/cache
	go test -run ^FuzzPayloadIDInvalidInput -fuzztime=${SHORT_FUZZ_TIME} github.com/berachain/beacon-kit/payload/cache
	go test -run ^FuzzPayloadIDCacheConcurrency -fuzztime=${SHORT_FUZZ_TIME} github.com/berachain/beacon-kit/payload/cache
	go test -run ^FuzzHashTreeRoot -fuzztime=${MEDIUM_FUZZ_TIME} github.com/berachain/beacon-kit/primitives/merkle

test-e2e: ## run e2e tests
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-no-build

test-e2e-no-build:
	go test -timeout 0 -tags e2e,bls12381,test ./testing/e2e/. -v

test-e2e-4844: ## run e2e tests
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-4844-no-build

test-e2e-4844-no-build:
	go test -timeout 0 -tags e2e,bls12381,test ./testing/e2e/. -v -testify.m Test4844Live

test-e2e-deposits: ## run e2e tests
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-deposits-no-build

test-e2e-deposits-no-build:
	go test -timeout 0 -tags e2e,bls12381,test ./testing/e2e/. -v -testify.m TestDepositRobustness
