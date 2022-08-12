
.PHONY: lint
lint: lint-go

.PHONY: lint-go
lint-go:
	golangci-lint run -v

.PHONY: test
test: 
	go test ./... -coverprofile=cover.out

.PHONY: cover
cover: 
	go tool cover -html=cover.out

.PHONY: build
build:
	go build .

.PHONY: serve
serve:
	go run .

.PHONY: generate
generate:
	go generate ./...

.PHONY: db-admin-start
db-admin-start:
	docker compose up -d > docker-compose-up.out
	$(info DB is running)

.PHONY: db-admin-stop
db-admin-stop:
	docker compose down > docker-compose-down.out

