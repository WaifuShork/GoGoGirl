# Clean the previous build
if (Test-Path -Path main.exe -PathType Leaf) 
{
    Remove-Item main.exe -Force    
}

# Output for luls
Write-Output("Building GoGoGirl Bot");

# Build the modules first
go.exe build ./cmd ./framework

# Build the executable
go.exe build -buildmode=exe  main.go;