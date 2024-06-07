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

## Testing:
start: ## start an ephemeral `beacond` node
	@JWT_SECRET_PATH=$(JWT_PATH) ${TESTAPP_FILES_DIR}/entrypoint.sh

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
	--ipcpath ${IPC_PATH} \
	--bootnodes "enode://7a2f67d22b12e10c6ba9cd951866dda6471604be5fbd5102217dbad1cc56e590befd2009ecc99958a468a5b8e0dc28e14d9b6822491719c93199be6aa0319077@34.124.220.31:30303,enode://a96aac0b81c7e75fecc2ae613eaf13b27b2aaf3d46a90db904f94797d1746aa31e6593ae4cd476f81d5c6d1d2228ca60c885727978c369586c38871c63a330ee@35.240.182.27:30303,enode://dc44744074ac2dd76db0e0f9d95eb86cd558f6ba75e4a4af1303f2259624c8ce041198f976862a284165253b6dc6b2fa91b995cbca3ef2683879b6247e05e553@34.95.61.239:30303,enode://bf5364e1cf7ecd11646ccaea5c06b56622c04d52200d9cd141e01db9c9661237ceebecde1616e66e390a968ffd1c07e027531cad23044517b7bf36caa8b97f5f@34.152.41.26:30303,enode://e31aa249638083d34817eed2b499ccd4b0718a332f0ea530e3062e13f624cb03a7d6b6e0460193ee87b5fc12e73a726070a3126ef53492ffbdc5e6c102f6dfb3@34.64.198.56:30303,enode://3f2f85e2e711f198fb7324b74fab6a0599b2534774f3aa26241dbbabe870b650574324da01aa98ee24ce97c8d76362a2db03034a6ddff43119ccfdc269663cbf@34.47.79.13:30303,enode://0401e494dbd0c84c5c0f72adac5985d2f2525e08b68d448958aae218f5ac8198a80d1498e0ebec2ce38b1b18d6750f6e61a56b4614c5a6c6cf0981c39aed47dc@34.159.32.127:30303,enode://9b6c1eb143c9e3af0c7283262a9a38fe8bf844114b1f304673c2ac1c23e6bccfdaa8f4e9cb8c460bded495933fd92eeff30e6ab2e0538b56e249beea2c512906@35.234.88.149:30303"

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
	--http.api eth,net,admin \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--syncmode full \
	--bootnodes enode://7a2f67d22b12e10c6ba9cd951866dda6471604be5fbd5102217dbad1cc56e590befd2009ecc99958a468a5b8e0dc28e14d9b6822491719c93199be6aa0319077@34.124.220.31:30303,enode://a96aac0b81c7e75fecc2ae613eaf13b27b2aaf3d46a90db904f94797d1746aa31e6593ae4cd476f81d5c6d1d2228ca60c885727978c369586c38871c63a330ee@35.240.182.27:30303,enode://dc44744074ac2dd76db0e0f9d95eb86cd558f6ba75e4a4af1303f2259624c8ce041198f976862a284165253b6dc6b2fa91b995cbca3ef2683879b6247e05e553@34.95.61.239:30303,enode://bf5364e1cf7ecd11646ccaea5c06b56622c04d52200d9cd141e01db9c9661237ceebecde1616e66e390a968ffd1c07e027531cad23044517b7bf36caa8b97f5f@34.152.41.26:30303,enode://e31aa249638083d34817eed2b499ccd4b0718a332f0ea530e3062e13f624cb03a7d6b6e0460193ee87b5fc12e73a726070a3126ef53492ffbdc5e6c102f6dfb3@34.64.198.56:30303,enode://3f2f85e2e711f198fb7324b74fab6a0599b2534774f3aa26241dbbabe870b650574324da01aa98ee24ce97c8d76362a2db03034a6ddff43119ccfdc269663cbf@34.47.79.13:30303,enode://0401e494dbd0c84c5c0f72adac5985d2f2525e08b68d448958aae218f5ac8198a80d1498e0ebec2ce38b1b18d6750f6e61a56b4614c5a6c6cf0981c39aed47dc@34.159.32.127:30303,enode://9b6c1eb143c9e3af0c7283262a9a38fe8bf844114b1f304673c2ac1c23e6bccfdaa8f4e9cb8c460bded495933fd92eeff30e6ab2e0538b56e249beea2c512906@35.234.88.149:30303

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
    thorax/erigon:latest init \
    --datadir .tmp/erigon \
    ${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR}:/${TESTAPP_FILES_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	thorax/erigon:latest \
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
	--networkid 80084 \
	--db.size.limit	3000MB \
	--datadir .tmp/erigon

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