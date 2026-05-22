#!/bin/bash
# OCMaster C端 构建脚本
# 产出: bin/ 目录下 Go CLI (三平台静态) + Electron (可选)
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCANNER_DIR="$SCRIPT_DIR/hardware-scanner"
OUT_DIR="$SCRIPT_DIR/bin"

mkdir -p "$OUT_DIR"

echo "========================================="
echo " OCMaster C端 构建"
echo "========================================="

cd "$SCANNER_DIR"

VERSION=$(cat "$SCRIPT_DIR/../VERSION" 2>/dev/null || echo "dev")
VERSION=$(echo "$VERSION" | tr -d '\n\r ')
LDFLAGS="-s -w -X main.version=$VERSION"
echo "    Version: $VERSION"

# ---------- 当前平台 CLI ----------
echo ""
echo "==> [1/2] 构建当前平台 CLI..."
HOST_OS=$(go env GOOS)
HOST_ARCH=$(go env GOARCH)
CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH}" ./cmd/cli
echo "    -> $OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH} ($(du -h "$OUT_DIR/ocmaster_${VERSION}_${HOST_OS}_${HOST_ARCH}" | cut -f1)) (static)"

# ---------- 跨平台 CLI ----------
echo ""
echo "==> [2/2] 跨平台编译..."

for target in "darwin/arm64" "darwin/amd64" "linux/amd64" "linux/arm64" "windows/amd64"; do
  os="${target%/*}"
  arch="${target#*/}"
  ext=""
  [[ "$os" == "windows" ]] && ext=".exe"
  echo -n "    $os/$arch ... "
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="$LDFLAGS" \
    -o "$OUT_DIR/ocmaster_${VERSION}_${os}_${arch}${ext}" ./cmd/cli 2>/dev/null && \
    echo "✓" || echo "✗"
done

echo ""
echo "========================================="
echo " 构建完成 — 产物:"
ls -lh "$OUT_DIR/" | awk '{print "   " $NF " (" $5 ")"}'

# ---------- Electron (可选) ----------
echo ""
ELECTRON_DIR="$SCRIPT_DIR/electron"
if [ -d "$ELECTRON_DIR/node_modules" ]; then
  echo "==> Electron 构建..."
  cd "$ELECTRON_DIR"
  npm run build 2>/dev/null && echo "    Electron build ✓" || echo "    Electron build ✗ (继续)"
else
  echo "    (Electron 跳过: cd c-end/electron && npm install)"
fi
echo "========================================="
