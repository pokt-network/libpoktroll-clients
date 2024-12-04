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

# Build deb/rpm/tar install packages.
cpack         # All
cpack -G DEB  # Debian
cpack -G RPM  # RHEL/CentOS
cpack -G TGZ  # tar.gz
# Produces:
#  - build/libpoktroll_clients-<version>-Linux.tar.gz
#  - build/libpoktroll_clients-<version>_.deb
#  - build/libpoktroll_clients-<version>_amd64.deb
#  - build/libpoktroll_clients-<version>.x86_64.rpm

# Build arch install package (depends on TGZ from cpack).
make pkgbuild
# Produces:
#  - build/PKGBUILD
#  - build/pkg/*

# Install the shared library and headers.

## Arch
sudo pacman -U ./pkg/libpoktroll_clients-0.1.0-1-x86_64.pkg.tar.zst   

## Debian
sudo dpkg -i ./pkg/libpoktroll_clients-0.1.0-Linux.deb

## RHEL/CentOS
sudo rpm -i ./pkg/libpoktroll_clients-0.1.0-Linux.rpm
```