#!/bin/bash
# OCMaster C端 构建脚本
# 产出: bin/ Go CLI (三平台静态) + Go GUI (CGO)
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCANNER_DIR="$SCRIPT_DIR/hardware-scanner"
OUT_DIR="$SCRIPT_DIR/bin"
mkdir -p "$OUT_DIR"

cd "$SCANNER_DIR"

VERSION=$(cat "$SCRIPT_DIR/../VERSION" 2>/dev/null || echo "dev")
VERSION=$(echo "$VERSION" | tr -d '\n\r ')
LDFLAGS="-s -w -X main.version=$VERSION"
echo "OCMaster C端 构建 — Version: $VERSION"

# --- Go CLI (三平台静态) ---
echo "==> Go CLI (静态编译)..."
for target in "darwin/arm64" "darwin/amd64" "linux/amd64" "linux/arm64" "windows/amd64"; do
  os="${target%/*}"; arch="${target#*/}"; ext=""
  [[ "$os" == "windows" ]] && ext=".exe"
  echo -n "  $os/$arch ... "
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="$LDFLAGS" \
    -o "$OUT_DIR/ocmaster_${VERSION}_${os}_${arch}${ext}" ./cmd/cli 2>/dev/null && echo "✓" || echo "✗"
done

# --- Go GUI (当前平台, CGO) ---
echo "==> Go GUI (CGO)..."
HOST_OS=$(go env GOOS)
HOST_ARCH=$(go env GOARCH)
GUI_EXT=""
[[ "$HOST_OS" == "windows" ]] && GUI_EXT=".exe"
CGO_ENABLED=1 go build -ldflags="$LDFLAGS" \
  -o "$OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH}${GUI_EXT}" ./cmd/gui 2>/dev/null && \
  echo "  -> $OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH}${GUI_EXT} ($(du -h "$OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH}${GUI_EXT}" | cut -f1))" || \
  echo "  (跳过: 需要 CGO + OpenGL)"

echo "==> 完成"
ls -lh "$OUT_DIR/" | awk '{print "  " $NF " (" $5 ")"}'
