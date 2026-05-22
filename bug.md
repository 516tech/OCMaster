文档版本：V4.1
更新日期：2026-05-22
所属项目：超频大师（OCMaster）

一、已解决 (Electron 迁移)

- [x] BUG-115: WinUI 3 XamlCompiler.exe crash → 迁移到 Electron
- [x] BUG-116: WinUI 3 NU1603 版本警告 → 不再使用 NuGet
- [x] BUG-117: Windows-only GUI → Electron 三平台
- [x] BUG-118: EnableXamlCompilation=false 不生效 → 不再使用 WinUI 3
- [x] BUG-119: 旧 c-end/OCMaster.App/ 代码已删除
- [x] BUG-120: 旧 build-windows.yml 已替换为 build-c-end.yml

二、Go 扫描器 (已解决)

- [x] BUG-100-105: Go DLL 三平台 + C# P/Invoke → Electron sidecar 替代
- [x] BUG-106-114: CI/静态编译/nulls/help 已修复

三、待解决

- [ ] 三平台 GUI 手动测试 (需真实硬件)
- [ ] macOS 代码签名 (需 Apple 开发者账号)
- [ ] CI 首次运行验证
