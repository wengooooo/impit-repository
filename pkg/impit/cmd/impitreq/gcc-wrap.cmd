@echo off
setlocal EnableDelayedExpansion

set "REAL_GCC=C:\w64devkit\bin\gcc.exe"
set "OBJCOPY=C:\w64devkit\bin\objcopy.exe"

set "OUT="
set "DO_CONVERT=0"

for %%A in (%*) do (
  if "%%~A"=="-c" set "DO_CONVERT=1"
)

:parse
if "%~1"=="" goto run
set "ARG=%~1"
if "%ARG:~0,2%"=="-o" (
  if not "%ARG%"=="-o" set "OUT=%ARG:~2%"
)
if "%~1"=="-o" (
  set "OUT=%~2"
  shift
)
shift
goto parse

:run
"%REAL_GCC%" %*
set "RC=%ERRORLEVEL%"
if not "%RC%"=="0" exit /b %RC%

if "%OUT%"=="" goto end

if /I "%OUT:~-2%"==".o" set "DO_CONVERT=1"
if /I "%OUT:~-4%"==".obj" set "DO_CONVERT=1"

if "%DO_CONVERT%"=="1" (
  if defined IMPIT_GCC_WRAP_DEBUG echo [gcc-wrap] convert "%OUT%" 1>&2
  "%OBJCOPY%" -O pe-x86-64 "%OUT%" "%OUT%.pe"
  if errorlevel 1 exit /b 1
  move /Y "%OUT%.pe" "%OUT%" >nul
)

:end
exit /b %RC%
