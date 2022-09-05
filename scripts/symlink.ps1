Clear-Host
Set-Location -Path $PSScriptRoot

try {
    $files = Get-ChildItem *.axs, *.jar, *.tko, *.xdd -File -Recurse -ErrorAction Stop
    if (!$files) { Write-Host "No files found"; exit }

    if (Test-Path -Path .\.env) {
        $env = Get-Content -Path .\.env -ErrorAction Stop | Out-String | ConvertFrom-StringData
        $sourceDirectory = $env.SOURCE_DIRECTORY_RELATIVE_PATH
    }
    else {
        $sourceDirectory = ".."
    }

    Write-Host "Symlinking $($files.Count) files to $sourceDirectory"
    Write-Host

    foreach ($file in $files) {
        $fileName = $file.Name
        $fileDirectory = Join-Path ".." $(($file.DirectoryName).Split("\")[-1])

        $linkPath = Join-Path $sourceDirectory $fileName
        if (Test-Path -Path $linkPath) {
            Write-Host "Removing existing symlink: $linkPath"
            Remove-Item -Path $linkPath -Force
        }

        $target = Join-Path $fileDirectory $fileName

        Write-Host "Creating symlink: $linkPath --> $target"
        New-Item -ItemType SymbolicLink -Path $linkPath -Target $target -ErrorAction Stop | Out-Null
    }
}
catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host
Write-Host "Complete" -ForegroundColor Green
Read-Host -Prompt "Press any key to exit..."
