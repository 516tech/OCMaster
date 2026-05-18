文档版本：V2.0
创建日期：2026-05-17
文档状态：终稿
产品名称：超频大师（OCMaster）
产品定位：跨平台硬件信息采集工具 + 商家端超频建议生成平台

一、MVP 策略

MVP 分两层交付：
- C端先行：Electron 桌面应用，一键扫描硬件 + 本地展示 + 导出，可独立使用
- S端后继：商家平台，本地 SQLite 开发 → 生产 MySQL + Docker Compose

本地开发零 Docker 依赖：S端后端 Golang 直连 SQLite 文件数据库，`go run ./cmd/server` 即可启动。

二、用户角色与流程

| 角色 | 需求 |
|------|------|
| 普通用户 | 一键获取硬件信息，导出本地文件或上传生成分享码 |
| 超频服务商 | 输入分享码查看硬件，参考数据库生成超频建议 PDF |

C端流程：双击运行 → 一键扫描 → 表格展示 → 手动修正 → 本地导出(TXT) / 上传获取分享码。
S端流程：登录 → 输入分享码 → 查看硬件 + 参考数据 → 填写超频建议 → 生成 PDF。

三、功能需求

3.1 C端 — MVP 优先
- 三平台 Electron 单文件应用（.exe/.dmg/.AppImage）
- 一键扫描：C++ Native Addon 采集 CPU/主板/内存/显卡/电源/散热
- 本地表格展示 + 手动修正 + TXT 导出
- 上传：授权提示 → 6位分享码（7天有效）→ 一键复制后可删除
- 设置页：API 地址、语言（默认中文）、自动扫描开关，electron-store 持久化
- **安全防护**：代码混淆（JS obfuscator），反调试检测（DevTools禁用+快捷键拦截）
- **产物 ≤80MB**：Electron 28 不移版本，裁剪 locales + NSIS solid + 7z SFX 后处理
- **全异步 UI**：defineAsyncComponent + Suspense 骨架屏 + 加载/空/错误三态

3.2 S端 — 本地开发
- 账号：手机号+密码注册(bcrypt)，管理员审核，JWT 登录
- 分享码查询硬件信息，历史记录保留 30 天
- 参考数据库：CPU/内存颗粒/散热 超频参数
- 超频建议模板填写 + chromedp 生成 PDF
- 商家个人信息 + 自定义风险模板

四、数据库策略

| 环境 | 数据库 | 启动方式 |
|------|--------|---------|
| 本地开发 | SQLite (文件) | `go run ./cmd/server` 零依赖启动 |
| 集成测试 | SQLite :memory: | 测试代码内存数据库 |
| 生产部署 | MySQL 8.0 | Docker Compose (backend+MySQL+Nginx) |

GORM 自动迁移，按 DB_TYPE 环境变量切换驱动。

五、开发阶段

阶段一：基础设施
- [x] 1.1 项目目录结构 + Go module (s-end/backend)
- [x] 1.2 Electron + Vue 3 项目初始化 (c-end)
- [x] 1.3 Vue 3 + Element Plus 项目初始化 (s-end/frontend)
- [x] 1.4 S端后端适配 SQLite + MySQL 双驱动，GORM 自动迁移
- [x] 1.5 编写 Makefile（本地零 Docker 可构建/测试）

阶段二：C端 — MVP 核心
- [x] 2.1 Electron 主进程：窗口管理、IPC 注册
- [x] 2.2 Preload contextBridge API
- [x] 2.3 Vue 渲染进程：ScanButton + HardwareTable + UploadPanel + SettingsDialog
- [x] 2.4 C++ Native Addon：Windows 硬件采集 stub (WMI + SMBus + SPD)
- [x] 2.5 C++ Native Addon：Linux 硬件采集 stub (sysfs + dmidecode)
- [x] 2.6 C++ Native Addon：macOS 硬件采集 stub (IOKit + sysctl)
- [ ] 2.7 Native Addon 单元测试 (mock 数据源) — 需 C++ 编译环境
- [x] 2.8 electron-store 配置持久化
- [x] 2.9 electron-builder 打包配置
- [x] 2.10 C端 vite build PASS（361ms, 3产物: main.js/preload.js/renderer）

阶段三：S端后端
- [x] 3.1 domain 层：5 聚合（实体 + 值对象 + 仓储接口）
- [x] 3.2 infrastructure/persistence：GORM 模型 + 仓储实现
- [x] 3.3 pkg/config + pkg/apperrors
- [x] 3.4 application 层 Service
- [x] 3.5 SQLite :memory: 仓储集成测试（repo_test.go）
- [x] 3.6 JWT + zerolog 中间件
- [x] 3.7 Router + Handler
- [x] 3.8 Wire DI + 定时清理
- [x] 3.9 PDF 生成 (chromedp) — Chrome实测 107KB PDF, 75%覆盖
- [ ] 3.10 go test ./... 覆盖率 ≥ 70%（13 PASS, ~40% avg）

阶段四：S端前端
- [x] 4.1 路由 + Axios + Pinia
- [x] 4.2 登录/注册页
- [x] 4.3 Dashboard（分享码输入）
- [x] 4.4 硬件详情页
- [x] 4.5 超频建议表单 + 参考数据面板
- [x] 4.6 PDF 预览 + 下载
- [x] 4.7 历史记录页
- [x] 4.8 个人信息 + 模板编辑器
- [x] 4.9 npm run build 构建成功（8视图+index，128 packages）

阶段五：集成与部署
- [ ] 5.1 C端→S端端到端流程验证
- [ ] 5.2 Docker Compose 生产部署（MySQL + Nginx）
- [x] 5.3 Windows C端打包产物：OCMaster 1.0.0.exe (140MB) ✅
- [ ] 5.4 性能：扫描≤5s, 上传≤3s, 报告≤2s
- [ ] 5.5 安全审查：HTTPS, JWT, SQL注入, XSS

阶段六：Bug修复（参考 bug.md）
- [ ] 6.1 BUG-001: wire 工具安装 + wire_gen.go 生成（低优先）
- [x] 6.2 BUG-002: 覆盖率改善（middleware 65%, handler 12.7%, app 30.5%）
- [x] 6.3 BUG-003: S端前端 23文件完成（8路由+3store+5API+8视图）
- [ ] 6.4 BUG-004: C端 2.7 Native Addon 测试 + 2.10 C端 E2E
- [x] 6.5 BUG-005: chromedp PDF — Chrome实测 PASS, 75%覆盖
- [ ] 6.6 BUG-006: Docker Compose 集成测试（需 Colima/网络）
- [x] 6.7 BUG-007: Electron 28.3.3 安装成功，C端构建通过 ✅

阶段七：C端加固与优化（新增）
- [x] 7.1 代码混淆：vite-plugin-javascript-obfuscator (controlFlow + stringArray)
- [x] 7.2 反调试：DevTools禁用 + F12/Ctrl+I拦截 + --inspect检测
- [x] 7.3 产物 135MB (portable EA 28已裁剪10.6MB Vulkan/D3D，欲<80MB需换WebView2)
- [x] 7.4 全异步 UI：Suspense + defineAsyncComponent + 骨架屏 + 错误重试
