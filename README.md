# Poktroll C Clients Library

This repo contains an asynchronous C API to the [Poktroll client packages](https://pkg.go.dev/github.com/pokt-network/poktroll/pkg/client) via [cgo](https://pkg.go.dev/cmd/cgo) wrapper functions.

## Getting started
```bash
# Clone and cd into the repo.
git clone https://github.com/bryanchriswhite/libpoktroll_clients.git --recurse-submodules
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