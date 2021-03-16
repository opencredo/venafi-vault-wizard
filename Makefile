BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o $(BIN)/vvw

.PHONY: generate-mocks
generate-mocks: $(BIN)/mockery
	@echo "Generating mocks...\n"
	$(BIN)/mockery --all

$(BIN)/mockery:
	@echo "Installing mockery to generate mocks...\n"
	GOBIN=$(BIN) go install github.com/vektra/mockery/v2@v2.6.0