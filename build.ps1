$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$guiProject = Join-Path $root "src\WC4SaveEditor.Gui\WC4SaveEditor.Gui.csproj"
$publishDir = Join-Path $root "src\WC4SaveEditor.Gui\bin\Release\net10.0-windows\win-x64\publish"

Write-Host "Publishing self-contained single-file build..."
dotnet publish $guiProject -c Release -r win-x64 --self-contained true `
    -p:PublishSingleFile=true `
    -p:IncludeNativeLibrariesForSelfExtract=true `
    -p:EnableCompressionInSingleFile=true

if ($LASTEXITCODE -ne 0) {
    Write-Error "Publish failed."
    exit 1
}

$exePath = Join-Path $publishDir "WC4SaveEditor.Gui.exe"
$destination = Join-Path $root "WC4SaveEditor.exe"
Copy-Item -Path $exePath -Destination $destination -Force

Write-Host "Done: $destination"
