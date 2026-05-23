<h1 align="center">OCMaster Bug Tracker</h1>

## Resolved

- [x] WinUI 3 XamlCompiler crash → 迁移到 giu
- [x] Electron 体积大 (70MB) → giu 13MB
- [x] c-end/electron/ 删除 (35 files)
- [x] CI Electron 构建复杂 → go build
- [x] scanner_test scanCPU/scanGPU → 改用 ScanAll()
- [x] CI covdata 错误 → 仅测试 scanner 包

## Pending

- [ ] GUI CGO CI 验证 (需要 MinGW + OpenGL)
- [ ] HTTP upload 实现 (uploadTab 目前是 stub)
- [ ] 三平台 GUI 手动测试
- [ ] macOS 代码签名
