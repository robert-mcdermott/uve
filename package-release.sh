#!/usr/bin/env bash

# Configuration
VERSION="0.1.5"
PACKAGE="main.go"
BINARY_NAME="uve-bin"
PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64")

echo "=== Building UVE v${VERSION} ==="

# Compile for all platforms
for platform in "${PLATFORMS[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    # Set output filename based on platform
    if [ "$GOOS" = "windows" ]; then
        output_name="${BINARY_NAME}-${GOOS}-${GOARCH}.exe"
    else
        output_name="${BINARY_NAME}-${GOOS}-${GOARCH}"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name $PACKAGE
    
    if [ $? -ne 0 ]; then
        echo "Error building for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "=== Creating release packages ==="

# Create a temporary README.txt with installation instructions
cat > README.txt << 'EOL'
UVE - UV Environment Manager

INSTALLATION:

Linux/macOS:
1. Copy uve-bin to a directory in your PATH (e.g., ~/bin/)
2. Run: uve-bin init
3. Restart your shell or source your shell config file

Windows:
1. Copy uve-bin.exe to a directory in your PATH (e.g., %USERPROFILE%\bin\)
2. Run: uve-bin.exe init
3. Start a new PowerShell session or run: Import-Module uve

For full documentation, visit: https://github.com/robert-mcdermott/uve
EOL

# Create distribution directories
mkdir -p dist

# Linux x86_64
mkdir -p tmp/uve-${VERSION}-linux-x86_64
cp ${BINARY_NAME}-linux-amd64 tmp/uve-${VERSION}-linux-x86_64/${BINARY_NAME}
cp README.txt tmp/uve-${VERSION}-linux-x86_64/
cp LICENSE tmp/uve-${VERSION}-linux-x86_64/
tar -zcvf dist/uve-${VERSION}-linux-x86_64.tar.gz -C tmp uve-${VERSION}-linux-x86_64

# macOS ARM64
mkdir -p tmp/uve-${VERSION}-macos-arm64
cp ${BINARY_NAME}-darwin-arm64 tmp/uve-${VERSION}-macos-arm64/${BINARY_NAME}
cp README.txt tmp/uve-${VERSION}-macos-arm64/
cp LICENSE tmp/uve-${VERSION}-macos-arm64/
tar -zcvf dist/uve-${VERSION}-macos-arm64.tar.gz -C tmp uve-${VERSION}-macos-arm64

# macOS x86_64
mkdir -p tmp/uve-${VERSION}-macos-x86_64
cp ${BINARY_NAME}-darwin-amd64 tmp/uve-${VERSION}-macos-x86_64/${BINARY_NAME}
cp README.txt tmp/uve-${VERSION}-macos-x86_64/
cp LICENSE tmp/uve-${VERSION}-macos-x86_64/
tar -zcvf dist/uve-${VERSION}-macos-x86_64.tar.gz -C tmp uve-${VERSION}-macos-x86_64

# Windows x86_64
mkdir -p tmp/uve-${VERSION}-windows-x86_64
cp ${BINARY_NAME}-windows-amd64.exe tmp/uve-${VERSION}-windows-x86_64/${BINARY_NAME}.exe
cp README.txt tmp/uve-${VERSION}-windows-x86_64/
cp LICENSE tmp/uve-${VERSION}-windows-x86_64/
(cd tmp && zip -r -X ../dist/uve-${VERSION}-windows-x86_64.zip uve-${VERSION}-windows-x86_64)

# Generate checksums
cd dist
sha256sum *.tar.gz *.zip > SHA256SUMS.txt
cd ..

# Clean up
rm -rf tmp
rm README.txt
rm ${BINARY_NAME}-*

echo "=== Release packages created in dist/ directory ==="
echo "Version: ${VERSION}"
echo "Platforms: ${PLATFORMS[*]}"
