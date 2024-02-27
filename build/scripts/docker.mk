#!/usr/bin/make -f

###############################################################################
###                                 Docker                                  ###
###############################################################################

# Variables
DOCKER_TYPE ?= base
ARCH ?= arm64
GO_VERSION ?= 1.22.0
IMAGE_NAME ?= beacond
IMAGE_VERSION ?= v0.0.0
BASE_IMAGE ?= beacond/base:$(IMAGE_VERSION)

# Docker Paths
BASE_DOCKER_PATH = ./examples/beacond
DOCKERFILE = ./examples/beacond/Dockerfile
EXEC_DOCKER_PATH = $(BASE_DOCKER_PATH)/Dockerfile
LOCAL_DOCKER_PATH = $(BASE_DOCKER_PATH)/local/Dockerfile
SEED_DOCKER_PATH =  $(BASE_DOCKER_PATH)/seed/Dockerfile
VAL_DOCKER_PATH =  $(BASE_DOCKER_PATH)/validator/Dockerfile
LOCALNET_CLIENT_PATH = ./e2e/precompile/beacond
LOCALNET_DOCKER_PATH = $(LOCALNET_CLIENT_PATH)/Dockerfile

# Image Build
docker-build:
	@echo "Build a release docker image for the Cosmos SDK chain..."
	@$(MAKE) docker-build-$(DOCKER_TYPE)

# Docker Build Types
docker-build-base:
	$(call docker-build-helper,$(EXEC_DOCKER_PATH),base)

docker-build-local:
	$(call docker-build-helper,$(LOCAL_DOCKER_PATH),local,--build-arg BASE_IMAGE=$(BASE_IMAGE))

docker-build-seed:
	$(call docker-build-helper,$(SEED_DOCKER_PATH),seed,--build-arg BASE_IMAGE=$(BASE_IMAGE))

docker-build-validator:
	$(call docker-build-helper,$(VAL_DOCKER_PATH),validator,--build-arg BASE_IMAGE=$(BASE_IMAGE))

docker-build-localnet:
	$(call docker-build-helper,$(LOCALNET_DOCKER_PATH),localnet,--build-arg BASE_IMAGE=$(BASE_IMAGE))

# Docker Build Function
define docker-build-helper
	docker build \
	--build-arg GO_VERSION=$(GO_VERSION) \
	--platform linux/$(ARCH) \
	--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
	--build-arg GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	--build-arg GOOS=linux \
	--build-arg GOARCH=$(ARCH) \
	-f $(1) \
	-t $(IMAGE_NAME)/$(2):$(IMAGE_VERSION) \
	$(if $(3),$(3)) \
	.

endef

.PHONY: docker-build-localnet