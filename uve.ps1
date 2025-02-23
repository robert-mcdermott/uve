# Store the original prompt function if not already stored
if (-not (Test-Path Variable:UVE_OLD_PROMPT)) {
    $Global:UVE_OLD_PROMPT = $function:prompt
}

# Main uve command function
function uve {
    param([string]$command, [string]$envName)
    
    switch ($command) {
        "activate" {
            if ([string]::IsNullOrEmpty($envName)) {
                Write-Error "Error: Environment name required"
                return
            }
            $activateScript = (uve-bin activate $envName)
            if ($LASTEXITCODE -eq 0) {
                Invoke-Expression $activateScript
                
                # Update prompt to show environment name
                $Global:UVE_ACTIVE_ENV = $envName
                $function:prompt = {
                    "($Global:UVE_ACTIVE_ENV) $($Global:UVE_OLD_PROMPT.InvokeReturnAsIs())"
                }
            }
        }
        "deactivate" {
            $deactivateScript = (uve-bin deactivate)
            if ($LASTEXITCODE -eq 0) {
                Invoke-Expression $deactivateScript
                
                # Restore original prompt
                $function:prompt = $Global:UVE_OLD_PROMPT
                Remove-Variable -Name UVE_ACTIVE_ENV -Scope Global -ErrorAction SilentlyContinue
            }
        }
        default {
            # Pass all other commands to the binary
            uve-bin $args
        }
    }
}

# Export the function
Export-ModuleMember -Function uve 