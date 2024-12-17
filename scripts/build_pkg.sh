#!/bin/bash

# Get the build directory from first argument
BUILD_DIR="$1"
VERSION="$2"

# Get architecture from uname -m
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
elif [ "$ARCH" = "x86_64" ]; then
    ARCH="x86_64"
fi

# Check if TGZ exists
if [ ! -f "${BUILD_DIR}/libpoktroll_clients-${VERSION}-${ARCH}-linux.tar.gz" ]; then
    cmake --build "${BUILD_DIR}" --target package
fi

# Ensure pkg directory exists
mkdir -p "${BUILD_DIR}/pkg"

# Copy files to pkg directory
cp "${BUILD_DIR}/PKGBUILD" "${BUILD_DIR}/pkg/"
cp "${BUILD_DIR}/libpoktroll_clients-${VERSION}-${ARCH}-linux.tar.gz" "${BUILD_DIR}/pkg/"

# Build the package
cd "${BUILD_DIR}/pkg" && makepkg -f