文档版本：V4.0
创建日期：2026-05-18
更新日期：2026-05-22
文档状态：设计稿
产品名称：超频大师（OCMaster）
产品定位：跨平台硬件采集工具 + 商家端超频建议平台

一、MVP 策略

单产物交付：
- C端：Electron + Vue 3 跨平台 GUI + Go CLI sidecar
- S端：商家 Web 平台，本地 SQLite 开发 → 生产 MySQL + Docker Compose
- S端零 Docker 依赖：`go run ./cmd/server` 直连 SQLite 文件启动

二、用户角色

| 角色 | 需求 |
|------|------|
| 普通用户 | 双击启动 GUI，一键扫描硬件，导出或上传获取分享码 |
| 进阶用户 | 命令行 `ocmaster scan --out json` 集成自动化 |
| 超频服务商 | 登录 Web 平台，输入分享码查询硬件，生成超频建议 PDF |

三、C端技术方案（Electron + Go Sidecar）

- GUI：Electron 33 + Vue 3 + Element Plus, 4 页面 (Scan/Upload/Settings/About)
- 硬件扫描：Go CLI 作为子进程，通过 stdout 获取 JSON 结果
- 跨平台：Windows/macOS/Linux GUI，统一代码库
- 打包：electron-builder (exe + dmg + AppImage)
- Go 扫描器三平台静态编译 (CGO_ENABLED=0)

四、S端需求（不变）

- 手机号+密码注册(bcrypt)，管理员审核，JWT 登录
- 分享码查询硬件信息，历史记录保留 30 天
- 参考数据库：CPU/内存颗粒/散热 超频参数
- 超频建议模板填写 + chromedp 生成 PDF

五、开发阶段

阶段〇：架构重构（Electron 迁移）
- [x] 0.1 废弃 c-end/OCMaster.App/ WinUI 3 代码 (保留 Go DLL 逻辑)
- [x] 0.2 Electron + Vue 3 项目初始化 (c-end/electron/)
- [x] 0.3 Electron 主进程：BrowserWindow + preload + contextBridge
- [x] 0.4 Go Sidecar 扫描器集成 (scanner.ts + IPC handlers)
- [x] 0.5 Vue 3 页面：ScanPage + UploadPage + SettingsPage + AboutPage
- [x] 0.6 ApiClient (TypeScript) + ConfigService (Electron userData)
- [x] 0.7 electron-builder 配置 (win/mac/linux)
- [x] 0.8 CI Matrix: Go 三平台构建 + Electron 打包 + Release
- [ ] 0.9 npm run build 验证通过 (electron-vite)
- [ ] 0.10 三平台手动测试 GUI

阶段一至六：S端 + 集成（已完成，不变）

六、CI Matrix (参考 ImHex)

| Job | Runner | 产物 |
|-----|--------|------|
| go-build | win/mac/linux 并行 | Go CLI + JSON 集成测试 |
| electron-build | win/mac/linux 并行 | Electron 打包产物 |
| s-end-test | ubuntu-latest | go test -cover |
| release | ubuntu-latest | GitHub Release |
