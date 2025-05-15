#!/usr/bin/make -f

## Linting:
format: ## run all configured formatters
	@$(MAKE) license-fix forge-lint-fix golines golangci-fix star-fix

lint: ## run all configured linters
	@$(MAKE) license markdownlint forge-lint golangci star-lint


#################
# golangci-lint #
#################

golangci-install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# TODO: Remove GODEBUG override once: https://github.com/golang/go/issues/68877 is resolved.
golangci: golangci-install
	@echo "--> Running linter on all modules"
	(GODEBUG=gotypesalias=0 golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --concurrency 8) || exit 1;
	@printf "All modules complete\n"


# TODO: Remove GODEBUG override once: https://github.com/golang/go/issues/68877 is resolved.
golangci-fix: golangci-install
	@echo "--> Running linter with fixes on all modules"
	(GODEBUG=gotypesalias=0 golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --fix --concurrency 8) || exit 1;
	@printf "All modules complete\n"

#################
#    golines    #
#################

golines-install:
	@go install github.com/segmentio/golines@latest

golines: golines-install
	@echo "--> Running golines"
	@./scripts/build/golines.sh

#################
#    license    #
#################

license-install:
	@go install github.com/google/addlicense@latest

license: license-install
	@echo "--> Running addlicense with -check"
	(addlicense -check -v -f $(ROOT_DIR)/LICENSE.header -ignore "contracts/**" -ignore ".idea/**" .) || exit 1;
	@printf "License check complete\n"

license-fix: license-install
	echo "--> Running addlicense"
	(addlicense -v -f $(ROOT_DIR)/LICENSE.header -ignore "contracts/**" -ignore ".idea/**" .) || exit 1;
	@printf "License check complete\n"

#################
#    nilaway    #
#################

nilaway-install:
	@go install go.uber.org/nilaway/cmd/nilaway@latest

nilaway: nilaway-install
	@echo "--> Running nilaway"
	(nilaway -test=false -exclude-errors-in-files "geth-primitives/deposit/","geth-primitives/ssztest/" -v ./...) || exit 1;
	@printf "Nilaway check complete\n"

#################
#     gosec     #
#################

gosec-install:
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

gosec: gosec-install
	@echo "--> Running gosec"
	@gosec ./...

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

lint-ci: lint slither gosec nilaway markdownlint generate-check \
    tidy-sync-check test-unit-cover test-unit-bench test-unit-fuzz \
	test-forge-cover test-forge-fuzz
