#!/usr/bin/make -f
include build/scripts/cosmos.mk build/scripts/constants.mk build/scripts/docker.mk

# Specify the default target if none is provided
.DEFAULT_GOAL := build

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(OUT_DIR)/beacond

build-linux-amd64:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

build-linux-arm64:
	GOOS=linux GOARCH=arm64 LEDGER_ENABLED=false $(MAKE) build

$(BUILD_TARGETS): $(OUT_DIR)/
	@echo "Building ${TESTAPP_CMD_DIR}"
	@cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./.

$(OUT_DIR)/:
	mkdir -p $(OUT_DIR)/

clean:
	@rm -rf .tmp/ 
	@rm -rf $(OUT_DIR)
	@$(MAKE) sszgen-clean proto-clean forge-clean

#################
#     forge     #
#################

forge-build: |
	@forge build --extra-output-files bin --extra-output-files abi  --root $(CONTRACTS_DIR)

forge-clean: |
	@forge clean --root $(CONTRACTS_DIR)


###############################################################################
###                                 CodeGen                                 ###
###############################################################################

generate:
	@$(MAKE) abigen-install mockery 
	@for module in $(MODULES); do \
		echo "Running go generate in $$module"; \
		(cd $$module && go generate ./...) || exit 1; \
	done
	@$(MAKE) sszgen

abigen-install:
	@echo "--> Installing abigen"
	@go install github.com/ethereum/go-ethereum/cmd/abigen@latest

mockery-install:
	@echo "--> Installing mockery"
	@go install github.com/vektra/mockery/v2@latest

mockery:
	@$(MAKE) mockery-install
	@echo "Running mockery..."
	@mockery


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
	@JWT_SECRET_PATH=$(JWT_PATH) ./examples/beacond/entrypoint.sh

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
	--http.api eth \
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
	--http.api eth \
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
#     unit      #
#################

SHORT_FUZZ_TIME=15s
MEDIUM_FUZZ_TIME=45s
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

test-unit-fuzz:
	@echo "Running fuzz tests with coverage..."
	go test ./cache/... -fuzz=FuzzPayloadIDCacheBasic -fuzztime=${SHORT_FUZZ_TIME}
	go test ./cache/... -fuzz=FuzzPayloadIDInvalidInput -fuzztime=${SHORT_FUZZ_TIME}
	go test ./cache/... -fuzz=FuzzPayloadIDCacheConcurrency -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzSSZUint64Marshal ./types/consensus/primitives/... -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzSSZUint64Unmarshal ./types/consensus/primitives/... -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzHashTreeRoot ./crypto/sha256/... -fuzztime=${MEDIUM_FUZZ_TIME}
	go test -fuzz=FuzzQueueSimple ./lib/store/collections/... -fuzztime=${SHORT_FUZZ_TIME}
	go test -fuzz=FuzzQueueMulti ./lib/store/collections/... -fuzztime=${SHORT_FUZZ_TIME}

#################
#     forge     #
#################

test-forge-cover:
	@echo "Running forge test with coverage..."
	@cd $(CONTRACTS_DIR) && FOUNDRY_PROFILE=coverage forge coverage --report lcov --report-file ../test-forge-cover.txt

test-forge-fuzz:
	@echo "Running forge fuzz tests..."
	@cd $(CONTRACTS_DIR) && FOUNDRY_PROFILE=fuzz forge test --mt testFuzz

#################
#      e2e      #
#################

test-e2e:
	@$(MAKE) test-e2e-no-build

test-e2e-no-build:
	@echo "Running e2e tests..."


###############################################################################
###                                Linting                                  ###
###############################################################################

format:
	@$(MAKE) license-fix buf-lint-fix forge-lint-fix golangci-fix

lint:
	@$(MAKE) license buf-lint forge-lint golangci

#################
#     forge     #
#################

forge-lint-fix:
	@echo "--> Running forge fmt"
	@cd $(CONTRACTS_DIR) && forge fmt

forge-lint:
	@echo "--> Running forge lint"
	@cd $(CONTRACTS_DIR) && forge fmt --check

#################
# golangci-lint #
#################

golangci-install:
	@echo "--> Installing golangci-lint"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint

golangci:
	@$(MAKE) golangci-install
	@echo "--> Running linter"
	@go list -f '{{.Dir}}/...' -m | xargs golangci-lint run  --timeout=10m --concurrency 8 -v 

golangci-fix:
	@$(MAKE) golangci-install
	@echo "--> Running linter"
	@go list -f '{{.Dir}}/...' -m | xargs golangci-lint run  --timeout=10m --fix --concurrency 8 -v 


