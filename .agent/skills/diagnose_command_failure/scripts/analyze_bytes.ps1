#Requires -Version 5
<#
.SYNOPSIS
    コマンド文字列のバイト列を解析し、BOMや文字化けの有無を検査するスクリプト
.PARAMETER InputString
    解析対象の文字列（失敗したコマンドなど）
.PARAMETER FilePath
    解析対象のファイルパス（ファイルから読み込む場合）
#>

param(
    [string]$InputString = "",
    [string]$FilePath = ""
)

$separator = "-" * 60

function Show-ByteAnalysis {
    param(
        [byte[]]$Bytes,
        [string]$Label
    )

    Write-Host ""
    Write-Host "=== $Label ==="
    Write-Host "バイト数: $($Bytes.Length)"

    # 先頭16バイトのHEX表示
    $headCount = [Math]::Min(16, $Bytes.Length)
    $hexStr = ($Bytes[0..($headCount - 1)] | ForEach-Object { "0x{0:X2}" -f $_ }) -join " "
    Write-Host "先頭${headCount}バイト (HEX): $hexStr"

    # BOM チェック
    if ($Bytes.Length -ge 3 -and $Bytes[0] -eq 0xEF -and $Bytes[1] -eq 0xBB -and $Bytes[2] -eq 0xBF) {
        Write-Host "⚠️  BOM検出: UTF-8 BOM (EF BB BF) が先頭にあります！"
        Write-Host "   → これが原因でコマンドの最初の文字が欠落する可能性があります"
        $contentBytes = $Bytes[3..($Bytes.Length - 1)]
        $content = [System.Text.Encoding]::UTF8.GetString($contentBytes)
        Write-Host "   BOM除去後の内容: $($content.Substring(0, [Math]::Min(80, $content.Length)))"
    }
    elseif ($Bytes.Length -ge 2 -and $Bytes[0] -eq 0xFF -and $Bytes[1] -eq 0xFE) {
        Write-Host "⚠️  BOM検出: UTF-16 LE BOM (FF FE) が先頭にあります！"
    }
    elseif ($Bytes.Length -ge 2 -and $Bytes[0] -eq 0xFE -and $Bytes[1] -eq 0xFF) {
        Write-Host "⚠️  BOM検出: UTF-16 BE BOM (FE FF) が先頭にあります！"
    }
    else {
        Write-Host "✅ BOMなし"
    }

    # 最初の数文字の ASCII 表示
    Write-Host ""
    Write-Host "先頭16バイトの文字表示:"
    for ($i = 0; $i -lt $headCount; $i++) {
        $b = $Bytes[$i]
        $char = if ($b -ge 0x20 -and $b -lt 0x7F) { [char]$b } else { "." }
        Write-Host "  Byte[$i]: 0x{0:X2} = '{1}' (dec: {2})" -f $b, $char, $b
    }

    # マルチバイト文字の途中切断チェック
    Write-Host ""
    Write-Host "エンコーディング別デコード試行:"
    @(
        @{ Name = "UTF-8"; Enc = [System.Text.Encoding]::UTF8 },
        @{ Name = "Shift-JIS"; Enc = [System.Text.Encoding]::GetEncoding(932) },
        @{ Name = "UTF-16 LE"; Enc = [System.Text.Encoding]::Unicode }
    ) | ForEach-Object {
        try {
            $decoded = $_.Enc.GetString($Bytes)
            $preview = $decoded.Substring(0, [Math]::Min(60, $decoded.Length)) -replace "`r|`n", " "
            Write-Host "  $($_.Name,-12): $preview"
        }
        catch {
            Write-Host "  $($_.Name,-12): (デコード失敗: $_)"
        }
    }
}

# --- メイン処理 ---
Write-Host "コマンド文字列 バイト解析ツール"
Write-Host "実行日時: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
Write-Host $separator

if ($FilePath -ne "") {
    if (Test-Path $FilePath) {
        Write-Host "対象ファイル: $FilePath"
        $bytes = [System.IO.File]::ReadAllBytes($FilePath)
        Show-ByteAnalysis -Bytes $bytes -Label "ファイルバイト解析"
    }
    else {
        Write-Host "エラー: ファイルが見つかりません: $FilePath"
        exit 1
    }
}

if ($InputString -ne "") {
    Write-Host ""
    Write-Host "対象文字列: $InputString"
    # UTF-8 エンコーディングでバイト列に変換
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($InputString)
    Show-ByteAnalysis -Bytes $bytes -Label "文字列バイト解析 (UTF-8 エンコード)"

    # PowerShell の デフォルトエンコーディングでも確認
    $defaultBytes = [System.Text.Encoding]::Default.GetBytes($InputString)
    if ($defaultBytes[0] -ne $bytes[0]) {
        Show-ByteAnalysis -Bytes $defaultBytes -Label "文字列バイト解析 (System Default エンコード)"
    }
}

if ($FilePath -eq "" -and $InputString -eq "") {
    Write-Host ""
    Write-Host "使い方:"
    Write-Host "  .\analyze_bytes.ps1 -InputString 'git status'"
    Write-Host "  .\analyze_bytes.ps1 -FilePath 'C:\tmp\script.ps1'"
}

Write-Host ""
Write-Host $separator
Write-Host "解析終了"
