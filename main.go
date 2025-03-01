// Package main implements a Python virtual environment manager using uv
// It provides commands to create, activate, deactivate and list virtual environments
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Version information
const (
	VERSION = "0.1.1"
)

// Environment variables used by the virtual environment manager
const (
	UVE_HOME_ENV     = "UVE_HOME"     // Directory where virtual environments are stored
	VIRTUAL_ENV_ENV  = "VIRTUAL_ENV"  // Currently active virtual environment
	PATH_ENV         = "PATH"         // System PATH variable
	OLD_PATH_ENV     = "UVE_OLD_PATH" // Backup of PATH before activation
	SHELL_ENV        = "SHELL"        // User's shell
	DEFAULT_UVE_HOME = ".uve"         // Default directory name for storing environments
)

// Shell integration scripts embedded in the binary
const (
	// Bash/Zsh shell integration script
	BASH_INTEGRATION = `#!/bin/bash

# Store the original prompt if not already stored
if [ -z "$UVE_OLD_PS1" ]; then
    export UVE_OLD_PS1="$PS1"
fi

# Main uve command function
uve() {
    case "$1" in
        "activate")
            if [ -z "$2" ]; then
                echo "Error: Environment name required"
                return 1
            fi
            eval "$(uve-bin activate "$2")"
            # Update prompt to show environment name
            export PS1="($2) $UVE_OLD_PS1"
            ;;
        "deactivate")
            eval "$(uve-bin deactivate)"
            # Restore original prompt
            export PS1="$UVE_OLD_PS1"
            ;;
        "delete")
            if [ -z "$2" ]; then
                echo "Error: Environment name required"
                return 1
            fi
            # Check if trying to delete the active environment
            if [ -n "$VIRTUAL_ENV" ] && [ "$(basename "$VIRTUAL_ENV")" = "$2" ]; then
                echo "Error: Cannot delete active environment. Deactivate it first."
                return 1
            fi
            # Pass to the binary
            uve-bin delete "$2"
            ;;
        *)
            # Pass all other commands to the binary
            uve-bin "$@"
            ;;
    esac
}

# Check if the script is being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "This script must be sourced. Use:"
    echo "  source ~/.uve.sh"
    exit 1
fi`

	// PowerShell module script
	POWERSHELL_MODULE = `# If the original prompt function isn't stored, store it now
if (-not (Test-Path Variable:UVE_OLD_PROMPT)) {
    $Global:UVE_OLD_PROMPT = $function:prompt
}

function uve {
    [CmdletBinding()]
    param(
        # Capture all arguments into an array
        [Parameter(ValueFromRemainingArguments=$true)]
        [string[]] $AllArgs
    )

    # If no arguments were provided, just call uve-bin with no args
    if (-not $AllArgs) {
        uve-bin
        return
    }

    switch ($AllArgs[0]) {
        "activate" {
            if ($AllArgs.Count -lt 2) {
                Write-Error "Error: Environment name required for 'activate'."
                return
            }

            $envName = $AllArgs[1]
            # Call uve-bin activate <envName>, capture output as a single string
            $activateScript = uve-bin activate $envName | Out-String

            if ($LASTEXITCODE -eq 0) {
                # Invoke the single string script to modify the session
                Invoke-Expression $activateScript

                # Update prompt to show the active environment
                $Global:UVE_ACTIVE_ENV = $envName
                $function:prompt = {
                    "($Global:UVE_ACTIVE_ENV) $($Global:UVE_OLD_PROMPT.InvokeReturnAsIs())"
                }
            }
        }
        "deactivate" {
            # Call uve-bin deactivate, also capture as a single string
            $deactivateScript = uve-bin deactivate | Out-String

            if ($LASTEXITCODE -eq 0) {
                Invoke-Expression $deactivateScript

                # Restore original prompt
                $function:prompt = $Global:UVE_OLD_PROMPT
                Remove-Variable -Name UVE_ACTIVE_ENV -Scope Global -ErrorAction SilentlyContinue
            }
        }
        "delete" {
            if ($AllArgs.Count -lt 2) {
                Write-Error "Error: Environment name required for 'delete'."
                return
            }
            
            $envName = $AllArgs[1]
            # Check if trying to delete the active environment
            if ($Global:UVE_ACTIVE_ENV -eq $envName) {
                Write-Error "Error: Cannot delete active environment. Deactivate it first."
                return
            }
            
            # Pass to the binary
            uve-bin delete $envName
        }
        default {
            # Pass all arguments along to uve-bin
            uve-bin $AllArgs
        }
    }
}

# Export the function so it's visible as a module command
Export-ModuleMember -Function uve`

	// ZSH specific integration (if needed beyond bash)
	ZSH_INTEGRATION = `#!/bin/zsh

# Store the original prompt if not already stored
if [ -z "$UVE_OLD_PS1" ]; then
    export UVE_OLD_PS1="$PS1"
fi

# Main uve command function
uve() {
    case "$1" in
        "activate")
            if [ -z "$2" ]; then
                echo "Error: Environment name required"
                return 1
            fi
            eval "$(uve-bin activate "$2")"
            # Update prompt to show environment name
            export PS1="($2) $UVE_OLD_PS1"
            ;;
        "deactivate")
            eval "$(uve-bin deactivate)"
            # Restore original prompt
            export PS1="$UVE_OLD_PS1"
            ;;
        "delete")
            if [ -z "$2" ]; then
                echo "Error: Environment name required"
                return 1
            fi
            # Check if trying to delete the active environment
            if [ -n "$VIRTUAL_ENV" ] && [ "$(basename "$VIRTUAL_ENV")" = "$2" ]; then
                echo "Error: Cannot delete active environment. Deactivate it first."
                return 1
            fi
            # Pass to the binary
            uve-bin delete "$2"
            ;;
        *)
            # Pass all other commands to the binary
            uve-bin "$@"
            ;;
    esac
}

# Check if the script is being sourced
if [[ "${(%):-%x}" == "${0}" ]]; then
    echo "This script must be sourced. Use:"
    echo "  source ~/.uve.sh"
    exit 1
fi`
)

