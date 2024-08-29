#!/usr/bin/make -f
MODULES := $(shell find . -type f -name 'go.mod' -exec dirname {} \;)
# Exclude root module
MODULES := $(filter-out ./,$(MODULES))

CONTRACTS_DIR := ./contracts
EXAMPLES_DIR = examples
OUT_DIR ?= $(CURDIR)/build/bin
BINDIR ?= $(GOPATH)/build/bin
TESTNAME = beacon
TESTAPP_DIR = beacond
TESTAPP_FILES_DIR = testing/files
TESTAPP_CMD_DIR = $(TESTAPP_DIR)/cmd
