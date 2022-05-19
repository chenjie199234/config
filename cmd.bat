@echo off
REM      Warning!!!!!!!!!!!This file is readonly!Don't modify this file!

cd %~dp0

where /q git.exe
if %errorlevel% == 1 (
	echo "missing dependence: git"
	goto :end
)

where /q go.exe
if %errorlevel% == 1 (
	echo "missing dependence: golang"
	goto :end
)

where /q protoc.exe
if %errorlevel% == 1 (
	echo "missing dependence: protoc"
	goto :end
)

where /q protoc-gen-go.exe
if %errorlevel% == 1 (
	echo "missing dependence: protoc-gen-go"
	goto :end
)

where /q codegen.exe
if %errorlevel% == 1 (
	echo "missing dependence: codegen"
	goto :end
)

if "%1" == "" (
	goto :help
)
if %1 == "" (
	goto :help
)
if %1 == "h" (
	goto :help
)
if "%1" == "h" (
	goto :help
)
if %1 == "-h" (
	goto :help
)
if "%1" == "-h" (
	goto :help
)
if %1 == "help" (
	goto :help
)
if "%1" == "help" (
	goto :help
)
if %1 == "-help" (
	goto :help
)
if "%1" == "-help" (
	goto :help
)
if %1 == "pb" (
	goto :pb
)
if "%1" == "pb" (
	goto :pb
)
if %1 == "kube" (
	goto :kube
)
if "%1" ==  "kube" (
	goto :kube
)
if %1 == "new" (
	if "%2" == "" (
		goto :help
	)
	if %2 == "" (
		goto :help
	)
	goto :new
)
if "%1" == "new" (
	if "%2" == "" (
		goto :help
	)
	if %2 == "" (
		goto :help
	)
	goto :new
)

:pb
	del >nul 2>nul .\api\*.pb.go
	del >nul 2>nul .\api\*.md
	go mod tidy
	for /F %%i in ('go list -m -f "{{.Dir}}" github.com/chenjie199234/Corelib') do ( set corelib=%%i )
	set workdir=%cd%
	cd %corelib%
	go install ./...
	cd %workdir%
	protoc -I ./ -I %corelib% --go_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-pbex_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-cgrpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-crpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-web_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-markdown_out=paths=source_relative:. ./api/*.proto
	go mod tidy
goto :end

:kube
	codegen -n config -p github.com/chenjie199234/config -k
goto :end

:new
	codegen -n config -p github.com/chenjie199234/config -s %2
goto :end

:help
	echo cmd.bat — every thing you need
	echo           please install git
	echo           please install golang
	echo           please install protoc           (github.com/protocolbuffers/protobuf)
	echo           please install protoc-gen-go    (github.com/protocolbuffers/protobuf-go)
	echo           please install codegen          (github.com/chenjie199234/Corelib)
	echo
	echo Usage:
	echo    ./cmd.bat <option^>
	echo
	echo Options:
	echo    pb                        Generate the proto in this program.
	echo    new <sub service name^>    Create a new sub service.
	echo    kube                      Update or add kubernetes config.
	echo    h/-h/help/-help/--help    Show this message.

:end
pause
exit /b 0