// getUveHome returns the path where virtual environments are stored.
// It checks for UVE_HOME environment variable first, then falls back to ~/.uve
func getUveHome() string {
	if uveHome := os.Getenv(UVE_HOME_ENV); uveHome != "" {
		return uveHome
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, DEFAULT_UVE_HOME)
}

// ensureUveHome ensures the virtual environments directory exists and returns its path.
// Creates the directory if it doesn't exist.
func ensureUveHome() string {
	uveHome := getUveHome()
	if err := os.MkdirAll(uveHome, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating UVE_HOME directory: %v\n", err)
		os.Exit(1)
	}
	return uveHome
}

// createEnv creates a new Python virtual environment using uv.
// Parameters:
//   - name: name of the virtual environment
//   - pythonVersion: optional Python version to use (empty string for default)
func createEnv(name string, pythonVersion string) {
	uveHome := ensureUveHome()
	envPath := filepath.Join(uveHome, name)

	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Environment '%s' already exists\n", name)
		os.Exit(1)
	}

	args := []string{"venv"}
	if pythonVersion != "" {
		args = append(args, "--python", pythonVersion)
	}
	args = append(args, envPath)

	cmd := exec.Command("uv", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating environment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created environment '%s' at %s\n", name, envPath)
}

// generateActivateScript generates a shell script to activate a virtual environment.
// The script modifies PATH and sets environment variables.
// Parameters:
//   - envPath: full path to the virtual environment
//
// Returns activation script as string, with OS-specific syntax
func generateActivateScript(envPath string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`$env:UVE_OLD_PATH = $env:PATH
$env:VIRTUAL_ENV = "%s"
$env:PATH = "%s\Scripts;" + $env:PATH
`, envPath, envPath)
	}
	return fmt.Sprintf(`export UVE_OLD_PATH="$PATH"
export VIRTUAL_ENV="%s"
export PATH="%s/bin:$PATH"
`, envPath, envPath)
}

