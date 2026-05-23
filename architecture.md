<h1 align="center">OCMaster Architecture</h1>

<p align="center">
    Core (Electron Shell) + Plugins (Go Scanner Backends) + S端 (DDD API)
    <br>
    <strong>参考 ImHex 设计模式</strong>
</p>

## Architecture

```
C端 Electron GUI (~70MB)
┌──────────────────────────────────────────┐
│  Vue 3 Renderer  ◄── IPC ──►  Main       │
│  ┌────────────┐              ┌──────────┐ │
│  │ ScanPage   │  scan:all    │ scanner  │ │
│  │ UploadPage │  config:*/   │ config   │ │
│  │ SettingsPg │  export:txt  │ ipc-mgr  │ │
│  │ AboutPage  │  window:*    │ window   │ │
│  └────────────┘              └────┬─────┘ │
│                                   │spawn  │
│  Pinia Stores:                    │       │
│  scan / settings / theme          │       │
│                                   │       │
│  Composables:              ┌──────▼──────┐
│  EventBus / useTheme       │  Go CLI     │
│  / useToast                │  (sidecar)  │
│                            └─────────────┘
└──────────────────────────────────────────┘

S端 (Docker Compose)
┌─────────┐     ┌──────────────┐     ┌─────────┐
│  Nginx  │────►│ Go Backend   │────►│  MySQL  │
│  :80    │     │ chi/GORM     │     │  :3306  │
│  / → SPA│     │ zerolog/Wire │     │         │
└─────────┘     └──────────────┘     └─────────┘
                      │
              ┌───────▼───────┐
              │  Vue 3 SPA    │
              │  Element Plus │
              │  Vite / Pinia │
              └───────────────┘
```

## File Tree

```
OCMaster/
├── c-end/
│   ├── electron/                         # Electron Shell (Core)
│   │   ├── electron/
│   │   │   ├── main.ts                   # BrowserWindow + 窗口状态持久化
│   │   │   ├── preload.ts                # contextBridge API 暴露
│   │   │   ├── scanner.ts                # spawn Go CLI → stdout JSON
│   │   │   ├── config.ts                 # userData/config.json 读写
│   │   │   └── ipc-handlers.ts           # ipcMain.handle 路由注册
│   │   ├── src/                          # Vue 3 Renderer
│   │   │   ├── App.vue                   # 导航 + 主题切换按钮
│   │   │   ├── views/                    # 4 页面组件
│   │   │   ├── stores/                   # Pinia (scan/settings/theme)
│   │   │   ├── composables/              # useEventBus/useTheme/useToast
│   │   │   ├── api/client.ts             # Axios HTTP 客户端
│   │   │   ├── types/                    # HardwareInfo 类型 + electron.d.ts
│   │   │   └── router/                   # Vue Router (Hash 模式)
│   │   ├── electron-builder.yml          # Win/Mac/Linux 打包配置
│   │   ├── electron.vite.config.ts       # Main/Preload/Renderer 构建
│   │   └── package.json
│   ├── hardware-scanner/                 # Go Scanner (Plugins)
│   │   ├── cmd/cli/main.go               # CLI 入口 (sidecar + 独立 CLI)
│   │   ├── scanner/
│   │   │   ├── types.go                  # HardwareInfo JSON struct
│   │   │   ├── scanner_darwin.go         # macOS: sysctl + system_profiler
│   │   │   ├── scanner_linux.go          # Linux: /proc + dmidecode + lspci
│   │   │   ├── scanner_windows.go        # Windows: WMI (wmic)
│   │   │   └── scanner_test.go           # 集成测试 (80.7% 覆盖)
│   │   └── go.mod
│   ├── bin/                              # 构建产物 (gitignore)
│   ├── build.sh / build.ps1
├── s-end/
│   ├── backend/                          # Go DDD API
│   │   ├── domain/                       # Entity + Repository 接口
│   │   ├── application/                  # Service 用例编排
│   │   ├── infrastructure/               # GORM + JWT + Chromedp PDF
│   │   └── interfaces/http/              # chi Handlers + Middleware
│   ├── frontend/                         # Vue 3 SPA
│   │   └── src/views/                    # 8 商家页面
├── .github/workflows/build-c-end.yml     # CI: Windows → Release
├── VERSION                               # 0.0.1
└── docker-compose.yml                    # MySQL + Backend + Nginx
```

## IPC & Data Flow

Renderer (Vue) ↔ Preload (contextBridge) ↔ Main (ipcMain)

| Channel | Direction | Returns |
|---------|-----------|---------|
| `scan:all` | invoke | `ScanResult { success, data?, error? }` |
| `config:load` | invoke | `AppConfig` JSON |
| `config:save` | invoke | `{ ok: true }` |
| `export:txt` | invoke | `ExportResult` |
| `app:version` | invoke | `string` |
| `window:*` | invoke | tab / state persistence |

Sidecar: `spawn('ocmaster', ['scan', '--out', 'json'])` → stdout `HardwareInfo` JSON → exit 0.

```
Scan Flow:
  User Click → ScanPage.scan() → AbortController
    → ipcRenderer.invoke('scan:all') → scanner.ts
    → spawn('ocmaster', ['scan', '--out', 'json'])
    → stdout JSON → HardwareInfo → Pinia store
    → events.scanCompleted.post() → Toast

Upload Flow:
  UploadPage → scanAll() → JSON → api/client.uploadHardware()
    → POST /api/v1/hardware/upload → shareCode → 复制/删除

Export Flow:
  ScanPage → window.electronAPI.exportTxt(content)
    → dialog.showSaveDialog → writeFileSync → Toast
```

## CI Pipeline

```
Push main → build (windows-latest)
  ├── Setup Go 1.22 + Node 20
  ├── Build Go CLI (CGO_ENABLED=0) → c-end/bin/ocmaster.exe
  ├── go test ./scanner/... -cover
  ├── CLI JSON integration test (scan --out json → validate fields)
  ├── Sync VERSION → package.json
  ├── npm ci → electron-vite build → electron-builder --win portable
  └── Rename → ocmaster_0.0.1_windows_amd64.exe → Upload artifact

Release (non-PR, needs build + s-end-test)
  └── Download artifact → softprops/action-gh-release@v2 (tag: v0.0.1)
```

## Design Decisions

| 决策 | 选择 | 原因 |
|------|------|------|
| UI 框架 | Electron + Vue 3 | 跨平台, 无 XAML 编译问题, HMR 开发体验 |
| Go 集成 | Child Process (sidecar) | Go CLI 已跨平台, 进程隔离, 独立可测试 |
| 状态管理 | Pinia | 官方推荐, Composition API 风格 |
| UI 组件库 | Element Plus | 完整中文支持, dark mode, 对标 WinUI 控件 |
| 主题 | CSS Variables + Pinia | 热切换, 参考 ImHex ThemeManager |
| 事件通信 | Typed EventBus | Event/Request 分离, 参考 ImHex EventManager |
| 打包 | electron-builder (portable) | 单文件 exe, 无安装器依赖 |
| S端数据库 | SQLite → MySQL | 零配置开发, Docker Compose 生产 |
