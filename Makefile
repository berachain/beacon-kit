#!/usr/bin/make -f

include build/scripts/build.mk 
include build/scripts/codegen.mk
include build/scripts/constants.mk
include build/scripts/devtools.mk 
include build/scripts/linting.mk
include build/scripts/protobuf.mk
include build/scripts/release.mk
include build/scripts/testing.mk
include contracts/Makefile 
include kurtosis/Makefile
include build/scripts/help.mk
include testing/scripts/fullnode.mk

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


