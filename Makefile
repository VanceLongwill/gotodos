BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
PKGS := $(shell go list ./... | grep -v /vendor)

build:
	go build

.PHONY: test
test: 
	go test $(PKGS) -test.v

coverage:
	go test $(PKGS) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage_report.html

$(GOMETALINTER): 
	go get -u github.com/alecthomas/gometalinter 
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor
