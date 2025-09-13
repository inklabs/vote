#!/bin/bash

set -ex

if [[ -z "$1" ]]; then
  echo "ERROR: Version not provided. Usage: $0 <version>"
  exit 1
fi

VERSION="$1"
OUTPUT_DIR="release"
BIN_NAME="api.vote.inklabs.dev"

mkdir -p "$OUTPUT_DIR"

TARGETS=(
  "linux arm"
#  "linux amd64"
#  "darwin amd64"
  "darwin arm64"
)

for t in "${TARGETS[@]}"; do
  read -r GOOS GOARCH <<< "$t"

  NAME="$BIN_NAME-$VERSION-$GOOS-$GOARCH"
  BINARY="$NAME.bin"
  ARCHIVE="$OUTPUT_DIR/$NAME.tar.gz"

  echo "Building $BINARY..."
  GOOS="$GOOS" GOARCH="$GOARCH" go build \
      -ldflags "-X github.com/inklabs/vote.Version=$VERSION" \
      -o "$OUTPUT_DIR/$BINARY" \
      cmd/httpapi/main.go

  echo "Packaging $ARCHIVE..."
  tar -C "$OUTPUT_DIR" -czf "$ARCHIVE" "$BINARY"

  rm "$OUTPUT_DIR/$BINARY"
done

echo "All builds completed. Check the '$OUTPUT_DIR' folder."
