#!/bin/bash
# 单源版本控制：从 VERSION 文件同步到所有模块
# 用法: ./scripts/sync-version.sh

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VER=$(cat "$ROOT/VERSION" | tr -d '[:space:]')
echo "Syncing version: $VER"

# S端前端
cd "$ROOT/s-end/frontend"
npm pkg set version="$VER" --json > /dev/null 2>&1 || true
echo "  s-end/frontend/package.json → $VER"

# C端 Go CLI: build.sh / build.ps1 自动读取 VERSION 文件嵌入 ldflags
echo "  C端 Go build 脚本自动读取 VERSION"

echo "Version synced to $VER"
