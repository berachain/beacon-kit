#!/usr/bin/make -f

## Linting:
format: ## run all configured formatters
	@$(MAKE) license-fix forge-lint-fix golines golangci-fix star-fix

lint: ## run all configured linters
	@$(MAKE) license markdownlint forge-lint golangci star-lint


#################
# golangci-lint #
#################

# TODO: Remove GODEBUG override once: https://github.com/golang/go/issues/68877 is resolved.
golangci:
	@echo "--> Running linter on all modules"
	(GODEBUG=gotypesalias=0 go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --concurrency 8) || exit 1;
	@printf "All modules complete\n"


# TODO: Remove GODEBUG override once: https://github.com/golang/go/issues/68877 is resolved.
golangci-fix:
	@echo "--> Running linter with fixes on all modules"
	(GODEBUG=gotypesalias=0 go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --fix --concurrency 8) || exit 1;
	@printf "All modules complete\n"

#################
#    golines    #
#################

golines:
	@echo "--> Running golines"
	@./scripts/build/golines.sh

#################
#    license    #
#################

license:
	@echo "--> Running addlicense with -check"
	(go run github.com/google/addlicense -check -v -f $(ROOT_DIR)/LICENSE.header -ignore "contracts/**" .) || exit 1;
	@printf "License check complete\n"

license-fix:
	echo "--> Running addlicense"
	(go run github.com/google/addlicense -v -f $(ROOT_DIR)/LICENSE.header -ignore "contracts/**" .) || exit 1;
	@printf "License check complete\n"

#################
#    nilaway    #
#################

nilaway:
	@echo "--> Running nilaway"
	(go run go.uber.org/nilaway/cmd/nilaway -exclude-errors-in-files "geth-primitives/deposit/" -v ./...) || exit 1;
	@printf "Nilaway check complete\n"

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
	chmod -R o+rw ./contracts
	docker run \
	-t \
	--platform linux/amd64 \
	-v ./contracts:/contracts:rw \
	trailofbits/eth-security-toolbox:edge \
	/bin/bash -c "cd /contracts/ && slither ./."

#################
# markdown-lint #
#################

markdownlint:
	@echo "--> Running markdownlint"
	@docker run --rm -v $(ROOT_DIR):/workspace -w /workspace -t markdownlint/markdownlint:latest --git-recurse **/**.md

#################
# all ci linters #
#################

lint-ci: lint slither gosec nilaway markdownlintg generate-check \
    tidy-sync-check test-unit-cover test-unit-bench test-unit-fuzz \
	test-forge-cover test-forge-fuzz
