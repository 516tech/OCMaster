文档版本：V3.0
创建日期：2026-05-18
文档状态：设计稿
所属项目：超频大师（OCMaster）

一、评估维度

| 维度 | 通过标准 | 当前状态 |
|------|---------|---------|
| 功能完整性 | 阶段〇 全部勾选 | 11/12 (WMI 待实现) |
| Go 测试 | go test 全PASS | scanner 80.7% 覆盖 |
| S端测试 | go test 全PASS | 15/15 PASS (81 tests) |
| 三平台 bin | win/mac/linux 编译 | 7 文件 1.9-2.1MB |
| 三平台 dll | c-shared 编译 | macOS .dylib 验证通过 |
| CI Windows 集成 | JSON+Table+help | 3-step 测试通过 |

二、本迭代新增

- 统一构建脚本: build.sh (bash) + build.ps1 (PowerShell)
- 三平台 CLI bin: ocmaster / ocmaster.exe
- Linux 扫描器: /proc/cpuinfo + /sys/class/dmi/ + dmidecode + lspci
- CI go-build job: Go setup → Build CLI → Build DLL → 3 集成测试 → Upload
- CI s-end-test job: go test ./... -cover

三、构建产物 (c-end/bin/)

| 文件 | 平台 | 大小 |
|------|------|------|
| ocmaster | macOS arm64 | 2.0M |
| ocmaster-windows-amd64.exe | Windows x64 | 1.9M |
| ocmaster-linux-amd64 | Linux x64 | 2.1M |
| libhardware_scanner.dylib | macOS arm64 | 2.1M |

四、待完成

- [ ] Windows WMI 扫描器实现
- [ ] WinUI 3 dotnet build 验证
