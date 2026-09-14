#Requires -Version 5
<#
.SYNOPSIS
    コマンド先頭バイト欠落問題の調査用に環境情報を収集するスクリプト
.DESCRIPTION
    シェル環境、エンコーディング設定、コードページ、WSL状態などを収集して表示します。
    問題発生時に diagnose_command_failure スキルから呼び出してください。
#>

param(
    [switch]$Verbose
)

$separator = "=" * 60

function Write-Section {
    param([string]$Title)
    Write-Host ""
    Write-Host $separator
    Write-Host "  $Title"
    Write-Host $separator
}

Write-Host "コマンド実行環境 診断レポート"
Write-Host "生成日時: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"

# --- PowerShell バージョン ---
Write-Section "PowerShell バージョン"
Write-Host "PSVersion    : $($PSVersionTable.PSVersion)"
Write-Host "PSEdition    : $($PSVersionTable.PSEdition)"
Write-Host "Host Name    : $($Host.Name)"
Write-Host "Host Version : $($Host.Version)"
Write-Host "実行ファイル  : $(Get-Process -Id $PID | Select-Object -ExpandProperty Path)"

# --- エンコーディング設定 ---
Write-Section "エンコーディング設定"
Write-Host "OutputEncoding    : $([Console]::OutputEncoding.EncodingName) (CodePage: $([Console]::OutputEncoding.CodePage))"
Write-Host "InputEncoding     : $([Console]::InputEncoding.EncodingName) (CodePage: $([Console]::InputEncoding.CodePage))"
Write-Host "Default (PS)      : $([System.Text.Encoding]::Default.EncodingName)"
try {
    $chcpOutput = & cmd /c chcp 2>&1
    Write-Host "chcp 現在値       : $chcpOutput"
} catch {
    Write-Host "chcp 取得失敗: $_"
}

# --- $OutputEncoding (PowerShell 変数) ---
Write-Section "PowerShell \$OutputEncoding"
if ($null -ne $OutputEncoding) {
    Write-Host "名前      : $($OutputEncoding.EncodingName)"
    Write-Host "CodePage  : $($OutputEncoding.CodePage)"
    Write-Host "BOM付き?  : $($OutputEncoding.GetPreamble().Length -gt 0)"
} else {
    Write-Host "(未設定)"
}

# --- 環境変数 ---
Write-Section "関連する環境変数"
@("LANG", "LC_ALL", "LC_CTYPE", "PYTHONUTF8", "PYTHONIOENCODING") | ForEach-Object {
    $val = [System.Environment]::GetEnvironmentVariable($_)
    Write-Host "${_,-20}: $(if ($val) { $val } else { '(未設定)' })"
}

# --- WSL 状態 ---
Write-Section "WSL 状態"
try {
    $wslStatus = & wsl --status 2>&1 | Out-String
    Write-Host $wslStatus
} catch {
    Write-Host "WSL ステータス取得失敗: $_"
}
try {
    $wslList = & wsl --list --verbose 2>&1 | Out-String
    Write-Host "WSL ディストリビューション:"
    Write-Host $wslList
} catch {
    Write-Host "WSL リスト取得失敗: $_"
}

# --- 重要コマンドの場所 ---
Write-Section "重要コマンドの実行ファイルパス"
@("git", "gh", "go", "wsl", "python", "node") | ForEach-Object {
    $cmd = $_
    try {
        $path = (Get-Command $cmd -ErrorAction Stop).Source
        Write-Host "${cmd,-10}: $path"
    } catch {
        Write-Host "${cmd,-10}: (見つかりません)"
    }
}

# --- 一時ファイルの状況 ---
Write-Section "一時ディレクトリ"
Write-Host "TEMP    : $env:TEMP"
Write-Host "TMP     : $env:TMP"
Write-Host "最近の .ps1 一時ファイル:"
try {
    Get-ChildItem -Path $env:TEMP -Filter "*.ps1" -ErrorAction SilentlyContinue |
        Sort-Object LastWriteTime -Descending |
        Select-Object -First 5 |
        ForEach-Object {
            $bytes = [System.IO.File]::ReadAllBytes($_.FullName)
            $bom = if ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) {
                "UTF-8 BOM あり"
            } elseif ($bytes.Length -ge 2 -and $bytes[0] -eq 0xFF -and $bytes[1] -eq 0xFE) {
                "UTF-16 LE BOM あり"
            } else {
                "BOM なし"
            }
            Write-Host "  $($_.Name) [$bom] - $($_.LastWriteTime)"
        }
} catch {
    Write-Host "  (取得失敗: $_)"
}

# --- 追加の詳細情報 ---
if ($Verbose) {
    Write-Section "PATH 環境変数 (各エントリ)"
    $env:PATH -split ';' | ForEach-Object { Write-Host "  $_" }

    Write-Section "PowerShell プロファイル"
    @(
        $PROFILE.AllUsersAllHosts,
        $PROFILE.AllUsersCurrentHost,
        $PROFILE.CurrentUserAllHosts,
        $PROFILE.CurrentUserCurrentHost
    ) | ForEach-Object {
        $exists = if (Test-Path $_) { "存在します" } else { "存在しません" }
        Write-Host "  $_ → $exists"
    }
}

Write-Host ""
Write-Host $separator
Write-Host "診断レポート終了"
Write-Host $separator
