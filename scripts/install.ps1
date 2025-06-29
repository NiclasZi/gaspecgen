# -------- Logging --------
function Log-Info($msg)  { Write-Host "[INFO] $msg" -ForegroundColor Green }
function Log-Warn($msg)  { Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Log-Err($msg)   { Write-Host "[ERROR] $msg" -ForegroundColor Red }

# -------- Spinner --------
function Show-Spinner($scriptBlock) {
  $spinner = @('|', '/', '-', '\')
  $i = 0
  $job = Start-Job $scriptBlock

  while ($job.State -eq 'Running') {
    Write-Host -NoNewline "`r[$($spinner[$i % $spinner.Length])] Working..."
    Start-Sleep -Milliseconds 150
    $i++
  }
  Receive-Job $job
  Remove-Job $job
  Write-Host "`r     `r" -NoNewline
}

# -------- Dependency check --------
function Check-Dependencies {
  param ($commands)
  foreach ($cmd in $commands) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
      Log-Err "Missing dependency: $cmd"
      exit 1
    }
  }
}

# -------- Argument parsing --------
$CustomOutputDir = "$env:ProgramFiles\gaspecgen\bin"
$CustomArch = $null
$CustomOS = $null

for ($i = 0; $i -lt $args.Length; $i++) {
  switch ($args[$i]) {
    '-o' | '--output' | '--install-dir' {
      $CustomOutputDir = $args[$i + 1]; $i++
    }
    '-a' | '--arch' | '--architecture' {
      $CustomArch = $args[$i + 1]; $i++
    }
    '--os' | '--operating-system' {
      $CustomOS = $args[$i + 1]; $i++
    }
    default {
      Log-Err "Unknown argument: $($args[$i])"
      exit 1
    }
  }
}

# -------- Binary installer --------
function Install-Binary {
  param ($BinaryName)

  $OS = if ($CustomOS) { $CustomOS } else { "windows" }
  $ARCH = if ($CustomArch) {
    switch ($CustomArch.ToLower()) {
      'x86_64' | 'amd64' { 'amd64' }
      'aarch64' | 'arm64' { 'arm64' }
      default {
        Log-Err "Unsupported architecture: $CustomArch"; exit 1
      }
    }
  } else {
    $env:PROCESSOR_ARCHITECTURE -match 'ARM64' ? 'arm64' : 'amd64'
  }

  $GITHUB_REPO = "NiclasZi/gaspecgen"

  Log-Info "Fetching latest release version..."
  $VERSION = (Invoke-RestMethod "https://api.github.com/repos/$GITHUB_REPO/releases/latest").tag_name
  if (-not $VERSION) {
    Log-Err "Failed to fetch latest version from $GITHUB_REPO"
    exit 1
  }

  $ZipName = "${BinaryName}_${OS}_${ARCH}.zip"
  $DownloadUrl = "https://github.com/$GITHUB_REPO/releases/download/$VERSION/$ZipName"
  $TmpZip = "$env:TEMP\$ZipName"
  $ExtractDir = "$env:TEMP\$BinaryName"

  Log-Info "Downloading $BinaryName $VERSION for $OS $ARCH..."

  Show-Spinner {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TmpZip -UseBasicParsing
  }

  if (-not (Test-Path $TmpZip)) {
    Log-Err "Failed to download binary"
    exit 1
  }

  Log-Info "Extracting..."
  Expand-Archive -Path $TmpZip -DestinationPath $ExtractDir -Force

  $BinaryPath = Join-Path $ExtractDir "$BinaryName.exe"
  if (-not (Test-Path $BinaryPath)) {
    Log-Err "Binary not found in archive"
    exit 1
  }

  if (-not (Test-Path $CustomOutputDir)) {
    New-Item -ItemType Directory -Path $CustomOutputDir -Force | Out-Null
  }

  Copy-Item -Path $BinaryPath -Destination "$CustomOutputDir\$BinaryName.exe" -Force

  Log-Info "$BinaryName installed successfully to $CustomOutputDir"
  Remove-Item $TmpZip -Force
  Remove-Item $ExtractDir -Recurse -Force
}

# -------- Main --------
Check-Dependencies @('Invoke-RestMethod', 'Expand-Archive')

$Binaries = @('gaspecgen')
foreach ($bin in $Binaries) {
  Install-Binary -BinaryName $bin
}
