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
	@go env -w GOPRIVATE=github.com/berachain
	@dirs=$$(find . -name 'go.mod' -exec dirname {} \;); \
	count=0; \
	total=$$(echo "$$dirs" | wc -l); \
	for dir in $$dirs; do \
		printf "[%d/%d modules complete] Running in %s\n" $$count $$total $$dir && \
		(cd $$dir && go mod tidy) || exit 1; \
		count=$$((count + 1)); \
	done
	@printf "go mod tidy complete for all modules\n"

update-dep: ## update a dependency in all of the go.mod files which import it
	@read -p "Enter go module path of the dependency (with an optional version specified by @): " dependency; \
	IFS='@' read -ra dependency_mod <<< "$$dependency"; \
	dependency_mod=$${dependency_mod[0]}; \
	for modfile in $$(find . -name go.mod); do \
		if grep -q $$dependency_mod $$modfile; then \
			echo "Updating $$modfile"; \
			DIR=$$(dirname $$modfile); \
			if [[ "$$dependency_mod" != *"$$(basename $$DIR)"* ]]; then \
				(cd $$DIR; go get -u $$dependency); \
			else \
				echo "Skipping $$DIR"; \
			fi; \
		else \
			echo "Skipping $$modfile"; \
		fi; \
	done

yap: ## the yap cave
	@go run ./mod/cli/pkg/utils/yap/yap.go

tidy-sync-check:
	@{ \
	pre_tidy_diff=$$(git diff --ignore-space-change); \
	$(MAKE) repo-rinse tidy sync; \
	post_tidy_diff=$$(git diff --ignore-space-change); \
	echo "$$pre_tidy_diff" > pre_tidy.diff; \
	echo "$$post_tidy_diff" > post_tidy.diff; \
	cmp -s pre_tidy.diff post_tidy.diff; \
	diff_status=$$?; \
	if [ $$diff_status -ne 0 ]; then \
		echo "Tidy and sync operations resulted in changes"; \
		diff pre_tidy.diff post_tidy.diff; \
	fi; \
	rm -f pre_tidy.diff post_tidy.diff; \
	exit $$diff_status; \
	}

.PHONY: format build test-unit bet
