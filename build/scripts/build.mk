#!/usr/bin/make -f

# cosmos.mk contains the build configuration for this project as
# forked from the cosmos-sdk https://github.com/cosmos/cosmos-sdk/tree/main/Makefile
# It is moved here to reduce clutter in the main Makefile and to make it easier to
# update.

export VERSION := $(shell echo $(shell git describe --tags --always --match "v*") | sed 's/^v//')
export COMMIT := $(shell git log -1 --format='%H')
CURRENT_DIR = $(shell pwd)
OUT_DIR ?= $(CURDIR)/build/bin
BINDIR ?= $(GOPATH)/build/bin
TESTAPP_DIR = beacond
TESTAPP_CMD_DIR = $(TESTAPP_DIR)/cmd
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

# process build tags
build_tags = netgo

ifeq (legacy,$(findstring legacy,$(COSMOS_BUILD_OPTIONS)))
  build_tags += app_v1
endif

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  build_tags += rocksdb grocksdb_clean_link
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

# always include pebble
build_tags += pebbledb

# always include blst
build_tags += blst

# always include ckzg
build_tags += ckzg
build_tags += cgo

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sim \
		-X github.com/cosmos/cosmos-sdk/version.AppName=simd \
		-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
# build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif