#!/usr/bin/env bash

mkdir -p dist/uve-0.1-linux-x86_64
cp uve-linux-amd64 dist/uve-0.1-linux-x86_64/uve-bin
cp uve.sh dist/uve-0.1-linux-x86_64/uve.sh
tar -zcvf dist/uve-0.1-linux-x86_64.tgz dist/uve-0.1-linux-x86_64

mkdir -p dist/uve-0.1-macos-arm64
cp uve-darwin-arm64 dist/uve-0.1-macos-arm64/uve-bin
cp uve.sh dist/uve-0.1-macos-arm64/uve.sh
tar -zcvf dist/uve-0.1-macos-arm64.tgz dist/uve-0.1-macos-arm64

mkdir -p dist/uve-0.1-macos-x86_64
cp uve-darwin-amd64 dist/uve-0.1-macos-x86_64/uve-bin
cp uve.sh dist/uve-0.1-macos-x86_64/uve.sh
tar -zcvf dist/uve-0.1-macos-x86_64.tgz dist/uve-0.1-macos-x86_64

mkdir -p dist/uve-0.1-windows-x86_64
cp uve-windows-amd64.exe dist/uve-0.1-windows-x86_64/uve-bin.exe
cp uve.ps1 dist/uve-0.1-windows-x86_64/uve.ps1
zip -r -X dist/uve-0.1-windows-x86_64.zip dist/uve-0.1-windows-x86_64
