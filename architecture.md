<h1 align="center">OCMaster Architecture</h1>

<p align="center">
    Go + giu (Dear ImGui) — 跟 ImHex 相同的立即模式 GUI
    <br>
    <strong>单二进制, 双模式 (GUI + CLI), ~13MB</strong>
</p>

## Architecture

```
ocmaster.exe (~13MB, Go + giu + CGO)
┌─────────────────────────────────────────┐
│  giu (Dear ImGui Go binding)            │
│  ┌───────────────────────────────────┐  │
│  │  TabBar: 扫描 | 上传 | 设置 | 关于  │  │
│  │  DefaultTheme() dark mode         │  │
│  │  scanner.ScanAll() 直接调用       │  │
│  │  goroutine + shared state         │  │
│  └───────────────────────────────────┘  │
│  OpenGL 3.2+ | GLFW backend            │
├─────────────────────────────────────────┤
│  双模式:                                 │
│    双击 → GUI                            │
│    ocmaster scan → CLI                   │
└─────────────────────────────────────────┘

S端 (不变)
┌─────────┐    ┌──────────────┐    ┌─────────┐
│  Nginx  │───►│ Go Backend   │───►│  MySQL  │
│  :80    │    │ chi/GORM     │    │  :3306  │
└─────────┘    └──────────────┘    └─────────┘
```

## File Tree

```
c-end/hardware-scanner/
├── cmd/
│   ├── cli/main.go          # CLI 入口 (CGO_ENABLED=0, ~2MB 静态)
│   └── gui/                 # GUI 入口 (CGO, ~13MB)
│       ├── main.go           # 双模式 + 主窗口 + loop
│       ├── pages.go           # 4 Tab 页面布局
│       └── theme.go           # 暗色主题
├── scanner/                  # 三平台扫描实现
│   ├── types.go              # HardwareInfo struct
│   ├── scanner_darwin.go     # macOS: sysctl + system_profiler
│   ├── scanner_linux.go      # Linux: /proc + dmidecode + lspci
│   ├── scanner_windows.go    # Windows: WMI (wmic)
│   └── scanner_test.go       # 80.7% coverage
├── go.mod / go.sum
├── bin/                      # 构建产物 (gitignore)
├── build.sh / build.ps1
├── .github/workflows/
│   └── build.yml             # CI: Go CLI + GUI + test → Release
└── VERSION
```

## Data Flow

```
GUI: User Click → goroutine → scanner.ScanAll()
       → shared state → giu loop picks up → Table display
       → exportTXT → os.WriteFile

CLI: ocmaster scan --out json
       → scanner.ScanAll() → json.Encode(os.Stdout)
```

## CI Pipeline

```
build (windows-latest)
  ├── Go CLI 静态编译 + JSON 集成测试
  ├── Go GUI CGO 编译 (MinGW)
  └── ocmaster_0.0.1_windows_amd64.exe → artifact

release (non-PR, needs build + s-end-test)
  └── GitHub Release (tag: v0.0.1)
```

## Design Decisions

| 决策 | 选择 | 原因 |
|------|------|------|
| UI 框架 | giu (Go Dear ImGui) | 跟 ImHex 相同架构, 单二进制 ~13MB |
| 状态管理 | goroutine + shared state | giu loop 每帧自动刷新, 无需 IPC/EventBus |
| 构建 | CGO + MinGW (Win) | giu 需要 OpenGL, 静态链接 |
| CLI | 独立 CGO_ENABLED=0 | ~2MB 纯静态, 服务器/headless 可用 |
| 打包 | 单 exe 直接分发 | 无需 electron-builder/npm |
