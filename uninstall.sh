#!/bin/sh
set -e

DEST="${DEST:-/usr/local/bin/jenkins-cli}"

echo "Uninstalling Jenkins CLI..."
if [ -f "$DEST" ]; then
    rm "$DEST"
    echo "Removed $DEST"
else
    echo "Jenkins CLI not found at $DEST"
fi
echo "Done"