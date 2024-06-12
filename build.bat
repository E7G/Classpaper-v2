@echo off
g++ cpp/setwallpaper.cpp -o setwallpaper.exe 
g++ cpp/deltaskbar.cpp -o deltaskbar.exe -lole32
go generate
go build -ldflags "-H windowsgui" .
pause