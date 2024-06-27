#!/usr/bin/make -f

## Linting:
format: ## run all configured formatters
	@$(MAKE) license-fix forge-lint-fix golines golangci-fix star-fix

lint: ## run all configured linters
	@$(MAKE) license markdownlint forge-lint golangci star-lint


#################
# golangci-lint #
#################

golangci:
	@echo "--> Running linter on all modules"
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		printf "[%d/%d modules complete] Running linter in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --concurrency 8) || exit 1; \
		count=$$((count + 1)); \
	done
	@printf "All modules complete\n"
	
golangci-fix:
	@echo "--> Running linter with fixes on all modules"
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		printf "[%d/%d modules complete] Running formatter in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config $(ROOT_DIR)/.golangci.yaml --timeout=10m --fix --concurrency 8) || exit 1; \
		count=$$((count + 1)); \
	done
	@printf "All modules complete\n"

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
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		printf "[%d/%d modules complete] Checking licenses in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/google/addlicense -check -v -f $(ROOT_DIR)/LICENSE.header ./. ) || exit 1; \
		count=$$((count + 1)); \
	done
	@printf "License check complete for all modules\n"

license-fix:
	@echo "--> Running addlicense"
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		printf "[%d/%d modules complete] Applying licenses in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run github.com/google/addlicense -v -f $(ROOT_DIR)/LICENSE.header ./. ) || exit 1; \
		count=$$((count + 1)); \
	done
	@printf "License application complete for all modules\n"


#################
#    nilaway    #
#################

nilaway:
	@echo "--> Running nilaway"
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	total=$$(echo "$$dirs" | wc -l); \
	count=0; \
	for dir in $$dirs; do \
		count=$$((count + 1)); \
		printf "[%d/%d modules complete] Running nilaway in %s\n" $$count $$total $$dir && \
		(cd $$dir && go run go.uber.org/nilaway/cmd/nilaway -exclude-errors-in-files "pkg/components/module,pkg/deposit" -v ./...) || exit 1; \
	done
	@printf "Nilaway complete for all modules\n"

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
