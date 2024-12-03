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
	(go run github.com/google/addlicense -check -v -f $(ROOT_DIR)/LICENSE.header cmd build kurtosis mod testing) || exit 1;
	@printf "License check complete\n"

license-fix:
	echo "--> Running addlicense"
	(go run github.com/google/addlicense -v -f $(ROOT_DIR)/LICENSE.header cmd build kurtosis mod testing) || exit 1;
	@printf "License check complete\n"

#################
#    nilaway    #
#################

nilaway:
	@echo "--> Running nilaway"
	(go run go.uber.org/nilaway/cmd/nilaway -exclude-errors-in-files "mod/geth-primitives/pkg/deposit/" -v ./...) || exit 1;
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
	docker run \
	-t \
	--platform linux/amd64 \
	-v ./contracts:/contracts \
	trailofbits/eth-security-toolbox:edge \
	/bin/bash -c "cd /contracts && slither ./."

#################
# markdown-lint #
#################

markdownlint:
	@echo "--> Running markdownlint"
	@docker run --rm -v $(ROOT_DIR):/workspace -w /workspace -t markdownlint/markdownlint:latest --git-recurse **/**.md
