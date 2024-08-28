#!/usr/bin/make -f

BERACHAIN_CONTRACTS_DIR = $(EXAMPLES_DIR)/berachain/contracts

#################
#     forge     #
#################

forge-build-berachain: |
	@forge build --extra-output-files bin --extra-output-files abi  --root $(BERACHAIN_CONTRACTS_DIR)

forge-clean-berachain: |
	@forge clean --root $(BERACHAIN_CONTRACTS_DIR)

test-forge-unit-berachain:
	@echo "Running forge unit tests..."
	@cd $(BERACHAIN_CONTRACTS_DIR) && FOUNDRY_PROFILE=coverage forge test 

test-forge-cover-berachain:
	@echo "Running forge test with coverage..."
	@cd $(BERACHAIN_CONTRACTS_DIR) && FOUNDRY_PROFILE=coverage forge build && \
	forge coverage --nmt testFuzz --report lcov --report-file ../test-forge-cover.txt

test-forge-fuzz-berachain:
	@echo "Running forge fuzz tests..."
	@cd $(BERACHAIN_CONTRACTS_DIR) && FOUNDRY_PROFILE=fuzz forge test

forge-snapshot-berachain:
	@echo "Running forge snapshot..."
	@cd $(BERACHAIN_CONTRACTS_DIR) && forge snapshot --isolate --fuzz-runs 1

forge-snapshot-diff-berachain:
	@echo "Running forge snapshot diff..."
	@cd $(BERACHAIN_CONTRACTS_DIR) && forge snapshot --diff --isolate --fuzz-runs 1

forge-lint-fix-berachain:
	@echo "--> Running forge fmt"
	@cd $(BERACHAIN_CONTRACTS_DIR) && forge fmt

forge-lint-berachain:
	@echo "--> Running forge lint"
	@cd $(BERACHAIN_CONTRACTS_DIR) && forge fmt --check