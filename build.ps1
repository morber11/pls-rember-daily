param(
    [switch]$AllowMultipleRuns,
    [switch]$SkipStartupDelay
)

$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$outputDirectory = Join-Path $projectRoot "dist"
$outputPath = Join-Path $outputDirectory "PlsRemberDaily.exe"
$iconPath = Join-Path $projectRoot "assets\images\plsrember.ico"
$resourcePath = Join-Path $projectRoot "cmd\pls-rember-daily\rsrc.syso"
$resourceDirectory = Split-Path -Parent $resourcePath

Set-Location $projectRoot
New-Item -ItemType Directory -Path $outputDirectory -Force | Out-Null

if (-not (Test-Path -LiteralPath $iconPath -PathType Leaf)) {
    throw "Icon file not found: $iconPath"
}

Get-ChildItem -LiteralPath $resourceDirectory -Filter "rsrc*.syso" |
    Remove-Item -Force

$rsrc = Get-Command rsrc -ErrorAction SilentlyContinue
if ($null -ne $rsrc) {
    & $rsrc.Source -ico $iconPath -o $resourcePath
}
else {
    go run github.com/akavel/rsrc@v0.10.2 -ico $iconPath -o $resourcePath
}
if ($LASTEXITCODE -ne 0) {
    throw "rsrc failed"
}

go test ./...
if ($LASTEXITCODE -ne 0) {
    throw "go test failed"
}

go vet ./...
if ($LASTEXITCODE -ne 0) {
    throw "go vet failed"
}

$buildTags = @()
if ($AllowMultipleRuns) {
    $buildTags += "allowmultipleruns"
    Write-Host "building with AllowMultipleRuns" -ForegroundColor Yellow
}

if ($SkipStartupDelay) {
    $buildTags += "skipstartupdelay"
    Write-Host "building with SkipStartupDelay" -ForegroundColor Yellow
}

$buildArguments = @("-ldflags=-H windowsgui", "-o", $outputPath)
if ($buildTags.Count -gt 0) {
    $buildArguments += @("-tags", ($buildTags -join ","))
}
$buildArguments += ".\cmd\pls-rember-daily"

go build @buildArguments
if ($LASTEXITCODE -ne 0) {
    throw "go build failed"
}

Write-Host "Built $outputPath"
