#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --tags --always --match "v*") | sed 's/^v//')
export COMMIT := $(shell git log -1 --format='%H')
CURRENT_DIR = $(shell pwd)
OUT_DIR ?= $(CURDIR)/build/bin
BINDIR ?= $(GOPATH)/build/bin
TESTAPP_DIR = beacond
TESTAPP_FILES_DIR = testing/files
TESTAPP_CMD_DIR = $(TESTAPP_DIR)/cmd
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

# process build tags
BUILD_TAGS = netgo

ifeq (legacy,$(findstring legacy,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += app_v1
endif

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS += rocksdb grocksdb_clean_link
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += boltdb
endif

# always include pebble
BUILD_TAGS += pebbledb

# always include blst
BUILD_TAGS += blst
BUILD_TAGS += bls12381

# always include ckzg
BUILD_TAGS += ckzg
BUILD_TAGS += cgo

whitespace :=
whitespace += $(whitespace)
comma := ,
BUILD_TAGS_comma_sep := $(subst $(whitespace),$(comma),$(BUILD_TAGS))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sim \
		-X github.com/cosmos/cosmos-sdk/version.AppName=simd \
		-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS_comma_sep)"

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath 
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

# This allows us to reuse the build target steps for both go build and go install
BUILD_TARGETS := build install

## Build: 
build: BUILD_ARGS=-o $(OUT_DIR)/beacond ## build `beacond`

$(BUILD_TARGETS): $(OUT_DIR)/
	@echo "Building ${TESTAPP_CMD_DIR}"
	@cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./.

$(OUT_DIR)/:
	mkdir -p $(OUT_DIR)/

	# Variables
ARCH ?= $(shell uname -m)
ifeq ($(ARCH),)
	ARCH = arm64
endif
IMAGE_NAME ?= beacond

# Docker Paths
DOCKERFILE = ./Dockerfile

build-docker: ## build a docker image containing `beacond`
	@echo "Build a release docker image for the Cosmos SDK chain..."
	docker build \
	--platform linux/$(ARCH) \
	--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
	--build-arg GIT_VERSION=$(shell git describe --tags --always --dirty) \
	--build-arg GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	--build-arg GOOS=linux \
	--build-arg GOARCH=$(ARCH) \
	--build-arg BUILD_TAGS="$(BUILD_TAGS)" \
	-f ${DOCKERFILE} \
	-t $(IMAGE_NAME):$(VERSION) \
	.
