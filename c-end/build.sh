#!/bin/bash
# OCMaster C端 统一构建脚本
# 产出: bin/ 目录下三平台可执行文件 + 动态库
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCANNER_DIR="$SCRIPT_DIR/hardware-scanner"
OUT_DIR="$SCRIPT_DIR/bin"

mkdir -p "$OUT_DIR"

echo "========================================="
echo " OCMaster C端 三平台构建"
echo "========================================="

cd "$SCANNER_DIR"

# 读取版本号
VERSION=$(cat "$SCRIPT_DIR/../VERSION" 2>/dev/null || echo "dev")
VERSION=$(echo "$VERSION" | tr -d '\n\r ')
LDFLAGS="-s -w -X main.version=$VERSION"
echo "    Version: $VERSION"

# ---------- 当前平台 CLI ----------
echo ""
echo "==> [1/3] 构建当前平台 CLI..."
HOST_OS=$(go env GOOS)
HOST_ARCH=$(go env GOARCH)
BIN_NAME="ocmaster"
[[ "$HOST_OS" == "windows" ]] && BIN_NAME="ocmaster.exe"

CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$OUT_DIR/$BIN_NAME" ./cmd/cli
echo "    -> $OUT_DIR/$BIN_NAME ($(du -h "$OUT_DIR/$BIN_NAME" | cut -f1)) (static)"

# ---------- 当前平台 DLL (需要 CGO) ----------
echo ""
echo "==> [2/3] 构建当前平台动态库 (CGO)..."
case "$HOST_OS" in
  darwin)
    DLL_NAME="libhardware_scanner.dylib"
    CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUT_DIR/$DLL_NAME" .
    ;;
  linux)
    DLL_NAME="libhardware_scanner.so"
    CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUT_DIR/$DLL_NAME" .
    ;;
  windows|mingw*)
    DLL_NAME="hardware_scanner.dll"
    CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUT_DIR/$DLL_NAME" .
    ;;
esac
echo "    -> $OUT_DIR/$DLL_NAME ($(du -h "$OUT_DIR/$DLL_NAME" | cut -f1))"

# ---------- 跨平台编译 ----------
echo ""
echo "==> [3/3] 跨平台编译..."

cross_build() {
  local os=$1 arch=$2 bin_suffix=$3 dll_ext=$4 dll_name=$5
  echo -n "    $os/$arch ... "

  # CLI binary (static)
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="$LDFLAGS" \
    -o "$OUT_DIR/ocmaster-${os}-${arch}${bin_suffix}" ./cmd/cli 2>/dev/null && \
    echo -n "bin✓ " || echo -n "bin✗ "

  # DLL
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -buildmode=c-shared \
    -o "$OUT_DIR/${dll_name}-${os}-${arch}${dll_ext}" . 2>/dev/null && \
    echo "dll✓" || echo "dll✗"
}

# macOS arm64 + amd64 (bin only, CGO cross is complex)
cross_build darwin arm64 "" ".dylib" "libhardware_scanner"
cross_build darwin amd64 "" ".dylib" "libhardware_scanner"

# Linux amd64 + arm64
cross_build linux amd64 "" ".so" "libhardware_scanner"
cross_build linux arm64 "" ".so" "libhardware_scanner"

# Windows amd64
cross_build windows amd64 ".exe" ".dll" "hardware_scanner"

echo ""
echo "========================================="
echo " 构建完成 — 产物:"
ls -lh "$OUT_DIR/" | awk '{print "   " $NF " (" $5 ")"}'
echo "========================================="
