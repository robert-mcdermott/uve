#!/bin/bash

# Store the original prompt if not already stored
if [ -z "$UVE_OLD_PS1" ]; then
    export UVE_OLD_PS1="$PS1"
fi

uve_activate() {
    if [ -z "$1" ]; then
        echo "Error: Environment name required"
        return 1
    fi
    eval "$(uve activate "$1")"
    # Update prompt to show environment name
    export PS1="($1) $UVE_OLD_PS1"
}

uve_deactivate() {
    eval "$(uve deactivate)"
    # Restore original prompt
    export PS1="$UVE_OLD_PS1"
}

# Check if the script is being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "This script must be sourced. Use:"
    echo "  source uve.sh"
    exit 1
fi 