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
    # ---------- Windows CLI (exe) ----------
    Write-Host "`n==> [1] 构建 Windows CLI (ocmaster.exe)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -ldflags="-s -w" -o "$OutDir\ocmaster.exe" .\cmd\cli
    Write-Host "    -> $OutDir\ocmaster.exe ($((Get-Item "$OutDir\ocmaster.exe").Length / 1KB) KB)" -ForegroundColor Green

    # ---------- Windows DLL ----------
    Write-Host "`n==> [2] 构建 Windows DLL (hardware_scanner.dll)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "1"
    go build -buildmode=c-shared -o "$OutDir\hardware_scanner.dll" .
    Write-Host "    -> $OutDir\hardware_scanner.dll ($((Get-Item "$OutDir\hardware_scanner.dll").Length / 1KB) KB)" -ForegroundColor Green

    # ---------- 集成测试 ----------
    Write-Host "`n==> [3] 集成测试..." -ForegroundColor Yellow
    $result = & "$OutDir\ocmaster.exe" scan --out json 2>&1
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
            $name = "ocmaster-$($target.OS)-$($target.Arch)$($target.Ext)"
            go build -ldflags="-s -w" -o "$OutDir\$name" .\cmd\cli
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
