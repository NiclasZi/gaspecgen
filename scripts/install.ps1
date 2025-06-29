# -------- Logging --------
function Log-Info($msg)  { Write-Host "[INFO]  $msg" -ForegroundColor Green }
function Log-Warn($msg)  { Write-Host "[WARN]  $msg" -ForegroundColor Yellow }
function Log-Err($msg)   { Write-Host "[ERROR] $msg" -ForegroundColor Red }

# -------- Spinner --------
function Show-Spinner {
  param (
    [ScriptBlock]$ScriptBlock
  )

  $spinner = @('|', '/', '-', '\')
  $i = 0
  $job = Start-Job $ScriptBlock

  while ($job.State -eq 'Running') {
    Write-Host -NoNewline "`r[$($spinner[$i % $spinner.Length])] Working..."
    Start-Sleep -Milliseconds 150
    $i++
  }

  Write-Host "`r     `r" -NoNewline
  Receive-Job $job
  Remove-Job $job
}

# -------- Dependency check --------
function Check-Dependencies {
  param ([string[]]$commands)
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
    '-o' { $CustomOutputDir = $args[$i + 1]; $i++ }
    '--output' { $CustomOutputDir = $args[$i + 1]; $i++ }
    '--install-dir' { $CustomOutputDir = $args[$i + 1]; $i++ }

    '-a' { $CustomArch = $args[$i + 1]; $i++ }
    '--arch' { $CustomArch = $args[$i + 1]; $i++ }
    '--architecture' { $CustomArch = $args[$i + 1]; $i++ }

    '--os' { $CustomOS = $args[$i + 1]; $i++ }
    '--operating-system' { $CustomOS = $args[$i + 1]; $i++ }

    default {
      Log-Err "Unknown argument: $($args[$i])"
      exit 1
    }
  }
}

# -------- Binary installer --------
function Install-Binary {
  param (
    [string]$BinaryName
  )

  $OS = if ($CustomOS) { $CustomOS } else { "windows" }
  $ARCH = if ($CustomArch) {
    switch ($CustomArch.ToLower()) {
      'x86_64' | 'amd64' { 'amd64' }
      'aarch64' | 'arm64' { 'arm64' }
      default {
        Log-Err "Unsupported architecture: $CustomArch"
        exit 1
      }
    }
  } else {
    if ($env:PROCESSOR_ARCHITECTURE -match 'ARM64') { 'arm64' } else { 'amd64' }
  }

  $GITHUB_REPO = "NiclasZi/gaspecgen"

  Log-Info "Fetching latest release version..."
  try {
    $VERSION = (Invoke-RestMethod "https://api.github.com/repos/$GITHUB_REPO/releases/latest").tag_name
  } catch {
    Log-Err "Failed to fetch release: $_"
    exit 1
  }

  if (-not $VERSION) {
    Log-Err "Could not determine release version"
    exit 1
  }

  $ZipName = "${BinaryName}_${OS}_${ARCH}.zip"
  $DownloadUrl = "https://github.com/$GITHUB_REPO/releases/download/$VERSION/$ZipName"
  $TmpZip = "$env:TEMP\$ZipName"
  $ExtractDir = "$env:TEMP\$BinaryName"

  Log-Info "Downloading $BinaryName $VERSION for $OS/$ARCH..."
  Show-Spinner {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TmpZip -UseBasicParsing
  }

  if (-not (Test-Path $TmpZip)) {
    Log-Err "Failed to download binary zip file."
    exit 1
  }

  Log-Info "Extracting archive..."
  Expand-Archive -Path $TmpZip -DestinationPath $ExtractDir -Force

  $BinaryPath = Join-Path $ExtractDir "$BinaryName.exe"
  if (-not (Test-Path $BinaryPath)) {
    Log-Err "Binary not found in archive"
    exit 1
  }

  if (-not (Test-Path $CustomOutputDir)) {
    New-Item -ItemType Directory -Path $CustomOutputDir -Force | Out-Null
  }

  Copy-Item -Path $BinaryPath -Destination (Join-Path $CustomOutputDir "$BinaryName.exe") -Force
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
