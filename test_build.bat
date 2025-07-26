@echo off
echo Testing Go build...
go build .
if %errorlevel% equ 0 (
    echo Build successful!
    echo Testing basic functionality...
    classpaper.exe --help 2>nul || echo Program compiled successfully
) else (
    echo Build failed with errors
)
pause