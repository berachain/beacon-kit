#!usr/bin/make -f 

GETH_GO_GENERATE_VERSION := $(shell grep 'github.com/ethereum/go-ethereum' beacond/go.mod | awk '{print $$2}')
GOPATH = $(shell go env GOPATH)
GETH_PKG_INCLUDE := $(GOPATH)/pkg/mod/github.com/ethereum/go-ethereum@$(GETH_GO_GENERATE_VERSION)

## Codegen:
generate: ## generate all the code
	@$(MAKE) forge-build
	@for module in $(MODULES); do \
		echo "Running go generate in $$module"; \
		(cd $$module && \
			GETH_PKG_INCLUDE=$(GETH_PKG_INCLUDE) go generate ./...) || exit 1; \
	done
	@go run github.com/vektra/mockery/v2@v2.49.0

generate-check:
	@$(MAKE) forge-build
	@$(MAKE) generate
	@if [ -n "$$(git status --porcelain | grep -vE '\.ssz\.go$$')" ]; then \
		echo "Generated files are not up to date"; \
		git status -s | grep -vE '\.ssz\.go$$'; \
		git diff -- . ':(exclude)*.ssz.go'; \
		exit 1; \
	fi
