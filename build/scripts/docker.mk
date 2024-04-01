#!/usr/bin/make -f

###############################################################################
###                                 Docker                                  ###
###############################################################################

# Variables
ARCH ?= $(shell uname -m)
ifeq ($(ARCH),)
	ARCH = arm64
endif
GO_VERSION ?= 1.22.1
IMAGE_NAME ?= beacond

# Docker Paths
DOCKERFILE = ./Dockerfile

# Image Build
docker-build:
	@echo "Build a release docker image for the Cosmos SDK chain..."
	docker build \
	--build-arg GO_VERSION=$(GO_VERSION) \
	--platform linux/$(ARCH) \
	--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
	--build-arg GIT_VERSION=$(shell git describe --tags --always --dirty) \
	--build-arg GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	--build-arg GOOS=linux \
	--build-arg GOARCH=$(ARCH) \
	-f ${DOCKERFILE} \
	-t $(IMAGE_NAME):$(VERSION) \
	.