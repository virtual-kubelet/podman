LINTER_BIN ?= golangci-lint

GO111MODULE := on
export GO111MODULE

.PHONY: build
build: clean bin/virtual-kubelet bin/virtual-kubelet-arm

.PHONY: clean
clean: files := bin/virtual-kubelet bin/virtual-kubelet-arm
clean:
	@rm $(files) &>/dev/null || exit 0

.PHONY: test
test:
	@echo running tests
	go test -v ./...

.PHONY: vet
vet:
	@go vet ./... #$(packages)

.PHONY: lint
lint:
	@$(LINTER_BIN) run --new-from-rev "HEAD~$(git rev-list master.. --count)" ./...

.PHONY: check-mod
check-mod: # verifies that module changes for go.mod and go.sum are checked in
	@hack/ci/check_mods.sh

.PHONY: mod
mod:
	@go mod tidy

bin/virtual-kubelet: BUILD_VERSION          ?= $(shell git describe --tags --always --dirty="-dev")
bin/virtual-kubelet: BUILD_DATE             ?= $(shell date -u '+%Y-%m-%d-%H:%M UTC')
bin/virtual-kubelet: VERSION_FLAGS    := -ldflags='-X "main.buildVersion=$(BUILD_VERSION)" -X "main.buildTime=$(BUILD_DATE)"'

bin/virtual-kubelet-arm: BUILD_VERSION          ?= $(shell git describe --tags --always --dirty="-dev")
bin/virtual-kubelet-arm: BUILD_DATE             ?= $(shell date -u '+%Y-%m-%d-%H:%M UTC')
bin/virtual-kubelet-arm: VERSION_FLAGS    := -ldflags='-X "main.buildVersion=$(BUILD_VERSION)" -X "main.buildTime=$(BUILD_DATE)"'
bin/virtual-kubelet-arm: GOARCH    := arm

bin/%:
	CGO_ENABLED=0 GOARCH=$(GOARCH) go build -ldflags '-extldflags "-static"' -o bin/$(*) $(VERSION_FLAGS) ./cmd/virtual-kubelet

run: clean build
	./bin/virtual-kubelet --provider podman --nodename podman --provider-config ./deploy/systemd/podman-cfg.json --full-resync-period=10s
