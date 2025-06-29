# Ga-spec-generator - Microsoft SQL CLI tool

## Overview

A tool CLI with tools for working with a Microsoft SQL db.

## Installation

### Using go install

This tool can be installed using `go`.

> ![NOTE]
> This requires go to be installed with correct PATH setup for go installed programs.

```bash
go install github.com/NiclasZi/gaspecgen
```

### Download and install binary

#### Mac and Linux

For Mac and Linux, there is an installation script that can be run to install the CLI.

##### Prerequisites

- bash
- curl

```bash
curl -fsSL https://raw.githubusercontent.com/Phillezi/gaspecgen/main/scripts/install.sh | bash

```

Check out what the script does [here](https://github.com/NiclasZi/gaspecgen/blob/main/scripts/install.sh).

#### Windows

There is a PowerShell installation script that can be run to install the CLI.

```powershell
powershell -c "irm https://raw.githubusercontent.com/Phillezi/gaspecgen/main/scripts/install.ps1 | iex"

```

Check out what the script does [here](https://github.com/NiclasZi/gaspecgen/blob/main/scripts/install.ps1).
