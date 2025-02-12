#!/usr/bin/make -f

repo-rinse: | ## dangerous!!! make sure you know what you are doing
	git clean -xfd
	git submodule foreach --recursive git clean -xfd
	git submodule foreach --recursive git reset --hard
	git submodule update --init --recursive

# tidy-sync-check runs `go mod tidy` and checks if go.mod/go.sum files are in sync with
# the dependencies in the codebase. If they are not in sync, it will exit with a
# non-zero status code.
tidy-sync-check:
	@{ \
	pre_tidy_diff=$$(git diff --ignore-space-change go.mod go.sum); \
	go mod tidy; \
	post_tidy_diff=$$(git diff --ignore-space-change go.mod go.sum); \
	if [ "$$pre_tidy_diff" != "$$post_tidy_diff" ]; then \
		echo "go.mod or go.sum changed after running 'go mod tidy'"; \
		git diff go.mod go.sum; \
		exit 1; \
	fi; \
	}
