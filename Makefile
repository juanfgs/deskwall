APP      := wallplay
PKG      := wallplay
VERSION  := 0.1.0
DIST     := dist

GO       ?= go
NFPM     ?= nfpm

.PHONY: all build run tidy clean fmt vet test deb install-tools help

all: build

# Build the binary into dist/.
build:
	@mkdir -p $(DIST)
	$(GO) build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(APP) .

# Build and run the application.
run: build
	$(DIST)/$(APP)

# Resolve and tidy module dependencies.
tidy:
	$(GO) mod tidy

# Format and lint.
fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

# Build a .deb package using nfpm.
deb: build
	@mkdir -p $(DIST)
	$(NFPM) pkg --packager deb --target $(DIST)/ --config nfpm.yaml

# Install nfpm (Debian packaging tool) and gotk4 deps helper.
install-tools:
	$(GO) install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest

# Remove build artifacts.
clean:
	rm -rf $(DIST)

help:
	@echo "Targets: build run tidy fmt vet test deb install-tools clean"
