# UVE - UV Environment Manager

UVE is a lightweight Python virtual environment manager that works with [UV](https://github.com/astral-sh/uv). It provides a conda-like experience for managing Python virtual environments from any directory, while leveraging UV's speed and reliability.

## Features

- Central storage of virtual environments (similar to conda)
- Create and manage Python environments from any directory
- Cross-platform support (Linux, macOS, Windows)
- Written in Go - no Python dependency for the manager itself
- Simple shell integration for environment activation/deactivation
- Uses UV for fast environment creation

## Prerequisites

- [Go](https://golang.org/doc/install)  - Optional, only needed if building from source
- [UV](https://github.com/astral-sh/uv) - Installed and available in PATH

## Installation

### Installing from source:

If you are building it yourself, clone the repo and build the binary:

```bash
git clone https://github.com/robert-mcdermott/uve.git
cd uve
go build -o uve-bin main.go
```

### Installation from pre-built binary:

If you've rather just use a pre-built binary, download the binary archive for your platform (os and architecture) from the [releases page](https://github.com/robert-mcdermott/uve/releases)
and extract the archive to a directory of your choice.

### Install the binary:

Linux/macOS:
```bash
# Create bin directory if it doesn't exist
mkdir -p ~/bin
# Copy the binary as uve-bin
cp uve-bin ~/bin/uve-bin
# Add to PATH (add this to your .bashrc or .zshrc)
export PATH="$HOME/bin:$PATH"
```

Windows (PowerShell Admin):
```powershell
# Create directory for binary
mkdir "$env:USERPROFILE\bin" -ErrorAction SilentlyContinue
# Copy binary as uve-bin
cp uve-bin.exe "$env:USERPROFILE\bin\uve-bin.exe"
# Add to PATH (permanent)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\bin",
    "User"
)
```

### Install shell integration:

Bash/Zsh:
```bash
# Copy shell script
cp uve.sh ~/.uve.sh
# Add to shell config (choose appropriate file)
echo 'source ~/.uve.sh' >> ~/.bashrc  # or ~/.zshrc on macOS
```

PowerShell:
```powershell
# Create PowerShell modules directory
$modulesDir = "$env:USERPROFILE\Documents\WindowsPowerShell\Modules\uve"
mkdir $modulesDir -ErrorAction SilentlyContinue
# Copy PowerShell module
cp uve.ps1 "$modulesDir\uve.psm1"

# Unblock the follow and adjust the PowerShell security policy
Unblock-File -Path "$env:USERPROFILE\Documents\WindowsPowerShell\Modules\uve\uve.psm1"
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# Import the module (for current session)
Import-Module uve

# Add auto-import to PowerShell profile (for future sessions)
if (!(Test-Path $PROFILE)) {
    New-Item -Type File -Path $PROFILE -Force
}
Add-Content $PROFILE "Import-Module uve"

# Verify installation
Get-Command uve
```

## Usage

### Command Line Interface

Create a new environment:

```bash
uve create myenv [python-version]
```

List all environments:

```bash
uve list
```

Activate an environment:

```bash
uve activate myenv
```

Deactivate current environment:

```bash
uve deactivate
```

Show version of uve you are running:

```base
uve version
```

### Examples

Create a new environment with Python 3.12:

```bash 
uve create ai-stuff 3.12

# output:
Using CPython 3.12.8
Created environment 'ai-stuff' at /home/rmcdermo/.uve/ai-stuff
```

List available environments:    

```bash
uve list

# output:
Available environments:
  - ai-stuff
  - base
  - test123
  - web-stuff
```

Activate an environment:

```bash
uve activate web-stuff

# output:
(web-stuff) rmcdermo:~/mycode/uve$
```

Verify activation

```bash
which python  # or `Get-Command python` on PowerShell

# output:
/home/rmcdermo/.uve/web-stuff/bin/python
```

Install packages in the activated environment using UV's pip command

```bash
uv pip install numpy pandas

# output:
Using Python 3.11.11 environment at: .uve/web-stuff
Resolved 6 packages in 763ms
Prepared 3 packages in 1.31s
Installed 6 packages in 239ms
 + numpy==2.2.3
 + pandas==2.2.3
 + python-dateutil==2.9.0.post0
 + pytz==2025.1
 + six==1.17.0
 + tzdata==2025.1
```

View installed packages

```bash
uv pip list

#output:
Using Python 3.11.11 environment at: .uve/web-stuff
Package         Version
--------------- -----------
numpy           2.2.3
pandas          2.2.3
python-dateutil 2.9.0.post0
pytz            2025.1
six             1.17.0
tzdata          2025.1
```

Deactivate the current environment

```bash
uve deactivate
```

## Configuration

Environments are stored in `~/.uve` by default. You can change this location by setting the `UVE_HOME` environment variable:

```bash
export UVE_HOME="/path/to/environments"
```

## Environment Structure

By default, environments are created in `~/.uve/<env-name>` with the following structure:

```
~/.uve/
├── env1/
│   ├── bin/          # Scripts and executables
│   ├── lib/          # Python packages
│   └── pyvenv.cfg    # Environment configuration
└── env2/
    └── ...
```

## How It Works

UVE manages Python virtual environments in a central location (`~/.uve` by default) and provides shell functions for activating/deactivating these environments. When you activate an environment:

1. Your original PATH is saved
2. The environment's bin/Scripts directory is added to PATH
3. VIRTUAL_ENV is set to the environment path

When deactivating:

1. Original PATH is restored
2. Environment variables are cleaned up

## Troubleshooting

### Common Issues

1. `command not found: uve`
   - Ensure the binary is in your PATH
   - Restart your terminal or source your shell config

2. Activation not working
   - For bash/zsh: Make sure you sourced `~/.uve.sh`
   - For PowerShell: Ensure you imported the module

3. UV not found
   - Install [UV](https://github.com/astral-sh/uv)


4. Package installation errors
   - Make sure you've activated the environment first
   - Use `uv pip install` instead of regular `pip install`

### Environment Variables

- `UVE_HOME`: Custom location for environments
- `VIRTUAL_ENV`: Set by activation (don't set manually)
- `UVE_OLD_PATH`: Stores original PATH during activation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

- [UV](https://github.com/astral-sh/uv) - The underlying Python package installer
- Inspired by conda's environment management approach



