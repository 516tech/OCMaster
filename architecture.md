# OCMaster 架构设计 V4.2

文档版本：V4.2 | 更新日期：2026-05-22 | 文档状态：设计稿

## 一、技术选型

| 模块 | 技术 |
|------|------|
| C端 GUI | Electron 33 + Vue 3 + TypeScript + Element Plus |
| C端硬件扫描 | Go CLI (sidecar, stdout JSON) |
| S端后端 | Golang / chi / GORM v2 / zerolog (DDD 四层) |
| S端前端 | Vue 3 / Element Plus / Vite / Pinia |
| 本地数据库 | SQLite (:memory: 测试 / 文件 开发) |
| 生产数据库 | MySQL 8.0 + Docker Compose |

## 二、C端架构

```
ocmaster_0.0.1_windows_amd64.exe (~70MB)
├── Electron Main Process
│   ├── BrowserWindow + contextBridge
│   ├── scanner.ts: spawn('ocmaster', ['scan', '--out', 'json'])
│   ├── config.ts: JSON 持久化 (userData)
│   └── ipc-handlers.ts: IPC 路由
├── Vue 3 Renderer (4 页面)
│   ├── ScanPage: 扫描 + 进度 + 取消 + TXT 导出
│   ├── UploadPage: 上传 + 分享码 + 复制/删除
│   ├── SettingsPage: API URL + 语言 + 主题 + 自动扫描
│   └── AboutPage: 版本 + 技术栈
├── UI 基础设施 (参考 ImHex)
│   ├── EventBus (Event/Request pub/sub)
│   ├── Theme (dark/light/system CSS var 热切换)
│   └── Toast (4s 自动过期 success/error/warning/info)
└── Go CLI (extraResources/scanner/)
    └── ocmaster scan --out json → stdout → HardwareInfo
```

## 三、目录结构

```
OCMaster/
├── c-end/
│   ├── electron/                    # Electron + Vue 3
│   │   ├── package.json
│   │   ├── electron-builder.yml     # Windows portable + NSIS
│   │   ├── electron.vite.config.ts
│   │   ├── electron/
│   │   │   ├── main.ts              # BrowserWindow + 窗口状态持久化
│   │   │   ├── preload.ts           # contextBridge API
│   │   │   ├── scanner.ts           # Go sidecar spawn
│   │   │   ├── config.ts            # 配置读写 (JSON)
│   │   │   └── ipc-handlers.ts      # IPC + last-tab 持久化
│   │   └── src/                     # Vue 3 渲染进程
│   │       ├── App.vue              # 导航 + 主题切换
│   │       ├── views/               # 4 个页面
│   │       ├── stores/              # Pinia (scan, settings, theme)
│   │       ├── composables/         # useEventBus, useTheme, useToast
│   │       ├── api/client.ts        # HTTP 客户端
│   │       ├── types/               # TypeScript 类型
│   │       └── router/index.ts      # Vue Router
│   ├── hardware-scanner/            # Go 扫描器
│   │   ├── cmd/cli/main.go          # CLI 入口 (sidecar)
│   │   ├── scanner/                 # 三平台实现 (WMI/sysfs/IOKit)
│   │   └── go.mod / go.sum
│   ├── bin/                         # 构建产物 (gitignore)
│   ├── build.sh / build.ps1
├── s-end/
│   ├── backend/                     # Go DDD 四层
│   └── frontend/                    # Vue 3 SPA
├── .github/workflows/
│   └── build-c-end.yml              # Windows build → Release
├── VERSION
└── docker-compose.yml
```

## 四、数据流

```
用户点击「开始扫描」
  → Vue ScanPage → window.electronAPI.scanAll()
  → IPC 'scan:all' → scanner.ts → spawn('ocmaster', ['scan', '--out', 'json'])
  → stdout JSON → HardwareInfo → Pinia store
  → events.scanCompleted.post() → Toast 通知

上传分享码:
  → UploadPage → scanAll() → JSON → api/client.uploadHardware()
  → POST /api/v1/hardware/upload → shareCode → 复制/删除
```

## 五、CI + 产物

```
build (windows-latest)
  ├── Go CLI 静态编译 + JSON 集成测试
  ├── npm ci → electron-vite build → electron-builder --win portable
  └── ocmaster_0.0.1_windows_amd64.exe → artifact

release (non-PR)
  └── GitHub Release (tag: v0.0.1)
```

## 六、UI 架构 (ImHex 模式)

| ImHex 模式 | Vue 3 实现 | 文件 |
|-----------|-----------|------|
| EventManager | Typed EventBus (Event/Request 分离) | `composables/useEventBus.ts` |
| ThemeManager | CSS Variables + Pinia dark/light/system | `composables/useTheme.ts`, `stores/theme.ts` |
| Toast/Banner | ElNotification 4s 过期 | `composables/useToast.ts` |
| TaskManager | AbortController + progress 0-100 | `stores/scan.ts` |
| LayoutManager | window-state.json + last-tab.json | `electron/main.ts`, `ipc-handlers.ts` |

## 七、S端 (不变)

DDD 四层: domain → application → infrastructure → interfaces/http
Go/chi/GORM/Zerolog | Vue 3 SPA | Docker Compose (MySQL + backend + Nginx)
