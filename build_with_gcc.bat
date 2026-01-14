@echo off
echo ============================================
echo   使用正确的GCC构建 impitreq.dll
echo ============================================

REM 显示当前GCC信息
echo [1] 检查GCC...
where gcc
gcc --version

REM 查找GCC的确切路径
for /f "tokens=*" %%i in ('where gcc') do (
    set GCC_PATH=%%i
)
echo 找到GCC在: %GCC_PATH%

REM 提取GCC目录
set GCC_DIR=%~dp0
for %%i in ("%GCC_PATH%") do set GCC_DIR=%%~dpi
set GCC_DIR=%GCC_DIR:~0,-1%

echo GCC目录: %GCC_DIR%

REM 设置环境变量
echo [2] 设置环境...
set CGO_ENABLED=1
set CC=%GCC_PATH%
set PATH=%GCC_DIR%;%PATH%

REM 清理
echo [3] 清理...
go clean -cache
if exist impitreq.dll del impitreq.dll
if exist impitreq.h del impitreq.h

echo.
echo [4] 测试CGO...
echo // +build cgo > _testcgo.c
echo #include <stdio.h> >> _testcgo.c
echo void TestCGO() { printf("CGO test OK\\n"); } >> _testcgo.c

echo package main > _testcgo.go
echo // #include "_testcgo.c" >> _testcgo.go
echo import "C" >> _testcgo.go
echo import "fmt" >> _testcgo.go
echo func main() { >> _testcgo.go
echo     C.TestCGO() >> _testcgo.go
echo     fmt.Println("Go test OK") >> _testcgo.go
echo } >> _testcgo.go

go run _testcgo.go
if %ERRORLEVEL% NEQ 0 (
    echo ✗ CGO测试失败
    echo 错误代码: %ERRORLEVEL%
    goto :error
)

del _testcgo.go _testcgo.c 2>nul
echo ✓ CGO测试通过

echo.
echo [5] 构建DLL...
echo 使用编译器: %CC%
go build -x -buildmode=c-shared -o impitreq.dll pkg/impit/cmd/impitreq/main.go 2>&1 | findstr /i "gcc cgo error"

if %ERRORLEVEL% EQU 0 (
    goto :error
)

if exist impitreq.dll (
    echo.
    echo ============================================
    echo   ✓ BUILD SUCCESSFUL!
    echo ============================================
    echo.
    echo 生成的文件:
    dir impitreq.*
    
    echo.
    echo DLL导出函数:
    dumpbin /exports impitreq.dll 2>nul | findstr /c:"ordinal" /c:"name"
    
    echo.
    echo 使用方法:
    echo 1. C/C++程序: #include "impitreq.h"
    echo 2. 链接: impitreq.dll
    echo 3. 其他语言: 使用FFI加载DLL
) else (
    :error
    echo.
    echo ============================================
    echo   ✗ BUILD FAILED
    echo ============================================
    echo.
    call :show_troubleshooting
)

pause
exit /b

:show_troubleshooting
echo 故障排除步骤:
echo.
echo 1. 检查Go代码中的CGO导入:
echo   确保文件开头有: import "C"
echo   并且没有实际的C代码需要编译
echo.
echo 2. 如果代码没有C代码，尝试禁用CGO:
echo   set CGO_ENABLED=0
echo   go build -buildmode=c-shared -o impitreq.dll ...
echo.
echo 3. 检查是否有其他GCC冲突:
echo   where /r C:\ gcc.exe
echo.
echo 4. 直接运行cgo测试:
echo   cd %GOROOT%\misc\cgo\test
echo   go test
goto :eof