// generateDeactivateScript generates a shell script to deactivate the current virtual environment.
// The script restores the original PATH and unsets environment variables.
// Returns deactivation script as string, with OS-specific syntax
func generateDeactivateScript() string {
	if runtime.GOOS == "windows" {
		return `if ($env:UVE_OLD_PATH) {
    $env:PATH = $env:UVE_OLD_PATH
    Remove-Item Env:\UVE_OLD_PATH
}
Remove-Item Env:\VIRTUAL_ENV -ErrorAction SilentlyContinue
`
	}
	return `if [ -n "$UVE_OLD_PATH" ]; then
    export PATH="$UVE_OLD_PATH"
    unset UVE_OLD_PATH
fi
unset VIRTUAL_ENV
`
}

func printVersion() {
	fmt.Printf("uve version %s\n", VERSION)
}

// deleteEnv safely removes a virtual environment.
// It checks if the environment exists and is a valid UVE environment before deletion.
// Parameters:
//   - name: name of the virtual environment to delete
func deleteEnv(name string) {
	uveHome := getUveHome()
	envPath := filepath.Join(uveHome, name)

	// Check if environment exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Environment '%s' does not exist\n", name)
		os.Exit(1)
	}

	// Safety check: ensure the path is within UVE_HOME
	absUveHome, err := filepath.Abs(uveHome)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UVE_HOME path: %v\n", err)
		os.Exit(1)
	}

	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving environment path: %v\n", err)
		os.Exit(1)
	}

	// Ensure the environment path is a subdirectory of UVE_HOME
	if !filepath.HasPrefix(absEnvPath, absUveHome) {
		fmt.Fprintf(os.Stderr, "Security error: Environment path is outside UVE_HOME\n")
		os.Exit(1)
	}

	// Check if the environment is currently active
	if os.Getenv(VIRTUAL_ENV_ENV) == envPath {
		fmt.Fprintf(os.Stderr, "Error: Cannot delete active environment. Deactivate it first.\n")
		os.Exit(1)
	}

	// Perform the deletion
	if err := os.RemoveAll(envPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting environment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted environment '%s'\n", name)
}

// detectShell attempts to detect the user's shell
func detectShell() string {
	// On Windows, we default to PowerShell
	if runtime.GOOS == "windows" {
		return "powershell"
	}

	// On macOS, default to zsh (since Catalina)
	if runtime.GOOS == "darwin" {
		return "zsh"
	}

	// Try to get the shell from environment variable
	shell := os.Getenv(SHELL_ENV)
	if shell != "" {
		// Extract the shell name from the path
		shellName := filepath.Base(shell)
		return shellName
	}

	// Default to bash if we can't determine
	return "bash"
}

// initShellIntegration sets up shell integration based on detected shell
func initShellIntegration() {
	shell := detectShell()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	switch shell {
	case "powershell", "pwsh":
		setupPowerShellIntegration(homeDir)
	case "zsh":
		setupZshIntegration(homeDir)
	case "bash", "sh":
		setupBashIntegration(homeDir)
	default:
		fmt.Printf("Unsupported shell: %s. Using bash integration.\n", shell)
		setupBashIntegration(homeDir)
	}
}

// setupBashIntegration sets up integration for Bash shell
func setupBashIntegration(homeDir string) {
	// Write the integration script
	scriptPath := filepath.Join(homeDir, ".uve.sh")
	err := os.WriteFile(scriptPath, []byte(BASH_INTEGRATION), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing shell integration file: %v\n", err)
		os.Exit(1)
	}

	// Check if integration is already in .bashrc
	bashrcPath := filepath.Join(homeDir, ".bashrc")
	bashrcContent, err := os.ReadFile(bashrcPath)
	if err == nil {
		if strings.Contains(string(bashrcContent), "source ~/.uve.sh") {
			fmt.Println("Shell integration already set up in .bashrc")
			return
		}
	}

	// Add source command to .bashrc
	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .bashrc: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString("\n# UVE shell integration\nsource ~/.uve.sh\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating .bashrc: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Bash shell integration set up successfully.")
	fmt.Println("Please restart your shell or run 'source ~/.bashrc' to activate.")
}

