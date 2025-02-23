function uve_activate {
    param([string]$envName)
    
    if ([string]::IsNullOrEmpty($envName)) {
        Write-Error "Error: Environment name required"
        return
    }
    
    $activateScript = (uve activate $envName)
    if ($LASTEXITCODE -eq 0) {
        Invoke-Expression $activateScript
    }
}

function uve_deactivate {
    $deactivateScript = (uve deactivate)
    if ($LASTEXITCODE -eq 0) {
        Invoke-Expression $deactivateScript
    }
}

# Export the functions
Export-ModuleMember -Function uve_activate, uve_deactivate 