<h1 align="center">超频大师 OCMaster</h1>

<p align="center">
    跨平台硬件信息采集工具 + 超频服务商建议平台
    <br>
    <strong>/'oʊvərklɒk 'mæstər/</strong>
</p>

## Features

<details>
<summary><strong>一键硬件扫描</strong></summary>

- CPU: 型号、核心数、线程数、基频
- 主板: 品牌、型号、芯片组、BIOS 版本
- 内存: 总容量、条数、频率、时序、颗粒类型、通道数
- 显卡: 型号、显存
- 电源: 额定功率 (SMBus 自动检测或手动输入)
- 散热: 风冷/水冷类型
- 扫描结果导出 TXT / JSON
</details>

<details>
<summary><strong>双模式运行</strong></summary>

- GUI 模式: Electron + Vue 3, 双击启动, 可视化操作
- CLI 模式: `ocmaster scan --out json`, 集成自动化工作流
- GUI 和 CLI 共享 Go 扫描器, 相同硬件检测逻辑
</details>

<details>
<summary><strong>上传分享码</strong></summary>

- 一键上传硬件信息到云端, 获取 6 位数字分享码
- 分享码 7 天有效, 用户可随时撤销删除
- 数据仅用于超频服务商提供建议, 到期自动清理
- 无 IP/地理位置采集, 上传完全可选
</details>

<details>
<summary><strong>超频建议平台 (S端)</strong></summary>

- 超频服务商通过分享码查询用户硬件配置
- 参考数据库: CPU 体质分 / 内存颗粒超频范围 / 散热器解热能力
- 模板化建议填写 + 一键 PDF 导出 (chromedp)
- 自定义风险提示模板, 商家个人资料管理
- 手机号+密码注册, bcrypt 加密, JWT 认证
</details>

<details>
<summary><strong>主题系统</strong></summary>

- 暗色模式 (默认) / 亮色模式 / 跟随系统
- CSS Variable 热切换, 无需重启
- 参考 ImHex ThemeManager 模式
- 主题偏好持久化到本地
</details>

<details>
<summary><strong>跨平台</strong></summary>

- Windows 10/11: Electron 单文件 exe (~70MB)
- macOS: .dmg (Apple Silicon arm64)
- Linux: .AppImage (amd64)
- Go CLI: 三平台静态编译 (~2MB, CGO_ENABLED=0)
</details>

## Architecture

OCMaster 采用 Core + Plugins 架构, 参考 ImHex 设计:

```
┌────────────────────────────────────────┐
│              Electron Shell             │
│  ┌──────────────────────────────────┐  │
│  │         Vue 3 Renderer           │  │
│  │  4 Pages + EventBus + Theme +    │  │
│  │  Toast + TaskManager + Pinia     │  │
│  └──────────────┬───────────────────┘  │
│                 │ IPC (invoke/handle)   │
│  ┌──────────────▼───────────────────┐  │
│  │        Electron Main Process     │  │
│  │  BrowserWindow + contextBridge   │  │
│  │  + scanner.ts + config.ts        │  │
│  │  + ipc-handlers.ts               │  │
│  └──────────────┬───────────────────┘  │
│                 │ child_process.spawn   │
│  ┌──────────────▼───────────────────┐  │
│  │     Go CLI (extraResources)       │  │
│  │  ocmaster scan --out json         │  │
│  │  → stdout → HardwareInfo JSON     │  │
│  └──────────────────────────────────┘  │
└────────────────────────────────────────┘

S端 (Go + Vue 3 SPA)
┌────────────────────┐
│  Nginx (port 80)    │
│  ├── /api/* → Go    │
│  └── /     → Vue    │
├────────────────────┤
│  Go Backend (8080)  │
│  chi / GORM /       │
│  zerolog / Wire DI  │
├────────────────────┤
│  MySQL 8.0 (3306)   │
└────────────────────┘
```

### Core (Electron Shell)

| 模块 | 技术 | 参考 ImHex |
|------|------|-----------|
| 窗口管理 | BrowserWindow + 状态持久化 | LayoutManager |
| IPC 通信 | contextBridge + ipcMain.handle | EventManager |
| 配置持久化 | JSON (userData) | ThemeManager |
| 主题切换 | CSS Variables + Pinia | ThemeManager |
| 任务管理 | AbortController + progress | TaskManager |
| 通知系统 | ElNotification 4s 过期 | Toast/Banner |

### Plugins (Go Scanner Backends)

| Backend | 平台 | 数据源 |
|---------|------|--------|
| scanner_windows | Windows | WMI (wmic) |
| scanner_darwin | macOS | sysctl + system_profiler |
| scanner_linux | Linux | /proc + sysfs + dmidecode + lspci |

### S端 (DDD 四层)

```
domain/ → application/ → infrastructure/ → interfaces/http/
```

Go/chi/GORM/Zerolog + Vue 3/Element Plus/Vite/Pinia

## Getting Started

**Requirements**

| 条件 | 要求 |
|------|------|
| OS | Windows 10 1809+ / macOS 12+ / Linux glibc 2.31+ |
| CPU | amd64 / arm64 (Apple Silicon) |
| RAM / Storage | ~50MiB / ~200MiB (GUI), ~10MiB / ~2MiB (CLI) |

**Install**
```
# C端 GUI — 下载 exe, 双击运行
https://github.com/516tech/OCMaster/releases/latest

# C端 CLI — 单文件静态编译
curl -LO <release>/ocmaster_0.0.1_$(uname -s)_$(uname -m)
chmod +x ocmaster_* && ./ocmaster_* scan --out json

# S端
cd s-end/backend && go run ./cmd/server   # SQLite 开发
docker compose up                          # MySQL + Nginx 生产
```

**Compiling**
```
cd c-end/hardware-scanner && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/ocmaster ./cmd/cli
cd c-end/electron && npm ci && npm run build && npx electron-builder --win portable
```
CI: GitHub Actions `build-c-end.yml`, windows-latest, push → auto Release.

**Contributing** — `scanner/scanner_<platform>.go` → `ScanAll() HardwareInfo` → PR.
Modules: `c-end/electron/` (GUI) | `c-end/hardware-scanner/` (CLI+sidecar) | `s-end/backend/` (API) | `s-end/frontend/` (SPA)

## License

MIT
