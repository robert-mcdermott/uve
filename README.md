# UVE - UV Environment Manager

UVE is a lightweight Python virtual environment manager that works with [UV](https://github.com/astral-sh/uv). It provides a conda-like experience for managing Python virtual environments from any directory, while leveraging UV's speed and reliability.

## Why?

While I appreciate UV for its clean, per-project virtual environments, it's still convenient at times to have long-lived, general-purpose conda style environments that you can activate from anywhere that aren't tied to an organized project, for general purpose hacking. Since I've completely switched from conda to UV, I created this companion utility to replicate conda-like workflows when needed—giving me the best of both worlds.

## Features

- Central storage of virtual environments (similar to conda)
- Create and manage Python environments from any directory
- Cross-platform support (Linux, macOS, Windows)
- Written in Go - no Python dependency for the manager itself
- Simple shell integration for environment activation/deactivation
- Uses UV for fast environment creation
- Automatic shell detection and integration setup

## Prerequisites

- [Go](https://golang.org/doc/install)  - Optional, only needed if building from source
- [UV](https://github.com/astral-sh/uv) - Installed and available in PATH

## Installation

### Installing from source:

If you are building it yourself, clone the repo and build the binary:

```bash
git clone https://github.com/robert-mcdermott/uve.git
cd uve
go build -ldflags="-s -w" -o uve-bin main.go
```

### Installation from pre-built binary:

If you'd rather just use a pre-built binary, download the archive for your platform (OS and architecture) from the [releases page](https://github.com/robert-mcdermott/uve/releases)
and extract the archive to a directory of your choice. Unless you have a reason to select a specific release, choose the latest release.

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

Windows (PowerShell):
```powershell
# Create directory for binary
mkdir "$env:USERPROFILE\bin" -ErrorAction SilentlyContinue
# Copy the uve-bin.exe binary to your profile
cp uve-bin.exe "$env:USERPROFILE\bin\uve-bin.exe"
# Remove the MOTW flag from the file 
Unblock-File -Path "$env:USERPROFILE\bin\uve-bin.exe"
# Add to PATH (permanent, user level)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\bin",
    "User"
)
```

### Set up shell integration:

UVE includes automatic shell integration with all shell scripts embedded in the binary. Simply run:

```bash
uve-bin init
```

#### *Note: running 'uve-bin' should only be used for the initial shell integration when first installing uve, after shell integration is complete, you'll only use the 'uve' command for all other operations*

This will:
- Detect your shell type (bash, zsh, or PowerShell)
- Create the appropriate shell integration file
- Update your shell configuration to load UVE automatically

After running `uve-bin init`, either restart your shell or source your shell configuration file to activate UVE:

Bash:
```bash
source ~/.bashrc
```

Zsh:
```bash
source ~/.zshrc
```

PowerShell:
```powershell
Import-Module uve
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

Delete an environment:

```bash
uve delete myenv
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

Delete an environment (must be deactivated first):

```bash
uve delete test123

# output:
Deleted environment 'test123'
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
   - Run `uve-bin init` to set up shell integration
   - Restart your terminal or source your shell config

2. Activation not working
   - Make sure you've run `uve-bin init` and restarted your shell
   - For PowerShell: Run `Import-Module uve` in the current session

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



