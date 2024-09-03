#!/usr/bin/make -f

#################
#     forge     #
#################

forge-build: |
	@forge build --extra-output-files bin --extra-output-files abi  --root $(CONTRACTS_DIR)

forge-clean: |
	@forge clean --root $(CONTRACTS_DIR)

test-forge-unit:
	@echo "Running forge unit tests..."
	@cd $(CONTRACTS_DIR) && FOUNDRY_PROFILE=coverage forge test 

test-forge-cover:
	@echo "Running forge test with coverage..."
	@cd $(CONTRACTS_DIR) && FOUNDRY_PROFILE=coverage forge build && forge coverage --nmt testFuzz --report lcov --report-file ../test-forge-cover.txt

test-forge-fuzz:
	@echo "Running forge fuzz tests..."
	@cd $(CONTRACTS_DIR) && FOUNDRY_PROFILE=fuzz forge test

forge-snapshot:
	@echo "Running forge snapshot..."
	@cd $(CONTRACTS_DIR) && forge snapshot --isolate --fuzz-runs 1

forge-snapshot-diff:
	@echo "Running forge snapshot diff..."
	@cd $(CONTRACTS_DIR) && forge snapshot --diff --isolate --fuzz-runs 1

forge-lint-fix:
	@echo "--> Running forge fmt"
	@cd $(CONTRACTS_DIR) && forge fmt .

forge-lint:
	@echo "--> Running forge lint"
	@cd $(CONTRACTS_DIR) && forge fmt --check