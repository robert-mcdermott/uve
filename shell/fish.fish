#!/usr/bin/env fish

function uve
    switch $argv[1]
        case "activate"
            if test (count $argv) -lt 2
                echo "Error: Environment name required"
                return 1
            end
            
            if not test -d "$HOME/.uve/$argv[2]"
                echo "Environment '$argv[2]' does not exist"
                return 1
            end

            # Evaluate the activation commands directly
            eval (uve-bin activate $argv[2])
            
            # Only modify prompt if activation succeeded
            if set -q VIRTUAL_ENV
                # Use the explicitly provided environment name ($argv[2]) instead of basename
                set -g __uve_env_name $argv[2]
                functions -c fish_prompt __fish_original_prompt
                function fish_prompt
                    echo -n "($__uve_env_name) "
                    __fish_original_prompt
                end
            end
            
        case "deactivate"
            # Evaluate the deactivation commands directly
            eval (uve-bin deactivate)
            
            # Restore original prompt if it was changed
            if functions -q __fish_original_prompt
                functions -e fish_prompt
                functions -c __fish_original_prompt fish_prompt
                functions -e __fish_original_prompt
            end
            set -e __uve_env_name  # Clean up our global variable
                
        
        case "delete"
            if test (count $argv) -lt 2
                echo "Error: Environment name required"
                return 1
            end
            
            # Check if trying to delete active environment
            if set -q VIRTUAL_ENV
                set current_env (basename $VIRTUAL_ENV)
                if test "$current_env" = "$argv[2]"
                    echo "Error: Cannot delete active environment. Deactivate it first."
                    return 1
                end
            end
            
            uve-bin delete $argv[2]
            
        case '*'
            uve-bin $argv
    end
end