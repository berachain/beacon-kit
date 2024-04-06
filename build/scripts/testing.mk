#!/usr/bin/make -f



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

#################
#    beacond     #
#################

# TODO: add start-erigon

JWT_PATH = ${TESTAPP_DIR}/jwt.hex
ETH_GENESIS_PATH = ${TESTAPP_DIR}/eth-genesis.json

# Start beacond
start:
	@JWT_SECRET_PATH=$(JWT_PATH) ./beacond/entrypoint.sh

# Start reth node
start-reth:
	@rm -rf .tmp/eth-home
	@docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	ghcr.io/paradigmxyz/reth node \
	--chain ${ETH_GENESIS_PATH} \
	--http \
	--http.addr "0.0.0.0" \
	--http.api eth,net \
	--authrpc.addr "0.0.0.0" \
	--authrpc.jwtsecret $(JWT_PATH) \
	
# Init and start geth node
start-geth:
	rm -rf .tmp/geth
	docker run \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go init \
	--datadir .tmp/geth \
	${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go \
	--http \
	--http.addr 0.0.0.0 \
	--http.api eth,net \
	--authrpc.addr 0.0.0.0 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir .tmp/geth

# Start nethermind node
start-nethermind:
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	nethermind/nethermind \
	--JsonRpc.Port 8545 \
	--JsonRpc.EngineEnabledModules "eth,net,engine" \
	--JsonRpc.EnginePort 8551 \
	--JsonRpc.EngineHost 0.0.0.0 \
	--JsonRpc.Host 0.0.0.0 \
	--JsonRpc.JwtSecretFile ../$(JWT_PATH) \
	--Sync.PivotNumber 0 \
	--Init.ChainSpecPath ../$(TESTAPP_DIR)/eth-nether-genesis.json

# Start besu node
start-besu:
	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	-v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
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


###############################################################################
###                                Testing                                  ###
###############################################################################


#################
#      unit     #
#################

SHORT_FUZZ_TIME=10s
MEDIUM_FUZZ_TIME=30s
LONG_FUZZ_TIME=3m

test:
	@$(MAKE) test-unit test-forge-fuzz
	
test-unit:
	@$(MAKE)
	@echo "Running unit tests..."
	go test ./...

test-unit-cover:
	@$(MAKE)
	@echo "Running unit tests with coverage..."
	go test -race -coverprofile=test-unit-cover.txt -covermode=atomic ./...

# On MacOS, if there is a linking issue on the fuzz tests, 
# use the old linker with flags -ldflags=-extldflags=-Wl,-ld_classic
test-unit-fuzz:
	@echo "Running fuzz tests with coverage..."
	go test ./mod/runtime/services/builder/local/cache/... -fuzz=FuzzPayloadIDCacheBasic -fuzztime=${SHORT_FUZZ_TIME}
	go test ./mod/runtime/services/builder/local/cache/... -fuzz=FuzzPayloadIDInvalidInput -fuzztime=${SHORT_FUZZ_TIME}
	go test ./mod/runtime/services/builder/local/cache/... -fuzz=FuzzPayloadIDCacheConcurrency -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzHashTreeRoot ./mod/trie/merkleize/... -fuzztime=${MEDIUM_FUZZ_TIME}
	go test -fuzz=FuzzQueueSimple ./mod/storage/statedb/collections/ -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzQueueMulti ./mod/storage/statedb/collections/ -fuzztime=${SHORT_FUZZ_TIME}

#################
#      e2e      #
#################

test-e2e:
	@$(MAKE) build-docker VERSION=kurtosis-local test-e2e-no-build

test-e2e-no-build:
	go test -tags e2e ./testing/e2e/. -v