APP      := deskwall
PKG      := deskwall
VERSION  := 0.1.0
DIST     := dist

GO       ?= go
NFPM     ?= nfpm
FYNE     ?= asdf exec fyne --

.PHONY: all build run tidy clean fmt vet test deb windows install-tools help

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

# Install nfpm (Debian packaging) and fyne CLI.
install-tools:
	$(GO) install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
	$(GO) install fyne.io/tools/cmd/fyne@latest

# Cross-compile for Windows (requires gcc-mingw-w64-x86-64).
# Install with: sudo apt install gcc-mingw-w64-x86-64
windows:
	@mkdir -p $(DIST)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
		$(GO) build -trimpath -ldflags "-s -w -H windowsgui" -o $(DIST)/$(APP).exe .

# Cross-compile for Windows and package with fyne (bundles GLFW DLLs).
# Produces dist/deskwall.zip with the exe + required DLLs.
windows-pkg: windows
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
		$(FYNE) package -os windows -icon assets/deskwall.png -executable $(DIST)/$(APP).exe

# Remove build artifacts.
clean:
	rm -rf $(DIST)

help:
	@echo "Targets: build run tidy fmt vet test deb windows windows-pkg install-tools clean"
