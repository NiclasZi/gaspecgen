# Variables
CLI_BINARY_NAME=gaspecgen
BUILD_DIR=bin
EXT=$(if $(filter windows,$(GOOS)),.exe,)
SERVICE_TEMPLATE_DIR=templates
VERSION=$(shell git describe --tags --abbrev=0)

# Targets
.PHONY: all build/* test release install clean lint docs

all: build/$(CLI_BINARY_NAME)

build/%:
	@echo "Building $*..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X github.com/Phillezi/gaspecgen/cmd/$*/cli.version=$(VERSION)" -o $(BUILD_DIR)/$*$(EXT) ./cmd/$*
	@echo "Build complete: $(BUILD_DIR)/$*$(EXT)"

test:
	@go test ./...

release/%:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -mod=readonly -ldflags "-w -s -X github.com/Phillezi/gaspecgen/cmd/$*/cli.version=$(VERSION)" -o $(BUILD_DIR)/$*$(EXT) ./cmd/$*
	@echo "Build complete."

install: release
	@echo "installing"
	@./scripts/escalate.sh cp ./$(BUILD_DIR)/$(CLI_BINARY_NAME)$(EXT) /usr/local/bin/$(CLI_BINARY_NAME)$(EXT)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

lint:
	@./scripts/check-lint.sh

docs:
	@go run ./cmd/docs
