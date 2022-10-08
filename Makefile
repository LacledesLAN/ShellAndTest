GO = CGO_ENABLED=0 go
GO_BUILD = $(GO) build

.PHONY : test
test :
	$(GO) test -mod=vendor ./...
	golangci-lint run --new-from-rev master --sort-results --out-format tab

.PHONY: lint
lint :
	golangci-lint run --sort-results --out-format tab

.PHONY: build
build:
	mkdir -p build
	$(GO_BUILD) -o build/shellandtest cmd/main.go
