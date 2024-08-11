#!/usr/bin/make -f



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

#################
#    beacond     #
#################

DEVNET_CHAIN_SPEC = devnet
JWT_PATH = ${TESTAPP_FILES_DIR}/jwt.hex
ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/eth-genesis.json
NETHER_ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/eth-nether-genesis.json
ETH_DATA_DIR = .tmp/eth-home
# URLs used for dialing the eth client
IPC_PATH = .tmp/eth-home/eth-engine.ipc
HTTP_URL = localhost:8551
IPC_PREFIX = ipc://
HTTP_PREFIX = http://
INTERNAL_IP := $(shell hostname -I | awk '{print $$1}')
BOOT_NODES = "enode://c32896cefb5006afc1a4400c5a87c17687740a0c64ed4273ee57fb2b53de889ad5807ac30b1dd09ce77ac128df17214626e69a1985bbad95ffec9bdcf3e8d5ec@10.0.2.106:30303,enode://61105d68f4dd9afb774b3bc7f6575730b7a3a0f0e5d6045d97ff5aa098cdcff31bafe80eb0c45f403d0cfe1b0c401538aa172ba6553390ad1770e7ed024c2932@10.0.13.254:30303,enode://d28d096fffb33467149a53ff133f1c3b9635f7333757010e41744cef77724e08b386938cf77f7a3fc6240142ce06f3694c7980f1e36a3d0d5628a01d9d884d97@10.0.1.173:30303,enode://724999756780c9ca03551d6c5ba0c92a115fbcca880557cdcd43f3ee4886fea1f458df4bca416d306e9ab9b9c7fc1238f617dad7818dfbb9493c7bfcc81f5011@10.0.15.167:30303"

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

start-validator-locally: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator locally
start-validator-1: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 1
start-validator-2: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 2
start-validator-3: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 3
start-validator-4: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 4
start-node-init: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 node onlyInit
start-node-run: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 node

# start-ipc is currently only supported while running eth client the host machine
# Only works with geth-host rn
start-ipc: ## start a local ephemeral `beacond` node with IPC
	@JWT_SECRET_PATH=$(JWT_PATH) \
	RPC_DIAL_URL=${IPC_PATH} \
	RPC_PREFIX=${IPC_PREFIX} \
	${TESTAPP_FILES_DIR}/entrypoint.sh 

start-reth: ## start an ephemeral `reth` node
	@rm -rf ${ETH_DATA_DIR}
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

start-geth-init:
	sudo chmod 777 -R .tmp
	rm -rf ${ETH_DATA_DIR}

	geth init --state.scheme "hash" --datadir ${ETH_DATA_DIR} ${ETH_GENESIS_PATH}

start-geth-init-local:
	geth init --state.scheme "hash" --datadir ${ETH_DATA_DIR} ${ETH_GENESIS_PATH}

start-geth-archive-run:
	sudo chmod 777 -R .tmp
	geth \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,web3,net,debug,txpool,admin \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--rpc.allow-unprotected-txs \
	--syncmode "full" \
	--state.scheme "hash" \
	--gcmode "archive" \
	--rpc.evmtimeout "10s" \
	--txpool.globalslots 1000000 \
	--txpool.globalqueue 3000000 \
	--http.vhosts "*" \
	--miner.gasprice 100000000 \
	--bootnodes "${BOOT_NODES}" \
	--nat extip:${INTERNAL_IP}

start-geth-node-snap-run:
	sudo chmod 777 -R .tmp
	geth \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,web3,net,debug,txpool,admin \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--rpc.allow-unprotected-txs \
	--syncmode "snap" \
	--state.scheme "hash" \
    --rpc.evmtimeout "10s" \
	--txpool.globalslots 1000000 \
	--txpool.globalqueue 3000000 \
	--http.vhosts "*" \
	--miner.gasprice 100000000 \
	--bootnodes "${BOOT_NODES}" \
	--nat extip:${INTERNAL_IP}

start-geth-run-local:
	geth \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net,debug,txpool \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--rpc.allow-unprotected-txs \
	--syncmode "snap"

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
    --rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
    -v $(PWD)/.tmp:/.tmp \
    thorax/erigon:v2.60.2 init \
    --datadir .tmp/erigon \
    ${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	thorax/erigon:v2.60.2 \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--http.vhosts "*" \
	--port 30303 \
	--http.corsdomain "*" \
	--http.port 8545 \
	--authrpc.addr	0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--networkid 2012 \
	--db.size.limit	3000MB \
	--datadir .tmp/erigon

start-erigon-init:
	sudo chmod 777 -R .tmp
	rm -rf .tmp/erigon
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	thorax/erigon:v2.60.2 init \
	--datadir .tmp/erigon \
	${ETH_GENESIS_PATH}

start-erigon-validator-run:
	sudo chmod 777 -R .tmp
	docker run --name execution-erigon \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -d -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	thorax/erigon:v2.60.2 \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,erigon,web3,net,debug,trace,txpool \
	--http.vhosts "*" \
	--port 30303 \
	--http.corsdomain "*" \
	--http.port 8545 \
	--authrpc.addr	0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--networkid 2012 \
	--db.size.limit	3000MB \
	--datadir .tmp/erigon \
	--rpc.allow-unprotected-txs \
	--miner.gaslimit 100000000 \
	--nat extip:${INTERNAL_IP}

start-erigon-node-run:
	sudo chmod 777 -R .tmp
	docker run --name execution-erigon \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -d -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	thorax/erigon:v2.60.2 \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,erigon,web3,net,debug,trace,txpool \
	--http.vhosts "*" \
	--port 30303 \
	--http.corsdomain "*" \
	--http.port 8545 \
	--authrpc.addr	0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--networkid 2012 \
	--db.size.limit	3000MB \
	--datadir .tmp/erigon \
	--rpc.allow-unprotected-txs \
	--nat extip:${INTERNAL_IP} \
	--verbosity "debug" \
    --log.console.verbosity "debug" \
    --miner.gaslimit 100000000 \
	--trustedpeers "" \
	--bootnodes ""

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
		go test

test-unit-cover: ## run golang unit tests with coverage
	@echo "Running unit tests with coverage..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -race -coverprofile=test-unit-cover.txt 

test-unit-bench: ## run golang unit benchmarks
	@echo "Running unit tests with benchmarks..."
	@go list -f '{{.Dir}}/...' -m | xargs \
		go test -bench=. -run=^$ -benchmem

# On MacOS, if there is a linking issue on the fuzz tests, 
# use the old linker with flags -ldflags=-extldflags=-Wl,-ld_classic
test-unit-fuzz: ## run fuzz tests
	@echo "Running fuzz tests with coverage..."
	go test ./mod/payload/pkg/cache/... -fuzz=FuzzPayloadIDCacheBasic -fuzztime=${SHORT_FUZZ_TIME}
	go test ./mod/payload/pkg/cache/... -fuzz=FuzzPayloadIDInvalidInput -fuzztime=${SHORT_FUZZ_TIME}
	go test ./mod/payload/pkg/cache/... -fuzz=FuzzPayloadIDCacheConcurrency -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzHashTreeRoot ./mod/primitives/pkg/merkle -fuzztime=${MEDIUM_FUZZ_TIME}





test-e2e: ## run e2e tests
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-no-build

test-e2e-no-build:
	go test -tags e2e,bls12381 ./testing/e2e/. -v