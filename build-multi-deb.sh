#!/bin/bash
set -e

echo "Building .deb packages for multiple architectures..."

# Build for amd64
echo "Building for amd64..."
export GOARCH=amd64
export CC=gcc
export CGO_ENABLED=1
export DEB_BUILD_ARCH=amd64
dpkg-buildpackage -a amd64 -us -uc -b

# Build for arm64
echo "Building for arm64..."
export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
export CGO_ENABLED=1
export DEB_BUILD_ARCH=arm64
dpkg-buildpackage -a arm64 -us -uc -b

# Build for armhf
echo "Building for armhf..."
export GOARCH=arm
export CC=arm-linux-gnueabihf-gcc
export CGO_ENABLED=1
export DEB_BUILD_ARCH=armhf
dpkg-buildpackage -a armhf -us -uc -b

echo "All packages built successfully!"
ls -la ../mynav_*.deb 