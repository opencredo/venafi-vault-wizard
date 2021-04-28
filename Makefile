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
	# --unroll-variadic=false means that variadic arguments get passed into testify/mock as a slice instead of as
	# individual arguments. This means that mock.Anything can be used to ignore the arguments, instead of having to type
	# mock.Anything as many times as there arguments and therefore having to specify how many arguments there will be.
	$(BIN)/mockery --all --keeptree --unroll-variadic=false

$(BIN)/mockery:
	@echo "Installing mockery to generate mocks...\n"
	GOBIN=$(BIN) go install github.com/vektra/mockery/v2@v2.6.0

.PHONY: clean
clean:
	rm -rf $(BIN)
