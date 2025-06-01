#!/usr/bin/env fish

function uve
    switch $argv[1]
        case "activate"
            if test (count $argv) -lt 2
                echo "Error: Environment name required"
                return 1
            end
            
            set env_name $argv[2]
            set activate_script (uve-bin activate $env_name)
            eval $activate_script
            
            # Update prompt to show environment name
            function fish_prompt
                echo -n "($env_name) "
                __fish_print_prompt
            end
            
        case "deactivate"
            set deactivate_script (uve-bin deactivate)
            eval $deactivate_script
            
            # Restore original prompt
            functions -e fish_prompt
            
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