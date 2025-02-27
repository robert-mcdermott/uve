// Package main implements a Python virtual environment manager using uv
// It provides commands to create, activate, deactivate and list virtual environments
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Version information
const (
	VERSION = "0.1.0"
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

// printUsage prints the command-line usage instructions
func printUsage() {
	fmt.Printf(`Usage: uve <command> [args]

Commands:
  create <name> [python-version]  Create a new environment
  activate <name>                 Print activation script for environment
  deactivate                     Print deactivation script
  list                          List all environments
  version                       Show version information
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

	case "list":
		listEnvs()

	case "version":
		printVersion()

	default:
		printUsage()
		os.Exit(1)
	}
}
