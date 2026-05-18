#!/bin/bash
# 单源版本控制：从 VERSION 文件同步到所有模块
# 用法: ./scripts/sync-version.sh

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VER=$(cat "$ROOT/VERSION" | tr -d '[:space:]')
echo "Syncing version: $VER"

# C端
cd "$ROOT/c-end"
npm pkg set version="$VER" --json > /dev/null 2>&1 || true

# S端前端
cd "$ROOT/s-end/frontend"
npm pkg set version="$VER" --json > /dev/null 2>&1 || true

# electron-builder 产物名（通过 env 传递）
echo "VERSION=$VER" > "$ROOT/.version.env"

echo "Version synced to $VER"
