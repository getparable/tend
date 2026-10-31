# Copyright (c) 2026 Parable, LLC
# go.mod's toolchain directive does not downgrade a newer host Go. Select the
# supported toolchain for every target to satisfy the exact-version gate,
# including recursive make, without depending on the caller's environment.
export GOTOOLCHAIN := go1.26.5
export MACOSX_DEPLOYMENT_TARGET := 26.0
export CGO_ENABLED := 1
export CGO_CFLAGS := -mmacosx-version-min=26.0
export CGO_LDFLAGS := -mmacosx-version-min=26.0
export CODE_SIGN_IDENTITY

.PHONY: check app install-dev verify-launchd verify-lifecycle format
check:
	@arch="$$(uname -m)"; test "$$arch" = arm64 || { echo "Tend requires Apple Silicon (arm64); found $$arch" >&2; exit 1; }
	@version="$$(go env GOVERSION)" || exit 1; \
		echo "Go: $$version (GOTOOLCHAIN=$(GOTOOLCHAIN))"; \
		test "$$version" = go1.26.5 || { echo "Tend requires Go 1.26.5; remove the GOTOOLCHAIN make override" >&2; exit 1; }
	@version="$$(swift --version)" || exit 1; echo "$$version"; \
		case "$$version" in *"Apple Swift version 6.4 "*|*"Apple Swift version 6.4."*) ;; \
		*) echo "Tend requires Apple Swift 6.4; select matching tools with DEVELOPER_DIR" >&2; exit 1 ;; esac
	@version="$$(xcrun --show-sdk-version)" || exit 1; \
		echo "macOS SDK: $$version ($$(xcrun --show-sdk-path))"; \
		test "$$version" = 27.0 || { echo "Tend requires macOS SDK 27.0; select it with DEVELOPER_DIR and SDKROOT" >&2; exit 1; }
	@test -z "$$(gofmt -l cmd internal tools tests)"
	xcrun swift-format lint --strict --recursive macos/Package.swift macos/Sources macos/Tests
	go vet ./...
	go test -race -v ./...
	go run ./tools/swifttest
	$(MAKE) app

app:
	go run ./tools/bundle

install-dev:
	go run ./tools/bundle install

verify-launchd:
	go run ./tools/lifecycle -mode launchd

verify-lifecycle:
	go run ./tools/lifecycle -mode attended

format:
	gofmt -w cmd internal tools tests
	xcrun swift-format format --in-place --recursive macos/Package.swift macos/Sources macos/Tests
