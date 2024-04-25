#!/usr/bin/make -f

## Linting:
format: ## run all configured formatters
	@$(MAKE) license-fix buf-lint-fix forge-lint-fix golines golangci-fix star-fix

lint: ## run all configured linters
	@$(MAKE) license buf-lint forge-lint golangci star-lint


#################
# golangci-lint #
#################

golangci:
	@echo "--> Running linter on all modules"
	@find . -name 'go.mod' -execdir sh -c 'echo "Linting in $$(pwd)" && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --concurrency 8 -v' \;

golangci-fix:
	@echo "--> Running linter with fixes on all modules"
	@find . -name 'go.mod' -execdir sh -c 'echo "Applying fixes in $$(pwd)" && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --fix --concurrency 8 -v' \;

#################
#    golines    #
#################

golines:
	@echo "--> Running golines"
	@./build/scripts/golines.sh

#################
#    license    #
#################

license: 
	@echo "--> Running addlicense with -check"
	@find . -name 'go.mod' -execdir go run github.com/google/addlicense -check -v -f $(ROOT_DIR)/LICENSE.header ./. \;

license-fix:
	@echo "--> Running addlicense"
	@find . -name 'go.mod' -execdir go run github.com/google/addlicense -v -f $(ROOT_DIR)/LICENSE.header ./. \;


#################
#    nilaway    #
#################

nilaway:
	@echo "--> Running nilaway"
	@find . -name 'go.mod' -execdir sh -c 'go run go.uber.org/nilaway/cmd/nilaway \
		-exclude-errors-in-files "x/beacon/api,runtime/services/staking/abi" \
		-v ./...' \;

#################
#     gosec     #
#################

gosec:
	@echo "--> Running gosec"
	@go run github.com/cosmos/gosec/v2/cmd/gosec -exclude G702 ./...

#################
#    slither    #
#################

slither:
	docker run \
	-t \
	--platform linux/amd64 \
	-v ./contracts:/contracts \
	trailofbits/eth-security-toolbox:edge \
	/bin/bash -c "cd /contracts && slither ./src/eip4788 && slither ./src/staking"
