# -------- Logging --------
function Log-Info($msg)  { Write-Host "[INFO]  $msg" -ForegroundColor Green }
function Log-Warn($msg)  { Write-Host "[WARN]  $msg" -ForegroundColor Yellow }
function Log-Err($msg)   { Write-Host "[ERROR] $msg" -ForegroundColor Red }

# -------- Spinner --------
function Show-Spinner {
    param (
        [ScriptBlock]$Until
    )

    $spinner = @('|', '/', '-', '\')
    $i = 0

    while (-not (& $Until)) {
        Write-Host -NoNewline "`r[$($spinner[$i % $spinner.Length])] Working..."
        Start-Sleep -Milliseconds 150
        $i++
    }

    Write-Host "`r     `r" -NoNewline
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
$CustomOutputDir = "$env:USERPROFILE\bin"
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
  param ([string]$BinaryName)

  # OS Detection
  $OS = if ($CustomOS) { $CustomOS.ToLower() } else { "windows" }

  # Arch Detection
  if ($CustomArch) {
    switch ($CustomArch.ToLower()) {
      'x86_64' { $ARCH = 'amd64' }
      'amd64'  { $ARCH = 'amd64' }
      'aarch64' { $ARCH = 'arm64' }
      'arm64'  { $ARCH = 'arm64' }
      default {
        Log-Err "Unsupported architecture: $CustomArch"
        exit 1
      }
    }
  } else {
    if ($env:PROCESSOR_ARCHITECTURE -match 'ARM64') {
      $ARCH = 'arm64'
    } else {
      $ARCH = 'amd64'
    }
  }

  Log-Info "Using architecture: $ARCH"
  Log-Info "Using OS: $OS"

  $GITHUB_REPO = "NiclasZi/gaspecgen"

  Log-Info "Fetching latest release version..."
  try {
    $response = Invoke-RestMethod "https://api.github.com/repos/$GITHUB_REPO/releases/latest"

    if (-not $response -or -not $response.tag_name) {
      Log-Err "Release response is invalid or missing tag_name"
      Log-Warn "Response content: $($response | Out-String)"
      exit 1
    }

    $VERSION = $response.tag_name
    Log-Info "Using version: $VERSION"
  } catch {
    Log-Err "Failed to fetch release: $_"
    exit 1
  }
  
  $ZipName = "${BinaryName}_${OS}_${ARCH}.zip"
  $DownloadUrl = "https://github.com/$GITHUB_REPO/releases/download/$VERSION/$ZipName"
  $TmpZip = "$env:TEMP\$ZipName"
  $ExtractDir = "$env:TEMP\$BinaryName"

  if (-not $DownloadUrl) {
    Log-Err "Download URL is empty. Cannot continue."
    exit 1
  }

  Log-Info "Download URL: $DownloadUrl"
  Log-Info "Downloading $BinaryName $VERSION for $OS/$ARCH..."

  $webClient = New-Object System.Net.WebClient
  $downloadTask = $webClient.DownloadFileTaskAsync($DownloadUrl, $TmpZip)

  Show-Spinner { $downloadTask.IsCompleted }

  # Wait for the download task to complete and handle errors
  try {
      $downloadTask.Wait()
  } catch {
      Log-Err "Failed to download the binary... :("
      Log-Warn "Check if the URL is correct:"
      Write-Host $DownloadUrl
      exit 1
  }

  if (-not (Test-Path $TmpZip)) {
    Log-Err "Failed to download binary zip file."
    exit 1
  }

  Log-Info "Extracting archive..."
  Expand-Archive -Path $TmpZip -DestinationPath $ExtractDir -Force

  $BinaryFile = Get-ChildItem -Path $ExtractDir -Filter "${BinaryName}_${OS}_${ARCH}.exe" -Recurse | Select-Object -First 1

  if (-not $BinaryFile) {
      Log-Err "Binary file '$BinaryName' not found in extracted folder."
      exit 1
  } else {
      Log-Info "Found binary at: $($BinaryFile.FullName)"
  }

  if (-not (Test-Path $CustomOutputDir)) {
    New-Item -ItemType Directory -Path $CustomOutputDir -Force | Out-Null
  }

  Copy-Item -Path "$($BinaryFile.FullName)" -Destination (Join-Path $CustomOutputDir "$BinaryName.exe") -Force
  Log-Info "$BinaryName installed successfully to $CustomOutputDir"

  Remove-Item $TmpZip -Force
  Remove-Item $ExtractDir -Recurse -Force

  $currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
  if (-not $currentUserPath) { $currentUserPath = "" }

  if (-not ($currentUserPath.Split(';') -contains $CustomOutputDir)) {
      $newUserPath = if ($currentUserPath -eq "") { $CustomOutputDir } else { "$currentUserPath;$CustomOutputDir" }
      [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
      Log-Info "Added $CustomOutputDir to user PATH. Restart your terminal to apply changes."
  }

}

# -------- Main --------
Check-Dependencies @('Invoke-RestMethod', 'Expand-Archive', 'Invoke-WebRequest')

$Binaries = @('gaspecgen')
foreach ($bin in $Binaries) {
  Install-Binary -BinaryName $bin
}
