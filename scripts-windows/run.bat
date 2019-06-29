@cd ..
go build -o cloudshell.exe
@if errorlevel 1 goto err_out

cloudshell info

@goto end

:err_out
@echo ***************************************************************
@echo Build Failed     Build Failed     Build Failed     Build Failed
@echo ***************************************************************

:end
