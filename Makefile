PROJECT_NAME := report-tfvars
OUTPUT_DIR := bin
OUTPUTS_DIR := outputs
COVER_OUT := $(OUTPUTS_DIR)/coverage.out
COVER_JSON := $(OUTPUTS_DIR)/coverage.json
GO_BUILD := go build -o
TARGETS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64


# Default target
all: prepare $(TARGETS) install

prepare:
	@go mod tidy
	@go mod download
	@gofmt -w report-tfvars.go
	@gofmt -w report-tfvars_test.go

install:
	@echo "Installed inside $(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)"
	@rm -rf $(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)
	@cp $(CI_PROJECT_DIR)/go/$(PROJECT_NAME)/bin/$(PROJECT_NAME)-linux-amd64 $(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)
	@chmod +x $(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)
	@bash -c '"$(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)" --help'

uninstall:
	@rm -rf $(CI_PROJECT_DIR)/bin/$(PROJECT_NAME)

remove: uninstall

delete: uninstall

# Build targets for each OS/Arch combination
$(TARGETS): 
	@echo "Building for GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@))..."
	GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) $(GO_BUILD) $(OUTPUT_DIR)/$(PROJECT_NAME)-$(word 1,$(subst /, ,$@))-$(word 2,$(subst /, ,$@)) .

# Clean up binaries
clean: uninstall
	@rm -rf $(OUTPUT_DIR)
	@rm -rf $(OUTPUTS_DIR)
	
# Run the package
run: prepare
	go run . $(ARGS)

test: prepare
	@mkdir -p $(OUTPUTS_DIR)
	@go test -json  ./... $(ARGS) > $(OUTPUTS_DIR)/tests.json 2> /dev/null
	@go test ./... $(ARGS)

# Help target
help:
	@echo "Makefile for $(PROJECT_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make [target]"
	@echo
	@echo "Targets:"
	@echo "  all       Build binaries for all target OS/Arch combinations"
	@echo "  run       Run the Go code directly"
	@echo "  clean     Remove all binaries"
	@echo "  help      Display this help message"
	@echo
	@echo "Target OS/Arch combinations:"
	@echo "  darwin/amd64"
	@echo "  darwin/arm64"
	@echo "  linux/amd64"
	@echo "  linux/arm64"
	@echo "  windows/amd64"

coverage:
	@mkdir -p $(OUTPUTS_DIR)
	@go test -coverprofile=$(COVER_OUT)
	@go tool cover -func=$(COVER_OUT)
	@go tool cover -o $(COVER_JSON) -func=$(COVER_OUT)

.PHONY: all run test coverage install uninstall remove delete clean help $(TARGETS)
