<h1 align="center">OCMaster Evaluate</h1>

<p align="center">本轮迭代: Electron → Go + giu (Dear ImGui) 迁移</p>

## Goals

| # | 目标 | 结果 |
|---|------|------|
| 1 | 删除 Electron + Vue 3 项目 | ✅ 35 files removed |
| 2 | Go + giu GUI 实现 (4 Tab) | ✅ cmd/gui/ (main/pages/theme) |
| 3 | 单二进制, 双模式 (GUI + CLI) | ✅ ocmaster scan → CLI |
| 4 | go build 编译通过 | ✅ 13MB binary |
| 5 | go test 通过 | ✅ 80.7% coverage |

## Compare

| 指标 | Electron | Go + giu |
|------|----------|----------|
| 二进制大小 | ~70MB | **~13MB** |
| 依赖 | Node.js + npm + Chromium | **Go + CGO** |
| 构建工具 | electron-vite + electron-builder | **go build** |
| UI 范式 | Web (Vue 3 + Element Plus) | **原生 ImGui** |
| 源码文件 | 35 files | **3 files** |
| 测试 | vitest 25 tests | go test 80.7% |
