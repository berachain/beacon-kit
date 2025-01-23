#!/usr/bin/make -f



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

#################
#    beacond     #
#################

DEVNET_CHAIN_SPEC = mainnet
JWT_PATH = ${TESTAPP_FILES_DIR}/jwt.hex
ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/eth-genesis.json
NETHER_ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/eth-nether-genesis.json
ETH_DATA_DIR = .tmp/eth-home
# URLs used for dialing the eth client
IPC_PATH = .tmp/eth-home/eth-engine.ipc
HTTP_URL = localhost:8551
IPC_PREFIX = ipc://
HTTP_PREFIX = http://

#################
#    bartio     #
#################

TESTNET_CHAIN_SPEC = testnet
BARTIO_NETWORK_FILES_DIR = ${TESTAPP_FILES_DIR}/../networks/80084
BARTIO_ETH_GENESIS_PATH = ${BARTIO_NETWORK_FILES_DIR}/eth-genesis.json

## Testing:
start: ## start an ephemeral `beacond` node
	@JWT_SECRET_PATH=$(JWT_PATH) \
	CHAIN_SPEC=$(DEVNET_CHAIN_SPEC) \
	${TESTAPP_FILES_DIR}/entrypoint.sh

start-bartio:
	@JWT_SECRET_PATH=$(JWT_PATH) \
	CHAIN_SPEC=$(TESTNET_CHAIN_SPEC) \
	${TESTAPP_FILES_DIR}/entrypoint.sh

# start-ipc is currently only supported while running eth client the host machine
# Only works with geth-host rn
start-ipc: ## start a local ephemeral `beacond` node with IPC
	@JWT_SECRET_PATH=$(JWT_PATH) \
	RPC_DIAL_URL=${IPC_PATH} \
	RPC_PREFIX=${IPC_PREFIX} \
	${TESTAPP_FILES_DIR}/entrypoint.sh

start-reth: ## start an ephemeral `reth` node
	# @rm -rf ${ETH_DATA_DIR}
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
	--ipcpath ${IPC_PATH} \
	--trusted-peers enode://0c5a4a3c0e81fce2974e4d317d88df783731183d534325e32e0fdf8f4b119d7889fa254d3a38890606ec300d744e2aa9c87099a4a032f5c94efe53f3fcdfecfe@34.22.104.177:30303,enode://b6a3137d3a36ef37c4d31843775a9dc293f41bcbde33b6309c80b1771b6634827cd188285136a57474427bd8845adc2f6fe2e0b106bd58d14795b08910b9c326@34.64.247.85:30303,enode://0b6633300614bc2b9749aee0cace7a091ec5348762aee7b1d195f7616d03a9409019d9bef336624bab72e0d069cd4cf0b0de6fbbf53f04f6b6e4c5b39c6bdca6@34.22.73.21:30303,enode://552b001abebb5805fcd734ad367cd05d9078d18f23ec598d7165460fadcfc51116ad95c418f7ea9a141aa8cbc496c8bea3322b67a5de0d3380f11aab1a797513@34.64.37.55:30303,enode://5b037f66099d5ded86eb7e1619f6d06ceb15609e8cc345ced22a4772b06178004e1490a3cd32fd1222789de4c6e4021c2d648a3d750f6d5323e64b771bbd8de7@35.247.182.34:30303,enode://846db253c53753d3ea1197aec296306dc84c25f3afdf142b65cb0fe0f984de55072daa3bbf05a9aea046a38a2292403137b6eafefd5646fcf62120b74e3b898d@34.87.9.231:30303,enode://64b7f6ee9bcd942ad4949c70f2077627f078a057dfd930e6e904e12643d8952f5ae87c91e24559765393f244a72c9d5c011d7d5176e59191d38f315db85a20f5@34.126.78.49:30303,enode://cf4d19bfb8ec507427ec882bac0bac85a0c8c9ddaa0ec91b773bb614e5e09d107cd9fbe323b96f62f31c493f8f42cc5495c18b87c08560c5dea1dfd25256dcf6@35.240.200.36:30303,enode://bb7e44178543431feac8f0ee3827056b7b84d8235b802a8bdbbcd4939dab7f7dd2579ff577a38b002bb0139792af67abd2dd5c9f4f85b8da6e914fa76dca82bc@34.40.14.50:30303,enode://8fef1f5df45e7b31be00a21e1da5665d5a5f5bf4c379086b843f03eade941bdd157f08c95b31880c492577edb9a9b185df7191eaebf54ab06d5bd683b289f3af@35.246.168.217:30303,enode://ce9c87cfe089f6811d26c96913fa3ec10b938d9017fc6246684c74a33679ee34ceca9447180fb509e37bf2b706c2877a82085d34bfd83b5b520ee1288b0fc32f@34.40.28.159:30303,enode://713657eb6a53feadcbc47e634ad557326a51eb6818a3e19a00a8111492f50a666ccbf2f5d334d247ecf941e68d242ef5c3b812b63c44d381ef11f79c2cdb45c7@35.234.82.236:30303,enode://d071fa740e063ce1bb9cdc2b7937baeff6dc4000f91588d730a731c38a6ff0d4015814812c160fab8695e46f74b9b618735368ea2f16db4d785f16d29b3fb7b0@35.203.86.197:30303,enode://ffc452fe451a2e5f89fe634744aea334d92dcd30d881b76209d2db7dbf4b7ee047e7c69a5bb1633764d987a7441d9c4bc57ccdbfd6442a2f860bf953bc89a9b9@34.118.187.161:30303,enode://da94328302a1d1422209d1916744e90b6095a48b2340dcec39b22002c098bb4d58a880dab98eb26edf03fa4705d1b62f99a8c5c14e6666e4726b6d3066d8a4d7@34.95.30.190:30303,enode://19c7671a4844699b481e81a5bcfe7bafc7fefa953c16ebbe1951b1046371e73839e9058de6b7d3c934318fe7e7233dde3621c1c1018eb8b294ea3d4516147150@34.47.60.196:30303

