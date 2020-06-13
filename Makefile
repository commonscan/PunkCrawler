ifeq ($(OS),Windows_NT)
  EXECUTABLE_EXTENSION := .exe
else
  EXECUTABLE_EXTENSION :=
endif

GO_FILES = $(shell find . -type f -name '*.go')
TEST_MODULES ?=

all: punkspider

.PHONY: all clean integration-test integration-test-clean docker-runner container-clean gofmt test

gofmt:
	goimports -w -l $(GO_FILES)

punkspider: $(GO_FILES)
	cd cmd/ && go build -o punk_crawler$(EXECUTABLE_EXTENSION) && cd ../..
	rm -f punk_crawler2
	ln -s cmd/punk_crawler$(EXECUTABLE_EXTENSION) punk_crawler