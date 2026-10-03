# Wallplay

A desktop application built with Go and [Fyne](https://fyne.io).

## Prerequisites

Fyne uses cgo (GLFW + OpenGL) but does NOT need GTK development headers.
Build dependencies on Debian/Ubuntu:

```sh
sudo apt install libgl1-mesa-dev libegl1-mesa-dev \
                 libx11-dev libxrandr-dev libxinerama-dev \
                 libxcursor-dev libxi-dev libxxf86vm-dev
```

Runtime needs only the userland libraries (listed in `nfpm.yaml`).

## Project structure

```
wallplay/
├── main.go              # Entry point
├── internal/
│   ├── app/             # Fyne application wiring and lifecycle
│   └── ui/              # Fyne windows and widgets
├── assets/             # Desktop entry, icons
├── nfpm.yaml           # Debian (.deb) packaging config
└── Makefile
```

## Build & run

```sh
make run      # build and run
make build    # build to dist/
make tidy     # resolve dependencies
make test
```

## Debian package

Packaging uses [nfpm](https://nfpm.goreleaser.com/), a single Go binary that
produces `.deb` (and `.rpm`/`.apk`) packages without the full Debian toolchain.

```sh
make install-tools   # installs nfpm
make deb             # builds binary + .deb into dist/
```

Install the resulting package:

```sh
sudo apt install ./dist/wallplay_0.1.0_amd64.deb
```
