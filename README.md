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

- [Go](https://golang.org/doc/install) (for building from source only)
- [UV](https://github.com/astral-sh/uv)

## Installation

### Clone and build:

```bash
git clone https://github.com/robert-mcdermott/uve.git
cd uve
go build -o uve main.go
```

### Install the binary:

Linux/macOS:
```bash
# Create bin directory if it doesn't exist
mkdir -p ~/bin
# Copy the binary as uve-bin
cp uve ~/bin/uve-bin
# Add to PATH (add this to your .bashrc or .zshrc)
export PATH="$HOME/bin:$PATH"
```

Windows (PowerShell Admin):
```powershell
# Create directory for binary
mkdir "$env:USERPROFILE\bin" -ErrorAction SilentlyContinue
# Copy binary as uve-bin
cp uve.exe "$env:USERPROFILE\bin/uve-bin.exe"
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
echo 'source ~/.uve.sh' >> ~/.bashrc  # or ~/.zshrc
```

PowerShell:
```powershell
# Create PowerShell modules directory
$modulesDir = "$env:USERPROFILE\Documents\PowerShell\Modules\uve"
mkdir $modulesDir -ErrorAction SilentlyContinue
# Copy PowerShell module
cp uve.ps1 "$modulesDir\uve.psm1"
```

## Usage

### Command Line Interface

```bash
# Create a new environment
uve create myenv [python-version]

# List all environments
uve list

# Activate an environment
uve_activate myenv

# Deactivate current environment
uve_deactivate
```

### Examples

```bash
# Create environment with Python 3.11
uve create ml-project 3.11

# List available environments
uve list

# Activate environment
uve_activate ml-project

# Verify activation
which python  # or `Get-Command python` on PowerShell

# Install packages using UV's pip command
uv pip install numpy pandas
uv pip install torch tensorflow
uv pip install -r requirements.txt  # Install from requirements file

# View installed packages
uv pip list

# Deactivate environment
uve_deactivate
```

## Configuration

Environments are stored in `~/.uve` by default. You can change this location by setting the `UVE_HOME` environment variable:

```bash
export UVE_HOME="/path/to/environments"
```

## How It Works

UVE manages Python virtual environments in a central location (`~/.uve` by default) and provides shell functions for activating/deactivating these environments. When you activate an environment:

1. Your original PATH is saved
2. The environment's bin/Scripts directory is added to PATH
3. VIRTUAL_ENV is set to the environment path

When deactivating:

1. Original PATH is restored
2. Environment variables are cleaned up

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


## Credits

- [UV](https://github.com/astral-sh/uv) - The underlying Python package installer
- Inspired by conda's environment management approach

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

## Roadmap

Future improvements under consideration:

- Environment removal command
- Environment cloning

