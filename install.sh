#! /bin/bash

mkdir ~/.mynav/ 2>/dev/null | true
cd ~/.mynav

git clone https://github.com/GianlucaP106/mynav src
cd src
go build
mv mynav ..
