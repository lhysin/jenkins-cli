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
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Failed to get latest version from GitHub API" >&2
        echo "This may be due to GitHub API rate limiting. Please try again later or download manually." >&2
        exit 1
    fi
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
echo "Binary: $BINARY"
echo "Installing Jenkins CLI to $DEST..."

# Check if we have permission to write to DEST
echo "Downloading..."
if ! curl -fsSL "https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}" -o "$DEST"; then
    echo "Error: Failed to download binary" >&2
    echo "URL: https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}" >&2
    
    # Check if it's a permission issue
    if [ ! -w "$(dirname "$DEST")" ]; then
        echo "" >&2
        echo "Permission denied. Try running with sudo:" >&2
        echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sudo sh" >&2
        echo "" >&2
        echo "Or install to a different location:" >&2
        echo "  DEST=\$HOME/bin/jenkins-cli curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sh" >&2
    fi
    exit 1
fi

chmod +x "$DEST"
echo "Installed successfully!"
echo "Run 'jenkins-cli${EXT} --help' to get started"