#################
#    license    #
#################

license-install:
	@echo "--> Installing google/addlicense"
	@go install github.com/google/addlicense

license:
	@$(MAKE) license-install
	@echo "--> Running addlicense with -check"
	@for module in $(MODULES); do \
		(cd $$module && addlicense -check -v -f ./LICENSE.header ./.) || exit 1; \
	done

license-fix:
	@$(MAKE) license-install
	@echo "--> Running addlicense"
	@for module in $(MODULES); do \
		(cd $$module && addlicense -v -f ./LICENSE.header ./.) || exit 1; \
	done


#################
#     gosec     #
#################

gosec-install:
	@echo "--> Installing gosec"
	@go install github.com/cosmos/gosec/v2/cmd/gosec 

gosec:
	@$(MAKE) gosec-install
	@echo "--> Running gosec"
	@gosec -exclude G702 ./...


#################
#     pkgsite     #
#################

pkgsite-install:
	@echo "--> Installing pkgsite"
	@go install golang.org/x/pkgsite/cmd/pkgsite
pkgsite:
	@$(MAKE) pkgsite-install
	@echo "Starting pkgsite server at http://localhost:6060/pkg/github.com/itsdevbear/bolaris/..."
	@pkgsite -http=:6060

#################
#    slither    #
#################

slither:
	docker run \
	-t \
	--platform linux/amd64 \
	-v ./contracts:/contracts \
	trailofbits/eth-security-toolbox \
	slither /contracts/src --config-file /contracts/slither.config.json


#################
#     proto     #
#################


protoImageName    := "ghcr.io/cosmos/proto-builder"
protoImageVersion := "0.14.0"
modulesProtoDir := "proto"

proto:
	@$(MAKE) buf-lint-fix buf-lint proto-build

proto-build:
	@docker run --rm -v ${CURRENT_DIR}:/workspace --workdir /workspace $(protoImageName):$(protoImageVersion) sh ./build/scripts/proto_generate.sh
	@./build/scripts/prysm_ssz_replacements.sh

proto-clean:
	@find . -name '*.pb.go' -delete
	@find . -name '*.pb.gw.go' -delete
	
buf-install:
	@echo "--> Installing buf"
	@go install github.com/bufbuild/buf/cmd/buf

buf-lint-fix:
	@$(MAKE) buf-install 
	@echo "--> Running buf format"
	@buf format -w --error-format=json $(modulesProtoDir)

buf-lint:
	@$(MAKE) buf-install 
	@echo "--> Running buf lint"
	@buf lint --error-format=json $(modulesProtoDir)

#################
#    sszgen    #
#################

sszgen-install:
	@echo "--> Installing sszgen"
	@go install github.com/prysmaticlabs/fastssz/sszgen

sszgen-clean:
	@find . -name '*.pb_encoding.go' -delete

sszgen:
	@$(MAKE) sszgen-install sszgen-clean
	@echo "--> Running sszgen on all structs with ssz tags"
	@sszgen -path ./types/consensus/v1/capella -objs BeaconKitBlockCapella,BlindedBeaconKitBlockCapella \
    --include ./types/consensus/primitives,./types/consensus/v1,\
	$(HOME)/go/pkg/mod/github.com/prysmaticlabs/prysm/v4@v4.2.1/proto/engine/v1
	@sszgen -path ./types/consensus/v1 -objs Deposit \
    --include $(HOME)/go/pkg/mod/github.com/prysmaticlabs/prysm/v4@v4.2.1/proto/engine/v1

##############################################################################
###                             Dependencies                                ###
###############################################################################

tidy: |
	go mod tidy

repo-rinse: |
	git clean -xfd
	git submodule foreach --recursive git clean -xfd
	git submodule foreach --recursive git reset --hard
	git submodule update --init --recursive


.PHONY: build build-linux-amd64 build-linux-arm64 \
	$(BUILD_TARGETS) clean \
	forge-build forge-clean proto proto-build docker-build generate \
	abigen-install mockery-install mockery \
	start test-unit test-unit-cover test-forge-cover test-forge-fuzz \
	test-e2e test-e2e-no-build hive-setup hive-view test-hive \
	test-hive-v test-localnet test-localnet-no-build format lint \
	forge-lint-fix forge-lint golangci-install golangci golangci-fix \
	license-install license license-fix \
	gosec-install gosec buf-install buf-lint-fix buf-lint sync tidy repo-rinse
