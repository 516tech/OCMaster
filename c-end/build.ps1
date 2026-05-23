# OCMaster C端 Windows 构建脚本
# 产出: bin\ Go CLI (静态) + Go GUI (CGO)
param(
    [switch]$Cross = $false
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
    $Version = (Get-Content "$ScriptDir\..\VERSION" -ErrorAction SilentlyContinue | ForEach-Object { $_.Trim() }) -join ""
    if (-not $Version) { $Version = "dev" }
    $LdFlags = "-s -w -X main.version=$Version"
    Write-Host "    Version: $Version" -ForegroundColor Gray

    # [1] Go CLI (静态)
    Write-Host "`n==> [1/2] Go CLI (静态编译)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "0"
    go build -ldflags="$LdFlags" -o "$OutDir\ocmaster.exe" .\cmd\cli
    Write-Host "    -> $OutDir\ocmaster.exe ($((Get-Item "$OutDir\ocmaster.exe").Length / 1KB) KB)" -ForegroundColor Green

    # CLI 集成测试
    Write-Host "`n==> 集成测试..." -ForegroundColor Yellow
    $result = & "$OutDir\ocmaster.exe" scan --out json 2>&1
    if ($LASTEXITCODE -eq 0) {
        $json = $result | ConvertFrom-Json
        Write-Host "    CPU: $($json.cpu.model) ($($json.cpu.cores)C/$($json.cpu.threads)T) PASS" -ForegroundColor Green
    } else {
        Write-Host "    FAILED: $result" -ForegroundColor Red
        exit 1
    }

    # [2] Go GUI (CGO + OpenGL)
    Write-Host "`n==> [2/2] Go GUI (CGO)..." -ForegroundColor Yellow
    $env:CGO_ENABLED = "1"
    go build -ldflags="$LdFlags" -o "$OutDir\ocmaster-gui.exe" .\cmd\gui
    if ($LASTEXITCODE -eq 0) {
        $name = "ocmaster_${Version}_windows_amd64.exe"
        Copy-Item "$OutDir\ocmaster-gui.exe" "$OutDir\$name"
        Write-Host "    -> $OutDir\$name ($([math]::Round((Get-Item "$OutDir\$name").Length / 1MB, 1)) MB)" -ForegroundColor Green
    } else {
        Write-Host "    GUI build skipped (需要 gcc/MinGW + OpenGL)" -ForegroundColor Yellow
    }

    Write-Host "`n=========================================" -ForegroundColor Cyan
    Write-Host " 构建完成" -ForegroundColor Cyan
    Get-ChildItem $OutDir -Filter "*.exe" | ForEach-Object { Write-Host "   $($_.Name) ($([math]::Round($_.Length/1KB,1)) KB)" }
    Write-Host "=========================================" -ForegroundColor Cyan
} finally {
    Pop-Location
}