// setupZshIntegration sets up integration for Zsh shell
func setupZshIntegration(homeDir string) {
	// Write the integration script
	scriptPath := filepath.Join(homeDir, ".uve.sh")
	err := os.WriteFile(scriptPath, []byte(ZSH_INTEGRATION), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing shell integration file: %v\n", err)
		os.Exit(1)
	}

	// Check if integration is already in .zshrc
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	zshrcContent, err := os.ReadFile(zshrcPath)
	if err == nil {
		if strings.Contains(string(zshrcContent), "source ~/.uve.sh") {
			fmt.Println("Shell integration already set up in .zshrc")
			return
		}
	}

	// Add source command to .zshrc
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening .zshrc: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString("\n# UVE shell integration\nsource ~/.uve.sh\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating .zshrc: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Zsh shell integration set up successfully.")
	fmt.Println("Please restart your shell or run 'source ~/.zshrc' to activate.")
}

// setupPowerShellIntegration sets up integration for PowerShell
func setupPowerShellIntegration(homeDir string) {
	// Create PowerShell modules directory
	modulesDir := filepath.Join(homeDir, "Documents", "WindowsPowerShell", "Modules", "uve")
	err := os.MkdirAll(modulesDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PowerShell modules directory: %v\n", err)
		os.Exit(1)
	}

	// Write the module file
	modulePath := filepath.Join(modulesDir, "uve.psm1")
	err = os.WriteFile(modulePath, []byte(POWERSHELL_MODULE), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing PowerShell module file: %v\n", err)
		os.Exit(1)
	}

	// Create or update PowerShell profile
	profileDir := filepath.Join(homeDir, "Documents", "WindowsPowerShell")
	profilePath := filepath.Join(profileDir, "Microsoft.PowerShell_profile.ps1")

	// Ensure profile directory exists
	err = os.MkdirAll(profileDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PowerShell profile directory: %v\n", err)
		os.Exit(1)
	}

	// Check if module import is already in profile
	var profileContent []byte
	profileContent, err = os.ReadFile(profilePath)
	if err == nil {
		if strings.Contains(string(profileContent), "Import-Module uve") {
			fmt.Println("PowerShell integration already set up in profile")
			fmt.Println("To use UVE in the current session, run: Import-Module uve")
			return
		}
	}

	// Add module import to profile
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PowerShell profile: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString("\n# UVE shell integration\nImport-Module uve\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating PowerShell profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("PowerShell integration set up successfully.")
	fmt.Println("To use UVE in the current session, run: Import-Module uve")
	fmt.Println("UVE will be automatically available in new PowerShell sessions.")
}

// printUsage prints the command-line usage instructions
func printUsage() {
	fmt.Printf(`Usage: uve <command> [args]

Commands:
  create <name> [python-version]  Create a new environment
  activate <name>                 Print activation script for environment
  deactivate                      Print deactivation script
  delete <name>                   Delete an environment
  list                            List all environments
  init                            Set up shell integration (auto-detects shell)
  version                         Show version information
`)
}

// listEnvs prints all available virtual environments in UVE_HOME
func listEnvs() {
	uveHome := getUveHome()
	entries, err := os.ReadDir(uveHome)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No environments found")
			return
		}
		fmt.Fprintf(os.Stderr, "Error reading environments: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available environments:")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}
}

// main is the entry point for the uve command-line tool.
// It parses command-line arguments and dispatches to appropriate functions.
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle different commands
	switch os.Args[1] {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: Environment name required")
			os.Exit(1)
		}
		pythonVersion := ""
		if len(os.Args) > 3 {
			pythonVersion = os.Args[3]
		}
		createEnv(os.Args[2], pythonVersion)

	case "activate":
		if len(os.Args) < 3 {
			fmt.Println("Error: Environment name required")
			os.Exit(1)
		}
		envPath := filepath.Join(getUveHome(), os.Args[2])
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Environment '%s' does not exist\n", os.Args[2])
			os.Exit(1)
		}
		fmt.Print(generateActivateScript(envPath))

	case "deactivate":
		fmt.Print(generateDeactivateScript())

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Environment name required")
			os.Exit(1)
		}
		deleteEnv(os.Args[2])

	case "list":
		listEnvs()

	case "init":
		initShellIntegration()

	case "version":
		printVersion()

	default:
		printUsage()
		os.Exit(1)
	}
}
