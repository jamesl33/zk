PACKAGE=

TESTS=

default: build

.PHONY: generate
generate:
	@go generate ./$(PACKAGE)/...

.PHONY: build
build: generate
	@go build -o build/zk

.PHONY: install
install:
	@go install .

.PHONY: test
test:
	@go test ./$(PACKAGE)/... -run=$(TESTS) -count=1 -covermode=atomic -coverprofile=coverage.out -failfast -shuffle=on 

.PHONY: coverage
coverage: test
	@go tool cover -html=coverage.out

.PHONY: clean
clean:
	@rm -f coverage.out
	@fd --no-ignore --glob '*_enum.go' | xargs rm -f
