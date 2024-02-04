#!/usr/bin/make -f

###############################################################################
###                                 Docker                                  ###
###############################################################################

# Variables
ARCH ?= arm64
GO_VERSION ?= 1.21.6
IMAGE_NAME ?= beacond

# Docker Paths
DOCKERFILE = ./examples/beacond/Dockerfile

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