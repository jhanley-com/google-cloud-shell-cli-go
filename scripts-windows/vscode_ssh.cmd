@echo off
cd  %~dp0
if "%1"=="" (
    type %UserProfile%\.ssh\config | find /i "host "
    exit /b
)
set vscode_host="%3"
set port="%2"
if "%vscode_host:cloudshell=%"=="%vscode_host%" (
    ssh %*
) else (
    REM cloudshell.exe ssh -D %port% --urlfetch https://****/urlfetch/ --v2ray
    cloudshell.exe ssh -D %port%  --proxy 127.0.0.1:1080 --debug
    REM cloudshell.exe ssh -D %port% -o LocalForward="127.0.0.1:22080 127.0.0.1:22080" -o ProxyCommand="C:\Program Files\Git\mingw64\bin\connect.exe -S 127.0.0.1:1080 %%h %%p"
)