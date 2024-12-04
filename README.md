## Poktroll C Clients Library

This repo contains an asynchronous C API to the [Poktroll client packages](https://pkg.go.dev/github.com/pokt-network/poktroll/pkg/client) via [cgo](https://pkg.go.dev/cmd/cgo) wrapper functions.

```bash
# Clone and cd into the repo.
git clone https://github.com/bryanchriswhite/libpoktroll_clients.git
cd libpoktroll_clients

# Make and cd to cmake build directory.
mkdir build
cd build

# Generate build files and build the library..
cmake ..
make

# Run tests (requires [running poktroll localnet](https://dev.poktroll.com/develop/developer_guide/quickstart#1-launch--inspect-localnet)).
ctest --output-on-failure

# Build tar/deb/rpm install packages.
cpack  # Produces:
#  - build/libpoktroll_clients-<version>-Linux.tar.gz
#  - build/libpoktroll_clients-<version>_.deb
#  - build/libpoktroll_clients-<version>_amd64.deb
#  - build/libpoktroll_clients-<version>.x86_64.rpm

# Build arch install package.
make pkgbuild  # Produces: build/pkg/*
```