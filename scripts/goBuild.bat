@echo off

rem BUILD ALL "cmd" executables
rem (the resulting executables are in the current folder)
@echo on
go clean -cache ./...
go build -o ./output/ -v ./...
@echo off
if %ERRORLEVEL% GEQ 1 echo !!!!! Finished !!!!!
rem dir *.exe