# Copyright (c) 2026 Parable, LLC
export MACOSX_DEPLOYMENT_TARGET := 26.0
export CGO_ENABLED := 1
export CGO_CFLAGS := -mmacosx-version-min=26.0
export CGO_LDFLAGS := -mmacosx-version-min=26.0
export CODE_SIGN_IDENTITY

.PHONY: check app install-dev verify-launchd verify-lifecycle format
check:
	@test "$$(uname -m)" = arm64
	@test "$$(go env GOVERSION)" = go1.26.5
	@swift --version | rg -q 'Swift version 6.4'
	@test "$$(xcrun --show-sdk-version)" = 27.0
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
