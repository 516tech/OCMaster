文档版本：V3.0
创建日期：2026-05-18
文档状态：设计稿
产品名称：超频大师（OCMaster）
产品定位：Windows 原生硬件采集工具 + 商家端超频建议平台

一、MVP 策略

单产物交付：
- C端：一个 exe，双模式。双击 → WinUI 3 GUI；`--cli scan` → 命令行
- S端：商家 Web 平台，本地 SQLite 开发 → 生产 MySQL + Docker Compose

S端零 Docker 依赖：`go run ./cmd/server` 直连 SQLite 文件启动。

二、用户角色

| 角色 | 需求 |
|------|------|
| 普通用户 | 双击 exe 启动 GUI，一键扫描硬件，导出或上传获取分享码 |
| 进阶用户 | 命令行 `ocmaster --cli scan --out json` 集成自动化 |
| 超频服务商 | 登录 Web 平台，输入分享码查询硬件，生成超频建议 PDF |

C端流程：双击 exe → GUI 一键扫描 → 表格展示 → TXT 导出 / 上传获取分享码。
CLI 流程：`ocmaster --cli scan [--out json]` → stdout 输出。

三、C端技术方案（单 exe 双模式）

- 一个产物：`ocmaster_0.0.1_windows_amd64.exe`（~100MB 单文件）
- GUI：WinUI 3 (NavigationView 四页) + Go DLL (嵌入资源，首次提取)
- CLI：`AttachConsole(-1)` 接管终端，复用 Go DLL 扫描逻辑
- Windows 10 1809+ / Windows 11，仅 x64
- Go 扫描器三平台源码 (`//go:build`)，macOS/Linux 用独立 Go bin

四、S端需求（不变）

- 手机号+密码注册(bcrypt)，管理员审核，JWT 登录
- 分享码查询硬件信息，历史记录保留 30 天
- 参考数据库：CPU/内存颗粒/散热 超频参数
- 超频建议模板填写 + chromedp 生成 PDF
- 商家个人信息 + 自定义风险模板

五、数据库策略

| 环境 | 数据库 | 启动方式 |
|------|--------|---------|
| 本地开发 | SQLite (文件) | `go run ./cmd/server` |
| 集成测试 | SQLite :memory: | go test 内存数据库 |
| 生产部署 | MySQL 8.0 | Docker Compose (backend+MySQL+Nginx) |

六、开发阶段

阶段〇：架构重构（C端取代 Electron）
- [x] 0.1 废弃 c-end/electron + Vue 3 代码（保留 Go DLL 逻辑参考）
- [x] 0.2 Go c-shared DLL 项目初始化（c-end/hardware-scanner/）
- [x] 0.3 C# WinUI 3 项目初始化（c-end/OCMaster.App/）
- [x] 0.4 实现 Go DLL C 导出接口：ScanAll + 分项扫描
- [x] 0.5 C# P/Invoke 调用 Go DLL，JSON 反序列化到 Model
- [x] 0.6 WinUI 3 页面：ScanPage + SettingsPage
- [x] 0.7 CLI 模式（oomaster scan / scan --out json）
- [x] 0.8 JSON 日志系统（Go zerolog-ready + C# ConfigService）
- [x] 0.9 S端后端 HardwareInfo schema 与 Go struct 一致（JSON tag 对齐）
- [x] 0.10 CI：Windows runner Go bin + DLL + JSON/Table/help 集成测试
- [x] 0.11 Linux 硬件扫描：/proc/cpuinfo + sysfs + dmidecode + lspci
- [x] 0.12 三平台统一构建：build.sh / build.ps1 → bin/ocmaster* + lib*

阶段一：基础设施（已完成，保留）
- [x] 1.1-1.5 S端后端 + 前端项目初始化

阶段二：C端 MVP（C# WinUI 3 + Go DLL/CLI）
- [x] 2.1 WinUI 3 主窗口 + 导航 (MainWindow + ScanPage + SettingsPage)
- [x] 2.2 Go DLL Windows 扫描 (wmic: cpu/baseboard/memorychip/videocontroller)
- [x] 2.3 Go DLL Linux 扫描 (/proc + sysfs + dmidecode + lspci)
- [x] 2.4 Go DLL macOS 扫描 (sysctl + system_profiler, Apple Silicon 实测)
- [x] 2.5 ScanPage: 表格展示 + Loading/Error/Empty + 骨架屏
- [x] 2.6 SettingsPage: JSON 文件持久化 (ConfigService)
- [x] 2.7 UploadPage: 扫描→上传→分享码→复制→删除，隐私提示
- [x] 2.7a AboutPage: 版本信息 + 技术栈
- [x] 2.7b MainWindow: NavigationView 四页导航
- [x] 2.8 TXT 导出 (FileSavePicker)
- [x] 2.9 CLI 三平台 bin + CI 集成测试 (JSON + Table + version)
- [x] 2.10 双模式合一: 双击→GUI, --cli scan→命令行, 单 exe 产物
- [x] 2.11 Go DLL 嵌入资源 → NativeLibrary.Load 首次提取

阶段三至六：S端后端 + 前端 + 集成（已完成，不变）
- [x] 3.1-3.10 S端后端（覆盖率 ~75%）
- [x] 4.1-4.9 S端前端
- [x] 5.3-5.5 打包 + 安全审查
- [x] 6.1-6.7 历史 Bug 修复
