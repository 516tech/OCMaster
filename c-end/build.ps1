# OCMaster C端 Windows 构建脚本
# 产出: bin\ 目录下 ocmaster_0.0.1_windows_amd64.exe
param(
    [switch]$Cross = $false  # -Cross 启用跨平台编译
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScannerDir = Join-Path $ScriptDir "hardware-scanner"
$ElectronDir = Join-Path $ScriptDir "electron"
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

    # ---------- [1] Go CLI ----------
    Write-Host "`n==> [1/3] 构建 Go CLI (ocmaster.exe)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "0"
    go build -ldflags="$LdFlags" -o "$OutDir\ocmaster.exe" .\cmd\cli
    Write-Host "    -> $OutDir\ocmaster.exe ($((Get-Item "$OutDir\ocmaster.exe").Length / 1KB) KB)" -ForegroundColor Green

    # ---------- [2] 集成测试 ----------
    Write-Host "`n==> [2/3] 集成测试..." -ForegroundColor Yellow
    $name = "ocmaster_${Version}_windows_amd64.exe"
    $testExe = Join-Path $OutDir "ocmaster.exe"
    $result = & $testExe scan --out json 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    CLI scan PASS" -ForegroundColor Green
        $json = $result | ConvertFrom-Json
        Write-Host "    CPU: $($json.cpu.model) ($($json.cpu.cores)C/$($json.cpu.threads)T)" -ForegroundColor Gray
    } else {
        Write-Host "    CLI scan FAILED" -ForegroundColor Red
        Write-Host "    $result" -ForegroundColor Red
        exit 1
    }

    # ---------- [3] Electron ----------
    Write-Host "`n==> [3/3] Electron 构建..." -ForegroundColor Yellow
    if (Test-Path "$ElectronDir\node_modules") {
        Push-Location $ElectronDir
        try {
            npm run build
            if ($LASTEXITCODE -eq 0) {
                Write-Host "    Electron build PASS" -ForegroundColor Green
            }
            Write-Host "    -> electron-builder --win portable ..." -ForegroundColor Gray
            npx electron-builder --win portable --publish=never
            if ($LASTEXITCODE -eq 0) {
                $exe = Get-ChildItem "$ElectronDir\dist\*.exe" | Where-Object { $_.Name -notmatch 'Setup|Installer' } | Select-Object -First 1
                if ($exe) {
                    $finalName = "ocmaster_${Version}_windows_amd64.exe"
                    Rename-Item $exe.FullName $finalName
                    Move-Item (Join-Path $exe.DirectoryName $finalName) "$OutDir\$finalName" -Force
                    Write-Host "    -> $OutDir\$finalName ($([math]::Round($exe.Length/1MB,1)) MB)" -ForegroundColor Green
                }
            }
        } finally {
            Pop-Location
        }
    } else {
        Write-Host "    (跳过: Electron 依赖未安装, 先 cd c-end/electron && npm install)" -ForegroundColor Yellow
    }

    # ---------- 跨平台 (可选) ----------
    if ($Cross) {
        Write-Host "`n==> [*] 跨平台编译..." -ForegroundColor Yellow
        foreach ($target in @(
            @{OS="linux"; Arch="amd64"; Ext=""},
            @{OS="darwin"; Arch="arm64"; Ext=""}
        )) {
            Write-Host "    $($target.OS)/$($target.Arch)..." -NoNewline
            $env:GOOS = $target.OS
            $env:GOARCH = $target.Arch
            $env:CGO_ENABLED = "0"
            $name = "ocmaster_${Version}_$($target.OS)_$($target.Arch)$($target.Ext)"
            go build -ldflags="$LdFlags" -o "$OutDir\$name" .\cmd\cli
            if ($LASTEXITCODE -eq 0) { Write-Host " OK" -ForegroundColor Green }
        }
    }

    Write-Host "`n=========================================" -ForegroundColor Cyan
    Write-Host " 构建完成" -ForegroundColor Cyan
    Get-ChildItem $OutDir | ForEach-Object { Write-Host "   $($_.Name) ($([math]::Round($_.Length/1KB,1)) KB)" -ForegroundColor Gray }
    Write-Host "=========================================" -ForegroundColor Cyan
} finally {
    Pop-Location
}
