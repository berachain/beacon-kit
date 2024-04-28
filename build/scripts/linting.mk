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
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		count=$$((count + 1)); \
		printf "[%d/%d modules complete] Linting in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --concurrency 8 -v) || exit 1; \
	done
	
golangci-fix:
	@echo "--> Running linter with fixes on all modules"
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		count=$$((count + 1)); \
		printf "[%d/%d modules complete] Applying fixes in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --fix --concurrency 8 -v) || exit 1; \
	done

define_run_linter = \
	@echo "--> Running $(1) on all modules"; \
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		count=$$((count + 1)); \
		printf "[%d/%d modules complete] $(2) in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m $(3) --concurrency 8 -v) || exit 1; \
	done

$(eval $(call define_run_linter,linter,Linting,))
$(eval $(call define_run_linter,linter with fixes,Applying fixes,--fix))


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
