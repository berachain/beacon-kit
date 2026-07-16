#!/usr/bin/make -f

# Tests a coordinated network upgrade end to end: a local multi-validator devnet halts itself at a configured
# halt point, the beacond binary is swapped, and the chain must resume from the same data directories. Only
# beacond is swapped, each validator's bera-reth execution client keeps running throughout. The test flow is
# documented in testing/upgrade/halt-swap-resume-test.sh.
#
#   make test-halt-swap-resume        halt at a fixed height (--halt-height)
#   make test-halt-swap-resume-time   halt at a wall-clock time (--halt-time)
#
# Both targets run the local build as both the old and the new beacond, gating the halt/swap/restart mechanics
# rather than cross-version compatibility. Extra script flags pass through HALT_SWAP_RESUME_ARGS, e.g.
# HALT_SWAP_RESUME_ARGS=--no-load to skip the cast-driven tx load.

# Extra arguments passed through to halt-swap-resume-test.sh, e.g. HALT_SWAP_RESUME_ARGS="--keep --num-vals 7".
HALT_SWAP_RESUME_ARGS ?=

# The bera-reth release used as each validator's execution client.
BERA_RETH_VERSION ?= v1.4.2
BERA_RETH_BIN_DIR = /tmp/.halt-upgrade-test
BERA_RETH_PLATFORM := $(shell uname -m | sed 's/arm64/aarch64/')-$(if $(filter Darwin,$(shell uname)),apple-darwin,unknown-linux-gnu)
BERA_RETH_ASSET = bera-reth-$(BERA_RETH_VERSION)-$(BERA_RETH_PLATFORM).tar.gz
BERA_RETH_URL = https://github.com/berachain/bera-reth/releases/download/$(BERA_RETH_VERSION)/$(BERA_RETH_ASSET)
BERA_RETH_BIN = $(BERA_RETH_BIN_DIR)/bera-reth-$(BERA_RETH_VERSION)

# The tarball comes from the GitHub release over HTTPS and is trusted as-is. The version-stamped mv at the end
# is the commit step, so an interrupted download never masquerades as a valid cached binary.
$(BERA_RETH_BIN):
	@mkdir -p $(BERA_RETH_BIN_DIR)
	curl -fsSL -o $(BERA_RETH_BIN_DIR)/$(BERA_RETH_ASSET) $(BERA_RETH_URL)
	@tar xzf $(BERA_RETH_BIN_DIR)/$(BERA_RETH_ASSET) -C $(BERA_RETH_BIN_DIR)
	@mv $(BERA_RETH_BIN_DIR)/bera-reth $@ && chmod +x $@
	@rm -f $(BERA_RETH_BIN_DIR)/$(BERA_RETH_ASSET)

test-halt-swap-resume: build $(BERA_RETH_BIN) ## halt the devnet at a fixed height, swap beacond, verify it resumes
	@OLD_BIN=$(CURDIR)/build/bin/beacond \
	NEW_BIN=$(CURDIR)/build/bin/beacond \
	RETH_BIN=$(BERA_RETH_BIN) \
	./testing/upgrade/halt-swap-resume-test.sh $(HALT_SWAP_RESUME_ARGS)

test-halt-swap-resume-time: build $(BERA_RETH_BIN) ## halt on halt-time instead of halt-height, swap, resume
	@OLD_BIN=$(CURDIR)/build/bin/beacond \
	NEW_BIN=$(CURDIR)/build/bin/beacond \
	RETH_BIN=$(BERA_RETH_BIN) \
	./testing/upgrade/halt-swap-resume-test.sh --halt-time-offset 25 $(HALT_SWAP_RESUME_ARGS)

.PHONY: test-halt-swap-resume test-halt-swap-resume-time