start-reth-bartio:
	@rm -rf ${ETH_DATA_DIR}
	@docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BARTIO_NETWORK_FILES_DIR}:/${BARTIO_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ghcr.io/paradigmxyz/reth node \
	--chain ${BARTIO_ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

start-reth-host: ## start a local ephemeral `reth` node on host machine
	rm -rf ${ETH_DATA_DIR}
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

start-geth: ## start an ephemeral `geth` node with docker
	rm -rf ${ETH_DATA_DIR}
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
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

start-geth-bartio:
	rm -rf ${ETH_DATA_DIR}
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BARTIO_NETWORK_FILES_DIR}:/${BARTIO_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go init \
	--datadir ${ETH_DATA_DIR} \
	${BARTIO_ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	--rm -v $(PWD)/${BARTIO_NETWORK_FILES_DIR}:/${BARTIO_NETWORK_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH}

start-geth-host: ## start a local ephemeral `geth` node on host machine
	rm -rf ${ETH_DATA_DIR}
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

start-nethermind: ## start an ephemeral `nethermind` node
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	nethermind/nethermind \
	--JsonRpc.Port 8545 \
	--JsonRpc.EngineEnabledModules "eth,net,engine" \
	--JsonRpc.EnginePort 8551 \
	--JsonRpc.EngineHost 0.0.0.0 \
	--JsonRpc.Host 0.0.0.0 \
	--JsonRpc.JwtSecretFile ../$(JWT_PATH) \
	--Sync.PivotNumber 0 \
	--Init.ChainSpecPath ../$(TESTAPP_FILES_DIR)/eth-nether-genesis.json

start-besu: ## start an ephemeral `besu` node
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	hyperledger/besu:latest \
	--data-path=.tmp/besu \
	--genesis-file=../../${ETH_GENESIS_PATH} \
	--rpc-http-enabled \
	--rpc-http-api=ETH,NET,ENGINE,DEBUG,NET,WEB3 \
	--host-allowlist="*" \
	--rpc-http-cors-origins="all" \
	--engine-rpc-port=8551 \
	--engine-rpc-enabled \
	--engine-host-allowlist="*" \
	--engine-jwt-secret=../../${JWT_PATH}

start-erigon: ## start an ephemeral `erigon` node
	rm -rf .tmp/erigon
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
	rm -rf .tmp/ethereumjs
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
	--dataDir .tmp/ethereumjs \
	--isSingleNode \
	--rpc \
	--rpcAddr 0.0.0.0

SHORT_FUZZ_TIME=10s
MEDIUM_FUZZ_TIME=30s
LONG_FUZZ_TIME=3m


test:
	@$(MAKE) test-unit test-forge-fuzz

test-unit: ## run golang unit tests
	@echo "Running unit tests..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -race -tags bls12381

test-unit-cover: ## run golang unit tests with coverage
	@echo "Running unit tests with coverage..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -race -coverprofile=test-unit-cover.txt -tags bls12381

test-unit-bench: ## run golang unit benchmarks
	@echo "Running unit tests with benchmarks..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -bench=. -run=^$ -benchmem -tags bls12381

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
	go test -timeout 0 -tags e2e,bls12381 ./testing/e2e/. -v

test-e2e-4844: ## run e2e tests
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-4844-no-build

test-e2e-4844-no-build:
	go test -timeout 0 -tags e2e,bls12381 ./testing/e2e/. -v -testify.m Test4844Live