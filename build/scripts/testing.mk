#!/usr/bin/make -f



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

#################
#    beacond     #
#################

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
BOOT_NODES = "enode://dea0da1a78ec8534ae607d330996f7a1bcf83ca092bbc6027a54d1e53a4d175f31507ed370f62783168809c6fd839c61dc4e89297e550ee9a57a5f33451ef97c@10.0.11.187:30303,enode://13204e74fd935d7c42759a423534ce684b7fb0d0d8eb94df952198cf18217796e4eaedbad7c50c27aa6bdf31d01d5e60cbc161aaee22bb9f9765a4d3eb9dd251@10.0.4.236:30303,enode://156db9abc289625b9848c7de18d4cfc709c1bb4d17dcbdd9c686f00c27df5b9aacec90f9ddd1a86ae29c1714c150551152ae8f07776e94be32937bceef4f85ce@10.0.8.119:30303"

## Testing:
start: ## start an ephemeral `beacond` node
	@JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh

start-validator-1: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 1
start-validator-2: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 2
start-validator-3: ## start an ephemeral `beacond` node
	JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh 1 validator 3
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
	--http.api eth,net,debug,txpool \
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
	--bootnodes "${BOOT_NODES}" \
	--nat extip:${INTERNAL_IP}

start-geth-node-snap-run:
	sudo chmod 777 -R .tmp
	geth \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net,debug,txpool,admin \
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
	ethpandaops/ethereumjs:stable \
	--gethGenesis ../../${ETH_GENESIS_PATH} \
	--rpcEngine \
	--jwtSecret ../../$(JWT_PATH) \
	--rpcEngineAddr 0.0.0.0 \
	--dataDir .tmp/ethereumjs

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