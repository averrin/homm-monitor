echo off
go generate ./...
go install  -ldflags -H=windowsgui ./...
rmdir HoMM-Monitor
mkdir HoMM-Monitor
xcopy .\assets .\HoMM-Monitor\assets\  /s /e /h /y
xcopy .\OBS .\HoMM-Monitor\OBS\  /s /e /h /y
copy .\README.md .\HoMM-Monitor\
copy .\CHANGELOG.md .\HoMM-Monitor\
mkdir .\HoMM-Monitor\plugin
copy .\plugin\cursors.dll .\HoMM-Monitor\plugin\

copy .\..\..\..\..\bin\player.exe .\HoMM-Monitor\
copy .\..\..\..\..\bin\commentator.exe .\HoMM-Monitor\