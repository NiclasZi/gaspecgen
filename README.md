# Ga-spec-generator

## Overview

`gaspecgen` is a CLI tool / basic html gui providing powerful capabilities for working with Microsoft SQL databases, including the ability to upload and execute arbitrary SQL query templates and fill in data from excel or csv files using go templating inside of SQL files.

> [!WARNING] **SECURITY RISK, SQL INJECTION CAPABILITIES BUILT-IN**
> This tool *deliberately* allows execution of arbitrary SQL queries from user-supplied files.
>
> **This means it must NEVER be exposed to any network-accessible environment** or shared with untrusted users.
>
> **ONLY run `gaspecgen` on your local machine where you are the sole user**.
>
> Exposing this tool over a network or to unauthorized users can have catastrophic consequences.

---

## Recommendations to modify to run as a server

If you want to run `gaspecgen` as a server exposing REST endpoints, you must implement **strong authentication and authorization controls**:

* Use authentication middleware that **only allows logged-in users** access to any endpoints.
* Implement **role-based access control (RBAC)**:

  * *Normal users* should only be able to execute **pre-approved queries** selected from a secure query repository or database.
  * *Privileged users* with higher roles can upload new query templates and manage the query repository.
* Never allow arbitrary query uploads or execution without authentication and authorization checks.

---

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
curl -fsSL https://raw.githubusercontent.com/NiclasZi/gaspecgen/main/scripts/install.sh | bash

```

Check out what the script does [here](https://github.com/NiclasZi/gaspecgen/blob/main/scripts/install.sh).

#### Windows

There is a PowerShell installation script that can be run to install the CLI.

```powershell
powershell -c "irm https://raw.githubusercontent.com/NiclasZi/gaspecgen/main/scripts/install.ps1 | iex"

```

Check out what the script does [here](https://github.com/NiclasZi/gaspecgen/blob/main/scripts/install.ps1).


## Development

### Dev Containers

To simplify the development setup, we provide a **Dev Container** configuration that sets up a ready-to-use development environment including a Microsoft SQL Server instance preloaded with mock data. This allows you to start developing and testing immediately without manual setup.

#### What is a Dev Container?

A **Dev Container** is a Docker-based development environment that runs in an isolated container but integrates seamlessly with your editor (e.g., VS Code). It ensures consistent dependencies and tools across machines and avoids cluttering your host system.

#### How to Use the Dev Container

**Prerequisites**:
Make sure you have these installed:

* [Docker or Docker Desktop](https://www.docker.com/)
* [Visual Studio Code (VS Code)](https://code.visualstudio.com/)
* [Remote Development Extension Pack for VS Code](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)

**Steps to Start the Dev Container**:

1. Clone this repository:

   ```bash
   git clone https://github.com/NiclasZi/gaspecgen.git
   cd gaspecgen
   ```

2. Open the repo folder in **VS Code**.

3. Press `Ctrl + Shift + P` (Windows/Linux) or `Cmd + Shift + P` (macOS) to open the Command Palette.

4. Select **"Dev Containers: Reopen in Container"**.

The devcontainer will start, including:

* The Microsoft SQL Server running on a container
* A sample database preloaded with mock data
* All necessary CLI tools and dependencies installed

You can then begin development with the environment fully configured.

For more info on devcontainers, see the [VS Code documentation](https://code.visualstudio.com/docs/devcontainers/containers).

