@echo on
mkdir build
mkdir test
set GOOS=windows
set GOARCH=amd64
set FYNE_FONT=C:\Windows\Fonts\Arial.ttf
go build -o %cd%\build\DFF_win_debug.exe %cd%\cmd\dff.go
go build -o %cd%\build\DFF_win.exe -ldflags -H=windowsgui %cd%\cmd\dff.go
7z a %cd%\out\DFF_windows.zip %cd%\build\DFF_win.exe
7z a %cd%\out\DFF_windows_debug.zip %cd%\build\DFF_win_debug.exe
pause
