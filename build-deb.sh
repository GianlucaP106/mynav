#!/bin/bash
set -e

# Build script for creating .deb packages
# Usage: ./build-deb.sh [architecture]

ARCH=${1:-"amd64"}
PACKAGE_NAME="mynav"

echo "Building .deb package for architecture: $ARCH"

# Check if required tools are installed
check_dependencies() {
    echo "Checking dependencies..."
    
    if ! command -v dpkg-buildpackage &> /dev/null; then
        echo "Error: dpkg-buildpackage not found"
        echo "Install with: sudo apt-get install dpkg-dev"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        echo "Error: Go not found"
        echo "Install with: sudo apt-get install golang-go"
        exit 1
    fi
    
    echo "Dependencies OK"
}

# Clean previous builds
clean_build() {
    echo "Cleaning previous builds..."
    rm -f ../${PACKAGE_NAME}_*.deb
    rm -f ../${PACKAGE_NAME}_*.tar.xz
    rm -f ../${PACKAGE_NAME}_*.dsc
    rm -f ../${PACKAGE_NAME}_*.changes
    rm -f ../${PACKAGE_NAME}_*.buildinfo
    rm -f ${PACKAGE_NAME}
}

# Build the package
build_package() {
    echo "Building package..."
    
    if [ "$ARCH" = "all" ]; then
        # Build for multiple architectures
        echo "Building for multiple architectures..."
        dpkg-buildpackage -a amd64 -us -uc -b
        dpkg-buildpackage -a arm64 -us -uc -b
        dpkg-buildpackage -a armhf -us -uc -b
    else
        # Build for specific architecture
        dpkg-buildpackage -a "$ARCH" -us -uc -b
    fi
}

# Verify the package
verify_package() {
    echo "Verifying package..."
    
    DEB_FILE=$(ls ../${PACKAGE_NAME}_*.deb | head -1)
    
    if [ -f "$DEB_FILE" ]; then
        echo "Package created successfully: $DEB_FILE"
        echo "Package info:"
        dpkg -I "$DEB_FILE"
        echo ""
        echo "Package contents:"
        dpkg -c "$DEB_FILE"
    else
        echo "Error: Package not found"
        exit 1
    fi
}

# Main execution
main() {
    echo "Starting .deb package build process..."
    
    check_dependencies
    clean_build
    build_package
    verify_package
    
    echo ""
    echo "Build completed successfully!"
    echo "To install: sudo dpkg -i ../${PACKAGE_NAME}_*.deb"
    echo "To remove: sudo dpkg -r ${PACKAGE_NAME}"
}

# Run main function
main "$@"