# UVE 0.1.4 Release Notes

## Overview

UVE 0.1.4 introduces significant improvements to the codebase structure and distribution process, making the tool more maintainable and easier to distribute. This release focuses on technical improvements while maintaining full compatibility with previous versions.

## What's New

- **New `init` Subcommand**: Added a new `uve-bin init` command that automatically detects your shell type and sets up the appropriate shell integration. No more manual copying of shell scripts!
- **Automatic Shell Detection**: UVE now automatically detects whether you're using Bash, Zsh, or PowerShell and configures the appropriate integration.
- **Simplified Setup Process**: The entire installation process is now just: download, place in PATH, run `uve-bin init`, and you're ready to go.
- **Self-contained Binary**: Shell integration scripts are now embedded directly into the binary using Go's embed directive, eliminating the need for separate script files during installation.
- **Simplified Installation**: The new "init" shell interation automatically creates the shell scripts and sets up shell integration for you.
- **SHA256 Checksums**: Added SHA256 checksums for all distribution packages to enhance security.
- **Improved Documentation**: Updated installation instructions to reflect the new self-contained binary approach.
- **Mark of the Web Handling**: Added explicit instructions to unblock downloaded executables on Windows using `Unblock-File`.


## Compatibility

UVE 0.1.4 is fully compatible with previous versions. No changes to existing environments or workflows are required.

## Installation

Download the appropriate package for your platform from the [releases page](https://github.com/robert-mcdermott/uve/releases/tag/0.1.4), extract it, and follow the installation instructions in the README.

### Quick Installation Steps

1. Download and extract the appropriate package for your platform
2. Copy the `uve-bin` executable to a directory in your PATH
3. Run `uve-bin init` to set up shell integration
4. Restart your shell or source your shell configuration file
5. Start using UVE with the `uve` command!

## Supported Platforms

- Windows (x86_64)
- macOS (Intel x86_64 and Apple Silicon ARM64)
- Linux (x86_64)

## Known Issues

None at this time.

---

Thank you for using UVE! If you encounter any issues or have suggestions, please submit them on the [GitHub repository](https://github.com/robert-mcdermott/uve/issues). 
