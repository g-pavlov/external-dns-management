EXECUTABLE=dns-controller-manager
PROJECT=github.com/gardener/external-dns-management
VERSION=$(shell cat VERSION)

.PHONY: revendor
revendor:
	@GO111MODULE=on go mod vendor
	@GO111MODULE=on go mod tidy


.PHONY: check
check:
	@.ci/check

.PHONY: build
build:
    @CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o $(EXECUTABLE) \
        -mod=vendor \
	    -ldflags "-X main.Version=$(VERSION)-$(shell git rev-parse HEAD)"\
	    ./cmd/dns

.PHONY: build-local
build-local:
	@CGO_ENABLED=1 GO111MODULE=on go build -o $(EXECUTABLE) \
	    -mod=vendor \
	    -race \
	    -gcflags="all=-N -l" \
	    -ldflags "-X main.Version=$(VERSION)-$(shell git rev-parse HEAD)"\
	    ./cmd/dns

.PHONY: build-local-compound
build-local-compound:
	@CGO_ENABLED=1 GO111MODULE=on go build -o $(EXECUTABLE)-compound \
	    -mod=vendor \
	    -race \
	    -gcflags="all=-N -l" \
	    -ldflags "-X main.Version=$(VERSION)-$(shell git rev-parse HEAD)"\
	    ./cmd/compound

.PHONY: release
release:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o $(EXECUTABLE) \
	    -a \
	    -mod=vendor \
	    -ldflags "-w -X main.Version=$(VERSION)" \
	    ./cmd/dns

.PHONY: test
test:
	GO111MODULE=on go test -mod=vendor ./pkg/...
	@echo ----- Skipping long running integration tests, use \'make alltests\' to run all tests -----
	test/integration/run.sh $(kindargs) -- -skip Many $(args)

.PHONY: generate
generate:
	@./hack/generate-code
	@GO111MODULE=on go generate ./pkg/apis/dns/...
	@GO111MODULE=on go generate ./charts/external-dns-management

.PHONY: install-requirements
install-requirements:
	@go install -mod=vendor github.com/onsi/ginkgo/ginkgo
	@GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	@./hack/install-requirements.sh

alltests:
	GO111MODULE=on go test -mod=vendor ./pkg/...
	test/integration/run.sh $(kindargs) -- $(args)
