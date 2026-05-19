文档版本：V3.0
创建日期：2026-05-18
文档状态：设计稿
所属项目：超频大师（OCMaster）— 架构重构

一、技术选型

| 模块 | 技术 |
|------|------|
| C端 UI | C# + WinUI 3 (Windows App SDK)，仅 Windows 10 1809+ |
| C端硬件扫描 | Go → cgo C shared library (.dll) → C# P/Invoke |
| S端后端 | Golang / chi / GORM v2 / zerolog |
| S端前端 | Vue 3 / Element Plus / TypeScript |
| 本地数据库 | SQLite (:memory: 测试 / 文件 开发) |
| 生产数据库 | MySQL 8.0 (Docker Compose) |

二、C端架构：C# WinUI 3 + Go DLL

```
┌──────────────────────────────────────────┐
│ C# WinUI 3 Desktop App (MSIX / exe)      │
│ ┌──────────┐  ┌──────────────────┐       │
│ │ WinUI 3  │  │  Config Store    │       │
│ │ Pages/   │  │  (JSON file)     │       │
│ │ Controls │  └──────────────────┘       │
│ │          │  ┌──────────────────┐       │
│ │          │  │  HTTP Client     │───→ S端 API
│ └──────────┘  └──────────────────┘       │
│      │ P/Invoke                          │
│      ▼                                   │
│ ┌──────────────────────────────────┐     │
│ │ hardware_scanner.dll (Go c-shared)│    │
│ │ - ScanAll() → JSON               │     │
│ │ - ScanCPU() → JSON               │     │
│ │ - ScanRAM() → JSON               │     │
│ │ - ScanGPU() → JSON               │     │
│ │ - ScanMotherboard() → JSON       │     │
│ └──────────────────────────────────┘     │
└──────────────────────────────────────────┘
```

Go DLL 包含三套平台实现，编译时 `//go:build` 选择：
- `scanner_windows.go` — WMI + Win32 API + SMBus + SPD
- `scanner_linux.go` — sysfs + dmidecode
- `scanner_darwin.go` — IOKit + sysctl

编译产物（Go 风格命名，静态编译）：

| 平台 | CLI (静态) | 动态库 (可选) |
|------|-----------|--------------|
| Windows | `ocmaster_0.0.1_windows_amd64.exe` | `hardware_scanner.dll` |
| macOS | `ocmaster_0.0.1_darwin_arm64` | `libhardware_scanner.dylib` |
| Linux | `ocmaster_0.0.1_linux_amd64` | `libhardware_scanner.so` |

构建命令：
- CLI:  `CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$VER" -o ocmaster_${VER}_${OS}_${ARCH} ./cmd/cli`
- DLL:  `CGO_ENABLED=1 go build -buildmode=c-shared -o hardware_scanner.dll .`
- WinUI 3 单文件: `dotnet publish -c Release -r win-x64 --self-contained -p:PublishSingleFile=true`
- 一键: `.\build.ps1` (Windows 全流程)

WinUI 3 单文件部署: Go DLL 作为 EmbeddedResource 嵌入 exe，首次运行时提取到 `%LOCALAPPDATA%\OCMaster\`。
- DLL:  `go build -buildmode=c-shared -o hardware_scanner.dll .`
- 一键: `bash build.sh` 或 `.\build.ps1`

三、目录结构

```
OCMaster/
├── c-end/
│   ├── OCMaster.App/             # C# WinUI 3 项目
│   │   ├── OCMaster.App.csproj
│   │   ├── MainWindow.xaml
│   │   ├── Pages/
│   │   │   ├── ScanPage.xaml
│   │   │   ├── SettingsPage.xaml
│   │   │   └── UploadPage.xaml
│   │   ├── Services/
│   │   │   ├── NativeInterop.cs        # P/Invoke 声明
│   │   │   ├── ApiClient.cs
│   │   │   └── ConfigService.cs
│   │   └── app.manifest
│   ├── hardware-scanner/         # Go DLL + CLI 项目
│   │   ├── go.mod / go.sum
│   │   ├── main.go               # c-shared DLL 入口 (//export ScanAll, ...)
│   │   ├── cmd/cli/main.go       # CLI 独立可执行文件入口
│   │   ├── scanner/
│   │   │   ├── types.go          # HardwareInfo JSON struct
│   │   │   ├── scanner_darwin.go # macOS: sysctl + system_profiler
│   │   │   ├── scanner_windows.go# Windows: WMI + Win32 (stub→implement)
│   │   │   ├── scanner_linux.go  # Linux: /proc + sysfs + dmidecode
│   │   │   └── scanner_test.go   # 测试 (80.7% 覆盖)
│   ├── bin/                      # 构建产物 (gitignore)
│   ├── build.sh                  # macOS/Linux 一键构建
│   └── build.ps1                 # Windows 一键构建 + 集成测试
├── s-end/
│   ├── backend/                  # (不变)
│   └── frontend/                 # (不变)
├── docker-compose.yml
├── prd.md / architecture.md / bug.md / evaluate.md
└── readme.md
```

四、Go DLL C 导出接口

```go
//export ScanAll
func ScanAll() *C.char {
    info := collectHardwareInfo()
    jsonBytes, _ := json.Marshal(info)
    return C.CString(string(jsonBytes))
}

//export FreeString
func FreeString(s *C.char) {
    C.free(unsafe.Pointer(s))
}
```

C# P/Invoke 调用：

```csharp
public static class NativeInterop
{
    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern IntPtr ScanAll();

    [DllImport("hardware_scanner.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern void FreeString(IntPtr ptr);

    public static string ScanAllManaged()
    {
        var ptr = ScanAll();
        var json = Marshal.PtrToStringUTF8(ptr);
        FreeString(ptr);
        return json;
    }
}
```

五、S端（不变）

DDD 四层，Go/chi/GORM/Zerolog，Vue 3 SPA。
详细参考 V2.0 architecture.md。

六、部署与 CI

| 组件 | 部署方式 |
|------|---------|
| C端 Windows | Go CLI (ocmaster.exe) + Go DLL (hardware_scanner.dll) + C# WinUI 3 |
| C端 macOS | Go CLI (ocmaster) + Go dylib (libhardware_scanner.dylib) |
| C端 Linux | Go CLI (ocmaster) + Go so (libhardware_scanner.so) |
| S端 | Docker Compose (MySQL + backend + nginx) |

CI (build-windows.yml):
- `go-build` job: Setup Go → Build CLI + DLL → `ocmaster.exe scan --out json` → 验证字段非空 → upload artifacts
- `s-end-test` job: go test ./... -cover
