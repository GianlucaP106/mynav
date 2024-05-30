#! /bin/bash

mkdir ~/.mynav/ 2>/dev/null | true
cd ~/.mynav

rm -rf src

latest_release=$(curl -sL "https://api.github.com/repos/GianlucaP106/mynav/releases/latest")
latest_tag=$(echo "$latest_release" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

curl -fsSL https://github.com/GianlucaP106/mynav/archive/refs/tags/${latest_tag}.zip -o src.zip

unzip -d src "src.zip"
rm -rf "src.zip"

tag_without_v=$(echo "$latest_tag" | tr -d "v")
cd "src/mynav-${tag_without_v}"
go build -o mynav
mv mynav ../../mynav
