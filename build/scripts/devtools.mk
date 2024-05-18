#!/usr/bin/make -f
## DevTools:

bet: build format test-unit ## yo bet
	@git add .
	@git commit -m 'bet'
	@git push

checkpoint: format  ## checkpoint and push to remote
	@git add -A
	@git commit -m "yo bet"
	@git push
	@echo "checkpointed and pushed to remote"

repo-rinse: | ## dangerous!!! make sure you know what you are doing
	git clean -xfd
	git submodule foreach --recursive git clean -xfd
	git submodule foreach --recursive git reset --hard
	git submodule update --init --recursive

sync: 
	@echo "Running go mod download && go work sync"
	@go mod download
	@go work sync

tidy: ## run go mod tidy in all modules
	@echo "Running go mod tidy in all modules"
	@find . -name 'go.mod' ! -path './go.mod' -execdir go mod tidy \;

yap: ## the yap cave
	@go run ./mod/node-builder/pkg/utils/yap/yap.go

tidy-sync-check:
	@$(MAKE) tidy sync repo-rinse
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Tidy and sync operations resulted in changes"; \
		git status -s; \
		git diff --exit-code; \
	fi


.PHONY: format build test-unit bet

