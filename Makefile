#!/usr/bin/make -f

include scripts/build/build.mk
include scripts/build/codegen.mk
include scripts/build/constants.mk
include scripts/build/devtools.mk
include scripts/build/linting.mk
include scripts/build/protobuf.mk
include scripts/build/release.mk
include scripts/build/testing.mk
include contracts/Makefile
include kurtosis/Makefile
include scripts/build/help.mk
include testing/forge-script/Makefile

# Specify the default target if none is provided
.DEFAULT_GOAL := build
ROOT_DIR := $(shell pwd)


##############################################################################
###                             Dependencies                                ###
###############################################################################

.PHONY: clean format lint \
	buf-install proto-clean \
	test-unit test-unit-cover test-forge-cover test-forge-fuzz \
	forge-snapshot forge-snapshot-diff \
	test-e2e test-e2e-no-build \
	forge-lint-fix forge-lint golangci-install golangci golangci-fix \
	license license-fix \
	gosec golines tidy repo-rinse proto build


