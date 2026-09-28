#Requires -Version 5
<#
.SYNOPSIS
    エージェントが生成した一時ファイルのBOM状態を確認するスクリプト
.DESCRIPTION
    TEMP ディレクトリの最近の .ps1 / .sh / .bat ファイルを検索し、
    BOMの有無・エンコーディングを確認します。
#>

param(
    [int]$RecentMinutes = 60,
    [string]$TempDir = $env:TEMP
)

$separator = "=" * 60
$cutoff = (Get-Date).AddMinutes(-$RecentMinutes)

Write-Host "一時ファイル BOM 確認ツール"
Write-Host "実行日時: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
Write-Host "対象: $TempDir 配下の直近 ${RecentMinutes}分以内のスクリプトファイル"
Write-Host $separator

function Get-BomInfo {
    param([string]$FilePath)
    try {
        $bytes = [System.IO.File]::ReadAllBytes($FilePath)
        if ($bytes.Length -eq 0) {
            return @{ HasBom = $false; Type = "(空ファイル)"; Bytes = @() }
        }
        if ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) {
            return @{ HasBom = $true; Type = "UTF-8 BOM"; Bytes = $bytes }
        }
        if ($bytes.Length -ge 2 -and $bytes[0] -eq 0xFF -and $bytes[1] -eq 0xFE) {
            return @{ HasBom = $true; Type = "UTF-16 LE BOM"; Bytes = $bytes }
        }
        if ($bytes.Length -ge 2 -and $bytes[0] -eq 0xFE -and $bytes[1] -eq 0xFF) {
            return @{ HasBom = $true; Type = "UTF-16 BE BOM"; Bytes = $bytes }
        }
        return @{ HasBom = $false; Type = "BOM なし (UTF-8 or other)"; Bytes = $bytes }
    }
    catch {
        return @{ HasBom = $false; Type = "読み取り失敗: $_"; Bytes = @() }
    }
}

# スクリプトファイルを収集
$extensions = @("*.ps1", "*.sh", "*.bat", "*.cmd")
$files = @()
foreach ($ext in $extensions) {
    $found = Get-ChildItem -Path $TempDir -Filter $ext -Recurse -ErrorAction SilentlyContinue |
    Where-Object { $_.LastWriteTime -ge $cutoff }
    $files += $found
}

if ($files.Count -eq 0) {
    Write-Host ""
    Write-Host "対象ファイルが見つかりませんでした。"
    Write-Host "(直近 ${RecentMinutes}分以内の .ps1/.sh/.bat/.cmd ファイル)"
}
else {
    Write-Host ""
    Write-Host "発見したファイル: $($files.Count) 件"
    Write-Host ""

    $bomFiles = @()
    foreach ($file in ($files | Sort-Object LastWriteTime -Descending)) {
        $info = Get-BomInfo -FilePath $file.FullName
        $icon = if ($info.HasBom) { "⚠️ " } else { "✅ " }
        $sizeKB = [Math]::Round($file.Length / 1024, 1)
        Write-Host "${icon} $($file.Name)"
        Write-Host "   パス     : $($file.FullName)"
        Write-Host "   更新日時 : $($file.LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss'))"
        Write-Host "   サイズ   : ${sizeKB} KB"
        Write-Host "   BOM状態  : $($info.Type)"

        if ($info.Bytes.Length -ge 4) {
            $preview = ($info.Bytes[0..3] | ForEach-Object { "0x{0:X2}" -f $_ }) -join " "
            Write-Host "   先頭4byte: $preview"
        }

        if ($info.HasBom) {
            $bomFiles += $file.FullName
            # ファイルの中身のプレビュー（BOM除去後）
            if ($info.Bytes.Length -ge 3) {
                $skipBytes = if ($info.Type -like "UTF-16*") { 2 } else { 3 }
                $content = [System.Text.Encoding]::UTF8.GetString($info.Bytes[$skipBytes..([Math]::Min(100, $info.Bytes.Length - 1))])
                $preview = $content -replace "`r|`n", " "
                Write-Host "   内容先頭  : $preview"
            }
        }
        Write-Host ""
    }

    # サマリー
    Write-Host $separator
    if ($bomFiles.Count -gt 0) {
        Write-Host "⚠️  BOM付きファイル: $($bomFiles.Count) 件 検出"
        Write-Host "これらがコマンド失敗の原因になっている可能性があります。"
        $bomFiles | ForEach-Object { Write-Host "   - $_" }
    }
    else {
        Write-Host "✅ BOM付きファイルは検出されませんでした。"
    }
}

# おまけ: $env:TEMP 以外の一般的な一時場所も確認
Write-Host ""
Write-Host $separator
Write-Host "その他の一時ファイルの場所:"
@("C:\tmp", "C:\Temp", "$env:USERPROFILE\AppData\Local\Temp") | ForEach-Object {
    if (Test-Path $_) {
        $count = (Get-ChildItem -Path $_ -Filter "*.ps1" -ErrorAction SilentlyContinue).Count
        Write-Host "  $_ → .ps1 ファイル数: $count"
    }
}

Write-Host ""
Write-Host $separator
Write-Host "確認終了"
