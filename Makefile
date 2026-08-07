APP := adobe-fonts-dump
VERSION ?= 0.1.0
BUILD_DIR := build
OUTPUT ?= $(BUILD_DIR)/$(APP)
BINARY ?= $(OUTPUT)
APP_BUNDLE ?= $(BUILD_DIR)/AdobeFontsDump.app
DMG ?= $(BUILD_DIR)/$(APP)-macos-arm64.dmg
GO ?= go
GO_PACKAGE ?= ./src

.PHONY: build run package-macos clean

build:
	mkdir -p "$(dir $(OUTPUT))"
	$(GO) build -trimpath -ldflags="-s -w" -o "$(OUTPUT)" $(GO_PACKAGE)

run:
	$(GO) run $(GO_PACKAGE)

package-macos:
	VERSION="$(VERSION)" ./scripts/package-macos.sh "$(BINARY)" "$(APP_BUNDLE)" "$(DMG)"

clean:
	rm -rf "$(BUILD_DIR)"
