#!/usr/bin/make -f

RELEASE_TARGETS := build-linux-amd64 build-linux-arm64 build-darwin-arm64

define build_release
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(3) \
	cd ${CURRENT_DIR}/$(TESTAPP_CMD_DIR) && \
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/beacond-$(1)-$(2)-$(3) ./.
endef

build-linux-amd64:
	$(call build_release,linux,amd64,0,$(shell git describe --tags --always --dirty))
build-linux-amd64-$(VERSION):
	$(call build_release,linux,amd64,0,$(VERSION))

build-linux-arm64:
	$(call build_release,linux,arm64,0,$(shell git describe --tags --always --dirty))
build-linux-arm64-$(VERSION):
	$(call build_release,linux,arm64,0,$(VERSION))

build-darwin-arm64:
	$(call build_release,darwin,arm64,1,$(shell git describe --tags --always --dirty))
build-darwin-arm64-$(VERSION):
	$(call build_release,darwin,arm64,1,$(VERSION))
