echo off
go generate ./...
go install  -ldflags -H=windowsgui ./...
mkdir dist
xcopy .\assets .\dist\assets\  /s /e /h /y
copy .\README.md .\dist\
copy .\CHANGELOG.md .\dist\

copy .\..\..\..\..\bin\player.exe .\dist\
copy .\..\..\..\..\bin\commentator.exe .\dist\