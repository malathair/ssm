BUILD_DIR=build
EXECUTABLE=ssm

DARWIN_AMD64=darwin_amd64
LINUX_AMD64=linux_amd64
WINDOWS_AMD64=windows_amd64

VERSION=$(shell git describe --tags --always)

.PHONY: all clean $(PLATFORMS)

all: build

build: darwin linux windows
	@echo "Generating SHA256 checksums..."
	@sha256sum $(BUILD_DIR)/* > $(BUILD_DIR)/$(EXECUTABLE)_$(VERSION)_checksums.txt

darwin: $(DARWIN_AMD64)

linux: $(LINUX_AMD64)

windows: $(WINDOWS_AMD64)

$(DARWIN_AMD64):
	@echo "Building binary for $(subst __amd64,,$@)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(subst _amd64,,$@) GOARCH=amd64 go build -o $(BUILD_DIR)/$(EXECUTABLE)_$@

$(LINUX_AMD64):
	@echo "Building binary for $(subst __amd64,,$@)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(subst _amd64,,$@) GOARCH=amd64 go build -o $(BUILD_DIR)/$(EXECUTABLE)_$@

$(WINDOWS_AMD64):
	@echo "Building binary for $(subst __amd64,,$@)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(subst _amd64,,$@) GOARCH=amd64 go build -o $(BUILD_DIR)/$(EXECUTABLE)_$@.exe

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
