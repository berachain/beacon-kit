#!/usr/bin/make -f

RELEASE_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-arm64

define build_release
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(3) \
	cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && \
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/beacond-$(1)-$(2) ./.
endef

build-linux-amd64:
	$(call build_release,linux,amd64,0)

build-linux-arm64:
	$(call build_release,linux,arm64,0)

build-darwin-arm64:
	$(call build_release,darwin,arm64,1)

