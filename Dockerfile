# Multi-architecture .deb package builder for MyNav
FROM debian:bookworm-slim

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV GOPATH=/go
ENV GO_VERSION=1.22
ENV PATH=$GOPATH/bin:$PATH

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    dpkg-dev \
    debhelper \
    devscripts \
    fakeroot \
    lintian \
    git \
    wget \
    curl \
    ca-certificates \
    golang-go \
    gcc-aarch64-linux-gnu \
    gcc-arm-linux-gnueabihf \
    g++-aarch64-linux-gnu \
    g++-arm-linux-gnueabihf \
    libc6-dev-arm64-cross \
    libc6-dev-armhf-cross \
    && rm -rf /var/lib/apt/lists/*

# Install Go cross-compilation tools
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set up cross-compilation environment
ENV GOOS=linux
ENV CGO_ENABLED=1

# Create build directory
WORKDIR /build

# Copy source code and build script
COPY . .
COPY build-multi-deb.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/build-multi-deb.sh

# Default command
# CMD ["/usr/local/bin/build-multi-deb.sh"] 
CMD ["/bin/bash"]
