#!/usr/bin/make -f

protoImageName    := "ghcr.io/cosmos/proto-builder"
protoImageVersion := "0.14.0"
modulesProtoDir := "mod/node-core/pkg/components/module/proto"

## Protobuf:
proto: ## run all the proto tasks
	@$(MAKE) proto-build

proto-build: ## build the proto files
	@docker run --rm -v ${CURRENT_DIR}:/workspace --workdir /workspace $(protoImageName):$(protoImageVersion) sh ./scripts/build/proto_generate_pulsar.sh

proto-clean: ## clean the proto files
	@find . -name '*.pb.go' -delete
	@find . -name '*.pb.gw.go' -delete

buf-install:
	@echo "--> Installing buf"
	@go install github.com/bufbuild/buf/cmd/buf