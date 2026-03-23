#!/bin/sh
set -e

REPO="lhysin/jenkins-cli"

detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)        echo " unsupported" && exit 1 ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)       echo " unsupported" && exit 1 ;;
    esac
}

OS=$(detect_os)
ARCH=$(detect_arch)
BINARY="jenkins-cli_${OS}_${ARCH}"
EXT=""
[ "$OS" = "windows" ] && EXT=".exe"

DEST="${DEST:=/usr/local/bin/jenkins-cli$EXT}"

echo "Detected: $OS/$ARCH"
echo "Installing Jenkins CLI to $DEST..."
curl -fsSL "https://github.com/${REPO}/releases/latest/download/${BINARY}$EXT" -o "$DEST"
chmod +x "$DEST"
echo "Installed!"
echo "Run 'jenkins-cli$EXT --help' to get started"