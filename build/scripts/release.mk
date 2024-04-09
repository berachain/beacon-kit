#!/usr/bin/make -f

RELEASE_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-arm64

build-release: $(RELEASE_TARGETS)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 \
	cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && \
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/beacond-linux-amd64  ./.

build-linux-arm64:
	GOOS=linux GOARCH=arm64 \
	cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && \
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/beacond-linux-arm64  ./.

build-darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
	cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && \
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/beacond-darwin-arm64  ./.
