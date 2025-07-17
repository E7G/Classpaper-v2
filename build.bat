@echo off
go generate
go build -ldflags "-H windowsgui" -o classpaper-win.exe . 
pause