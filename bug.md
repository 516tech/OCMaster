<h1 align="center">OCMaster Bug Tracker</h1>

## Resolved

- [x] WinUI 3 XamlCompiler.exe crash → 迁移到 Electron
- [x] WinUI 3 NU1603 版本警告 → 不再使用 NuGet
- [x] Windows-only GUI → Electron 三平台
- [x] EnableXamlCompilation=false 不生效 → 不再使用 WinUI 3
- [x] c-end/OCMaster.App/ 旧代码 → 已删除 (18 files)
- [x] build-windows.yml 旧 CI → 替换为 build-c-end.yml
- [x] Go DLL c-shared main.go → 已删除 (sidecar 替代)
- [x] scanner_test.go scanCPU/scanGPU undefined → 改用 ScanAll()
- [x] CI go test cmd/cli covdata 错误 → 仅测试 scanner 包
- [x] electron-builder entry main/index.js → 改为 main/main.js
- [x] electron-builder build/ → dist/ 路径
- [x] build.sh DLL 构建段 → 已清理

## Pending

- [ ] 三平台 GUI 手动测试 (需真实硬件)
- [ ] macOS 代码签名 (需 Apple 开发者账号)
