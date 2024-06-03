#!/usr/bin/make -f



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

#################
#    beacond     #
#################

JWT_PATH = ${TESTAPP_FILES_DIR}/jwt.hex
ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/genesis.eth.json
NETHER_ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR}/nether-genesis.eth.json
ETH_DATA_DIR = .tmp/eth-home

# URLs used for dialing the eth client
IPC_PATH = .tmp/eth-home/eth-engine.ipc
HTTP_URL = localhost:8551
IPC_PREFIX = ipc://
HTTP_PREFIX = http://

# Variables
PREDEPLOY_ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02,0x4e59b44847b379578588920cA78FbF26c0B4956C,0x00000000219ab540356cbb839cbe05303d7705fa
NONCE = 1,1,1
PREDEPLOY_BALANCE = 0,0,0
CODE = 0x3373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500,0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3,0x608060405260043610610028575f3560e01c80632dfdf0b51461002c5780635b70fa2914610068575b5f80fd5b348015610037575f80fd5b505f5461004b9067ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200160405180910390f35b61007b610076366004610318565b61007d565b005b603086146100b7576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602084146100f1576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060811461012b576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6101358461020d565b905064077359400067ffffffffffffffff82161015610180576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008116600167ffffffffffffffff928316908101909216179091556040517f68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46916101fb918b918b918b918b9188918b918b919061040a565b60405180910390a15050505050505050565b5f61021c633b9aca003461049c565b15610253576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610262633b9aca00346104af565b905067ffffffffffffffff8111156102a6576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6102b05f346102b6565b92915050565b5f385f3884865af16102cf5763b12d13eb5f526004601cfd5b5050565b5f8083601f8401126102e3575f80fd5b50813567ffffffffffffffff8111156102fa575f80fd5b602083019150836020828501011115610311575f80fd5b9250929050565b5f805f805f805f6080888a03121561032e575f80fd5b873567ffffffffffffffff80821115610345575f80fd5b6103518b838c016102d3565b909950975060208a0135915080821115610369575f80fd5b6103758b838c016102d3565b909750955060408a01359150808216821461038e575f80fd5b909350606089013590808211156103a3575f80fd5b506103b08a828b016102d3565b989b979a50959850939692959293505050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60a081525f61041d60a083018a8c6103c3565b828103602084015261043081898b6103c3565b905067ffffffffffffffff808816604085015283820360608501526104568287896103c3565b9250808516608085015250509998505050505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f826104aa576104aa61046f565b500690565b5f826104bd576104bd61046f565b50049056fea264697066735822122022654475565591de6c8c68085b8494831b4b47f9ad39618a8f84e361e7d0382464736f6c63430008190033
ACCOUNT = 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
BALANCE = 100

generate-genesis-cmd:
	@if [ ! -f $(OUTPUT) ]; then \
		go run ./testing/generate-genesis $(FORMAT) --predeployAddresses $(PREDEPLOY_ADDRESS) \
		--predeployNonces $(NONCE) --predeployBalances $(PREDEPLOY_BALANCE) --predeployCodes $(CODE) \
		--accountAddresses $(ACCOUNT) --accountBalances $(BALANCE) --output $(OUTPUT); \
	fi
	

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
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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

start-geth-host: ## start a local ephemeral `geth` node on host machine
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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
	$(MAKE) FORMAT=nethermind OUTPUT=${NETHER_ETH_GENESIS_PATH} generate-genesis-cmd
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
	--Init.ChainSpecPath ../$(NETHER_ETH_GENESIS_PATH)

start-besu: ## start an ephemeral `besu` node
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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
	--networkid 80086 \
	--db.size.limit	3000MB \
	--datadir .tmp/erigon

start-ethereumjs:
	$(MAKE) FORMAT=geth OUTPUT=${ETH_GENESIS_PATH} generate-genesis-cmd
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