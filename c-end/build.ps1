# OCMaster C端 Windows 构建脚本
# 产出: bin\ 目录下 ocmaster.exe + hardware_scanner.dll
param(
    [switch]$Cross = $false  # -Cross 启用跨平台编译
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScannerDir = Join-Path $ScriptDir "hardware-scanner"
$OutDir = Join-Path $ScriptDir "bin"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " OCMaster C端 Windows 构建" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

Push-Location $ScannerDir

try {
    # 读取版本号
    $Version = (Get-Content "$ScriptDir\..\VERSION" -ErrorAction SilentlyContinue | ForEach-Object { $_.Trim() }) -join ""
    if (-not $Version) { $Version = "dev" }
    $LdFlags = "-s -w -X main.version=$Version"
    Write-Host "    Version: $Version" -ForegroundColor Gray

    # ---------- Windows CLI (exe) ----------
    Write-Host "`n==> [1] 构建 Windows CLI (ocmaster.exe)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $name = "ocmaster_${Version}_windows_amd64.exe"
    go build -ldflags="$LdFlags" -o "$OutDir\$name" .\cmd\cli
    Write-Host "    -> $OutDir\$name ($((Get-Item "$OutDir\$name").Length / 1KB) KB)" -ForegroundColor Green

    # ---------- Windows DLL ----------
    Write-Host "`n==> [2] 构建 Windows DLL (hardware_scanner.dll)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "1"
    go build -buildmode=c-shared -o "$OutDir\hardware_scanner.dll" .
    Write-Host "    -> $OutDir\hardware_scanner.dll ($((Get-Item "$OutDir\hardware_scanner.dll").Length / 1KB) KB)" -ForegroundColor Green

    # 复制到 C# 项目目录 (供 embed + publish)
    Copy-Item "$OutDir\hardware_scanner.dll" "$ScriptDir\OCMaster.App\" -Force
    Write-Host "    -> 复制到 OCMaster.App\ (供嵌入资源)" -ForegroundColor Gray

    # ---------- 集成测试 ----------
    Write-Host "`n==> [3] 集成测试..." -ForegroundColor Yellow
    $name = "ocmaster_${Version}_windows_amd64.exe"
    $result = & "$OutDir\$name" scan --out json 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    CLI scan PASS" -ForegroundColor Green
        $json = $result | ConvertFrom-Json
        Write-Host "    CPU: $($json.cpu.model) ($($json.cpu.cores)C/$($json.cpu.threads)T)" -ForegroundColor Gray
        Write-Host "    RAM: $($json.ram.totalCapacity)" -ForegroundColor Gray
        Write-Host "    GPU: $($json.gpu.model)" -ForegroundColor Gray
    } else {
        Write-Host "    CLI scan FAILED" -ForegroundColor Red
        Write-Host "    $result" -ForegroundColor Red
        exit 1
    }

    # ---------- WinUI 3 发布 (单文件 exe) ----------
    Write-Host "`n==> [4] dotnet publish (单文件 exe)..." -ForegroundColor Yellow
    Push-Location "$ScriptDir\OCMaster.App"
    try {
        dotnet publish -c Release -r win-x64 --self-contained `
            -p:PublishSingleFile=true `
            -p:WindowsAppSDKSelfContained=true `
            -o "$OutDir\publish" 2>&1 | Select-Object -Last 3
        if ($LASTEXITCODE -eq 0) {
            $exe = Get-ChildItem "$OutDir\publish\OCMaster.App.exe" -ErrorAction SilentlyContinue
            if ($exe) {
                Write-Host "    -> $($exe.FullName) ($($exe.Length / 1MB) MB)" -ForegroundColor Green
            }
            Write-Host "    WinUI 3 single-file exe published" -ForegroundColor Green
        } else {
            Write-Host "    dotnet publish failed (可能未安装 .NET SDK，跳过)" -ForegroundColor Yellow
        }
    } finally {
        Pop-Location
    }

    # ---------- 跨平台 (可选) ----------
    if ($Cross) {
        Write-Host "`n==> [4] 跨平台编译..." -ForegroundColor Yellow
        foreach ($target in @(
            @{OS="linux"; Arch="amd64"; Ext=""},
            @{OS="darwin"; Arch="arm64"; Ext=""},
            @{OS="darwin"; Arch="amd64"; Ext=""}
        )) {
            Write-Host "    $($target.OS)/$($target.Arch)..." -NoNewline
            $env:GOOS = $target.OS
            $env:GOARCH = $target.Arch
            $env:CGO_ENABLED = "0"
            $name = "ocmaster_${Version}_$($target.OS)_$($target.Arch)$($target.Ext)"
            go build -ldflags="$LdFlags" -o "$OutDir\$name" .\cmd\cli
            if ($LASTEXITCODE -eq 0) {
                Write-Host " OK" -ForegroundColor Green
            }
        }
    }

    Write-Host "`n=========================================" -ForegroundColor Cyan
    Write-Host " 构建完成" -ForegroundColor Cyan
    Get-ChildItem $OutDir | ForEach-Object { Write-Host "   $($_.Name) ($($_.Length / 1KB) KB)" -ForegroundColor Gray }
    Write-Host "=========================================" -ForegroundColor Cyan
} finally {
    Pop-Location
}
