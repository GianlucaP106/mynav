#! /bin/bash

echo "Installing mynav..."

echo "Initializing configuration directory at ~/.mynav"
mkdir ~/.mynav/ 2>/dev/null | true
cd ~/.mynav

latest_release=$(curl -sL "https://api.github.com/repos/GianlucaP106/mynav/releases/latest")
latest_tag=$(echo "$latest_release" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

echo "Installing latest release ${latest_tag}"

os="$(uname -s)"
arch="$(uname -m)"

if [ "$arch" = "aarch64" ]; then
    arch="arm64"
fi

echo "Detected platform ${os}-${arch}"

file="mynav_${os}_${arch}.tar.gz"

curl -s -L -o build.tar.gz https://github.com/GianlucaP106/mynav/releases/download/${latest_tag}/${file}

tar -xzf build.tar.gz 2>/dev/null
if [ $? -ne 0 ]; then
    echo "Failed to install mynav. Platform not supported."
    exit
fi

rm build.tar.gz

directory="$HOME/.mynav"
if echo "$PATH" | grep -q "$directory"; then
    echo "Successfully installed mynav at ~/.mynav!"
else
    echo "Adding ~/.mynav to PATH"
    rc_files=(~/.bashrc ~/.zshrc ~/.profile ~/.bash_profile ~/.bash_login ~/.cshrc ~/.tcshrc)
    for file in "${rc_files[@]}"; do
        if [ -f "$file" ]; then
            echo "Adding to $file"
            echo '# mynav path export' >>"$file"
            echo 'export PATH="$PATH:$HOME/.mynav"' >>"$file"
            break
        fi
    done
    echo "Successfully installed mynav at ~/.mynav!"
    echo "Restart your terminal session."
fi
