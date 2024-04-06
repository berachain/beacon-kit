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
	@echo "--> Running linter"
	@go list -f '{{.Dir}}/...' -m | grep -v '**/contracts' | \
		xargs go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout=10m --concurrency 8 -v 

golangci-fix:
	@echo "--> Running linter"
	@go list -f '{{.Dir}}/...' -m | grep -v '**/contracts' | \
		xargs go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout=10m --fix --concurrency 8 -v 

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
	@go run go.uber.org/nilaway/cmd/nilaway \
		-exclude-errors-in-files beacond/x/beacon/api,mod/runtime/services/staking/abi \
		-v ./...

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
