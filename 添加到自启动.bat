@echo off
setlocal

rem 设置要添加到用户自启动的程序名（必须在当前目录下）
set "programName=ClassPaper.exe"

rem 获取当前用户的用户文件夹路径
for /f "tokens=*" %%a in ('echo %USERPROFILE%') do set "userProfile=%%a"

rem 设置启动文件夹的路径
set "startupFolder=%userProfile%\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup"

rem 构建程序的完整路径
set "programPath=%cd%\%programName%"

rem 创建程序的快捷方式
echo Set objShell = WScript.CreateObject("WScript.Shell") > CreateShortcut.vbs
echo strDesktop = objShell.SpecialFolders("Startup") >> CreateShortcut.vbs
echo Set objShortCut = objShell.CreateShortcut(strDesktop ^& "\%~n0.lnk") >> CreateShortcut.vbs
echo objShortCut.TargetPath = "%programPath%" >> CreateShortcut.vbs
echo objShortCut.Save >> CreateShortcut.vbs
cscript CreateShortcut.vbs
del CreateShortcut.vbs

if exist "%startupFolder%\%~n0.lnk" (
    echo %programPath% 添加到用户自启动成功。
) else (
    echo %programPath% 添加到用户自启动失败。
)

endlocal
