# `libpoktroll` - C Clients Shared Library

An asynchronous C API which wraps the [`poktroll` client packages](https://pkg.go.dev/github.com/pokt-network/poktroll/pkg/client) (via [cgo](https://pkg.go.dev/cmd/cgo) wrapper functions).

## Table of contents <!-- omit in toc -->

- [Installation](#installation)
- [Getting started (development environment)](#getting-started-development-environment)
- [Linux](#linux)
    - [Building shared libraries](#building-shared-libraries-1)
    - [Building installers](#building-installers-1)
- [MacOS](#macos)
- [Cross-compiling from (arch) linux](#cross-compiling-from-arch-linux-1)
    - [Targeting Windows](#targeting-windows)
    - [Targeting macOS](#targeting-macos)


## Installation

You can _EITHER_ **download** or **build** an OS/architecture-specific installer or shared library.

### Download

Downloads are available via the [releases page](https://github.com/bryanchriswhite/libpoktroll-clients/releases).

Installers are preferred as they only need to be run, whereas shared libraries **need to be renamed and saved to one of the OS's library search paths** (e.g. `/usr/lib/`, `/usr/local/lib/`, `~/.local/lib/`, etc.).
Depending on your OS, the file extension will either be `.so`, `.dylib`, or `.dll` for linux, macOS, and Windows, respectively.

E.g.: For v0.1.0 x86_64 linux, `libpoktroll_clients-v0.1.0-amd64.so` -> `/usr/local/lib/libpoktroll_clients.so`

### Build

To build from source, complete the [getting started](#getting-started-development-environment) section below, then run `sudo make install` or see the platform-specific sections below.

As with downloaded shared libraries (see above), the steps outlined in the platform-specific sections below will produce an OS/architecture-specific shared library, which will **need to be renamed and saved to one of the OS's library search paths**.

E.g.: For v0.1.0 x86_64 linux, `libpoktroll_clients-v0.1.0-amd64.so` -> `/usr/local/lib/libpoktroll_clients.so`

## Getting started (development environment)
```bash
# Clone and cd into the repo.
git clone https://github.com/bryanchriswhite/libpoktroll-clients.git --recurse-submodules
cd libpoktroll_clients

# If you cloned but didn't pull the submodules, run:
git submodule update --init --recursive

# (optional) Update protobufs ("pull" from buf.build)
buf export buf.build/pokt-network/poktroll

# Generate protobufs
buf generate

# Make and cd to cmake build directory.
mkdir build
cd build

# Generate build files and build the library..
cmake ..
make

# Run tests (requires running poktroll localnet.
# (see: https://dev.poktroll.com/develop/developer_guide/quickstart#1-launch--inspect-localnet)
ctest --output-on-failure
```

## Linux

### Building shared libraries

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./build/libpoktroll_clients-<version>-amd64.so -buildmode=c-shared .
```

### Building installers

```bash
cd libpoktroll_clients/build
cmake ..

# Build deb/rpm/tar install packages.
make package  # ALL

## Produces:
#  - build/libpoktroll_clients-<version>-<arch>-linux.{sh,tar.gz}
#  - build/libpoktroll_clients-<version>-<arch>-linux.deb
#  - build/libpoktroll_clients-<version>-<arch>-linux.rpm

## OR
cpack -G "TGZ;DEB;RPM"  # All
cpack -G DEB            # Debian/Ubuntu
cpack -G RPM            # RHEL/Fedora
cpack -G TGZ            # tar.gz

# Build arch install package (depends on TGZ from cpack).
make pkgbuild

## Produces:
#  - build/PKGBUILD
#  - build/pkg/*

# Install the shared library and headers.

## Arch
sudo pacman -U ./pkg/libpoktroll_clients-<version>-<arch>-Linux.pkg.tar.zst   

## Debian
sudo dpkg -i ./libpoktroll_clients-<version>-<arch>-Linux.deb

## RHEL/CentOS
sudo rpm -i ./libpoktroll_clients-<version>-<arch>-Linux.rpm
```

## MacOS

### Building shared libraries (on a macOS host)

```bash
# MacOS (Intel)
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ./build/libpoktroll_clients-<version>-amd64.dylib -buildmode=c-shared .

# MacOS (ARM)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ./build/libpoktroll_clients-<version>-arm64.dylib -buildmode=c-shared .
```

### Building installers

```bash
cd poktroll-clients-py
mkdir ./build && cd ./build
cmake ..

cpack -G "productbuild;TGZ"
```

## Cross-compiling from (arch) linux

### Targeting Windows

```bash
## Dependencies
sudo pacman -S mingw-w64-gcc wine wine-mono wine_gecko winetricks
yay -S llvm-mingw-w64-toolchain-ucrt-bin

## Shared library
CGO_ENABLED=1 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
GOOS=windows \
GOARCH=amd64 \
go build -o ./build/libpoktroll_clients-<version>-amd64.dll -buildmode=c-shared .

## Installer
TODO: While this is possible, it's a bit involved and out of scope for the initial release.
```

### Targeting macOS

**TODO: While this is possible, it's a bit involved and out of scope for the initial release.**