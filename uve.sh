#!/bin/bash

uve_activate() {
    if [ -z "$1" ]; then
        echo "Error: Environment name required"
        return 1
    fi
    eval "$(uve activate "$1")"
}

uve_deactivate() {
    eval "$(uve deactivate)"
}

# Check if the script is being sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "This script must be sourced. Use:"
    echo "  source uve.sh"
    exit 1
fi 