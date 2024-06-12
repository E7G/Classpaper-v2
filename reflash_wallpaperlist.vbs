Option Explicit

Dim fso, folder, files, file, output, outputFile, backupFile
Set fso = CreateObject("Scripting.FileSystemObject")

' 定义输出文件名和备份文件名
outputFile = "res/config/wallpaperlist.js"
backupFile = "res/config/wallpaperlist_backup.js"

' 如果目标JS文件存在，则将其备份
If (fso.FileExists(outputFile)) Then
    fso.CopyFile outputFile, backupFile, True
End If

' 初始化output变量
output = "const wallpaperlist=["

' 遍历res/wallpaper下的所有图片文件
Set folder = fso.GetFolder("res/wallpaper")
Set files = folder.Files
For Each file In files
    If LCase(fso.GetExtensionName(file.Name)) = "jpg" Or LCase(fso.GetExtensionName(file.Name)) = "png" Then
        output = output & """wallpaper/" & file.Name & ""","
    End If
Next

' 移除最后一个逗号，并添加结束括号
output = Left(output, Len(output) - 1) & "];"

' 输出到目标JS文件
Set outputFile = fso.OpenTextFile(outputFile, 2, True)
outputFile.Write output
outputFile.Close
