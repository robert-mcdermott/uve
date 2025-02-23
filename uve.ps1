# Store the original prompt function if not already stored
if (-not (Test-Path Variable:UVE_OLD_PROMPT)) {
    $Global:UVE_OLD_PROMPT = $function:prompt
}

function uve_activate {
    param([string]$envName)
    
    if ([string]::IsNullOrEmpty($envName)) {
        Write-Error "Error: Environment name required"
        return
    }
    
    $activateScript = (uve activate $envName)
    if ($LASTEXITCODE -eq 0) {
        Invoke-Expression $activateScript
        
        # Update prompt to show environment name
        $Global:UVE_ACTIVE_ENV = $envName
        $function:prompt = {
            "($Global:UVE_ACTIVE_ENV) $($Global:UVE_OLD_PROMPT.InvokeReturnAsIs())"
        }
    }
}

function uve_deactivate {
    $deactivateScript = (uve deactivate)
    if ($LASTEXITCODE -eq 0) {
        Invoke-Expression $deactivateScript
        
        # Restore original prompt
        $function:prompt = $Global:UVE_OLD_PROMPT
        Remove-Variable -Name UVE_ACTIVE_ENV -Scope Global -ErrorAction SilentlyContinue
    }
}

# Export the functions
Export-ModuleMember -Function uve_activate, uve_deactivate 