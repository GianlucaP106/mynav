# MyNav Makefile

.PHONY: build deb clean install

# Build binary
build:
	go build -ldflags="-s -w" -o mynav .

# Build Debian package
deb: clean
	dpkg-buildpackage -us -uc -b

# Build for multiple architectures
deb-multi:
	dpkg-buildpackage -a amd64 -us -uc -b
	dpkg-buildpackage -a arm64 -us -uc -b
	dpkg-buildpackage -a armhf -us -uc -b

# Install locally
install: build
	install -D -m 755 mynav $(DESTDIR)/usr/bin/mynav

# Clean build artifacts
clean:
	rm -f mynav
	rm -f ../mynav_*.deb
	rm -f ../mynav_*.tar.xz
	rm -f ../mynav_*.dsc
	rm -f ../mynav_*.changes
	rm -f ../mynav_*.buildinfo

# Development build with debug info
dev:
	go build -o mynav .

# Test package installation
test-install:
	sudo dpkg -i ../mynav_*.deb

# Remove installed package
uninstall:
	sudo dpkg -r mynav