@echo off
go build -ldflags "-H windowsgui" -o classpaper-win.exe . 
pause