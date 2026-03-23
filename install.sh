#!/bin/sh
set -e

REPO="lhysin/jenkins-cli"

detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)        echo "unsupported" && exit 1 ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)       echo "unsupported" && exit 1 ;;
    esac
}

get_latest_version() {
    VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    echo "$VERSION"
}

OS=$(detect_os)
ARCH=$(detect_arch)
VERSION=${VERSION:-$(get_latest_version)}
VERSION_NO_V=$(echo "$VERSION" | sed 's/^v//')
EXT=""
[ "$OS" = "windows" ] && EXT=".exe"
BINARY="jenkins-cli_${VERSION_NO_V}_${OS}_${ARCH}${EXT}"
DEST="${DEST:-/usr/local/bin/jenkins-cli${EXT}}"

echo "Detected: $OS/$ARCH"
echo "Latest version: $VERSION"
echo "Installing Jenkins CLI to $DEST..."
curl -fsSL "https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}" -o "$DEST"
chmod +x "$DEST"
echo "Installed!"
echo "Run 'jenkins-cli${EXT} --help' to get started"