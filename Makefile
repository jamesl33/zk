# TODO
PACKAGE=

# TODO
TESTS=

# TODO
default: build

# TODO
.PHONY: generate
generate:
	@go generate ./$(PACKAGE)/...

# TODO
.PHONY: build
build: generate
	@go build -o build/zk

# TODO
.PHONY: test
test:
	@go test ./$(PACKAGE)/... -run=$(TESTS) -count=1 -covermode=atomic -coverprofile=coverage.out -failfast -shuffle=on 

# TODO
.PHONY: coverage
coverage: test
	@go tool cover -html=coverage.out

# TODO
.PHONY: clean
clean:
	@rm -f coverage.out
	@fd --no-ignore --glob '*_enum.go' | xargs rm -f
