#!/bin/zsh

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
            export UV_PROJECT_ENVIRONMENT="$VIRTUAL_ENV"
            # Update prompt to show environment name
            export PS1="($2) $UVE_OLD_PS1"
            ;;
        "deactivate")
            eval "$(uve-bin deactivate)"
            unset UV_PROJECT_ENVIRONMENT
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

# Sync UV_PROJECT_ENVIRONMENT on session start if VIRTUAL_ENV already set
if [ -n "$VIRTUAL_ENV" ]; then
    export UV_PROJECT_ENVIRONMENT="$VIRTUAL_ENV"
fi

# Check if the script is being sourced
if [[ $ZSH_EVAL_CONTEXT =~ :file$ ]]; then
    : # Script is being sourced, do nothing
else
    echo "This script must be sourced. Use:"
    echo "  source ~/.uve.sh"
    exit 1
fi 