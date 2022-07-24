
.PHONY: lint
lint: lint-go

.PHONY: lint-go
lint-go:
	golangci-lint run -v

.PHONY: test
test: 
	go test ./...

.PHONY: build
build:
	go build .

.PHONY: serve
serve:
	go run .

.PHONY: generate
generate:
	go generate ./...

