# If the original prompt function isn't stored, store it now
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
Export-ModuleMember -Function uve 