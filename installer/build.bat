@echo off
rem SPDX-License-Identifier: MIT
rem Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.

setlocal enabledelayedexpansion
set PATHEXT=.exe
set BUILDDIR=%~dp0
cd /d %BUILDDIR% || exit /b 1

for /f "tokens=3" %%a in ('findstr /r "Number.*=.*[0-9.]*" ..\version\version.go') do set WIREGUARD_VERSION=%%a
set WIREGUARD_VERSION=%WIREGUARD_VERSION:"=%

set WIX_CANDLE_FLAGS=-nologo -dWIREGUARD_VERSION="%WIREGUARD_VERSION%"
set WIX_LIGHT_FLAGS=-nologo -spdb
set WIX_LIGHT_FLAGS=%WIX_LIGHT_FLAGS% -sice:ICE39
set WIX_LIGHT_FLAGS=%WIX_LIGHT_FLAGS% -sice:ICE61
set WIX_LIGHT_FLAGS=%WIX_LIGHT_FLAGS% -sice:ICE03

if exist .deps\prepared goto build
"%~dp0\..\nhcolor.exe" 0E "[+] Installing dependencies"
rmdir /s /q .deps 2> NUL
mkdir .deps || goto error
cd .deps || goto error
call :download wix-binaries.zip https://github.com/wixtoolset/wix3/releases/download/wix3141rtm/wix314-binaries.zip 6ac824e1642d6f7277d0ed7ea09411a508f6116ba6fae0aa5f2c7daa2ff43d31 || goto error
"%~dp0\..\nhcolor.exe" 0E "[+] Extracting wix-binaries.zip"
mkdir wix\bin || goto error
tar -xf wix-binaries.zip -C wix\bin || goto error
"%~dp0\..\nhcolor.exe" 0E "[+] Cleaning up wix-binaries.zip"
del wix-binaries.zip || goto error
copy /y NUL prepared > NUL || goto error
cd .. || goto error

:build
"%~dp0\..\nhcolor.exe" 0E "[+] Starting build process"
if exist ..\sign.bat call ..\sign.bat
set PATH=%BUILDDIR%..\.deps\llvm-mingw-20231128-ucrt-x86_64\bin;%PATH%
set WIX=%BUILDDIR%.deps\wix\
set CFLAGS=-O3 -Wall -std=gnu11 -DWINVER=0x0601 -D_WIN32_WINNT=0x0601 -municode -DUNICODE -D_UNICODE -DNDEBUG
set LDFLAGS=-shared -s -Wl,--kill-at -Wl,--major-os-version=6 -Wl,--minor-os-version=1 -Wl,--major-subsystem-version=6 -Wl,--minor-subsystem-version=1 -Wl,--tsaware -Wl,--dynamicbase -Wl,--nxcompat -Wl,--export-all-symbols
set LDLIBS=-lmsi -lole32 -lshlwapi -lshell32 -luuid -lntdll

"%~dp0\..\nhcolor.exe" 0A "[+] Building for x86"
call :build_msi x86 i686 x86
if errorlevel 1 goto error

"%~dp0\..\nhcolor.exe" 0A "[+] Building for amd64"
call :build_msi amd64 x86_64 x64
if errorlevel 1 goto error

"%~dp0\..\nhcolor.exe" 0A "[+] Building for arm64"
call :build_msi arm64 aarch64 arm64
if errorlevel 1 goto error

if "%SigningProvider%"=="" goto success
if "%TimestampServer%"=="" goto success
"%~dp0\..\nhcolor.exe" 0E "[+] Signing MSI packages"
signtool sign %SigningProvider% /fd sha256 /tr "%TimestampServer%" /td sha256 /d "AmneziaWG Setup" "dist\amneziawg-*-%WIREGUARD_VERSION%.msi" || goto error

:success
"%~dp0\..\nhcolor.exe" 0A "[+] Build successful"
pause
exit /b 0

:download
"%~dp0\..\nhcolor.exe" 0E "[+] Downloading %1"
curl -#fLo %1 %2 || exit /b 1
"%~dp0\..\nhcolor.exe" 0E "[+] Verifying %1"
for /f %%a in ('CertUtil -hashfile %1 SHA256 ^| findstr /r "^[0-9a-f]*$"') do if not "%%a"=="%~3" exit /b 1
exit /b 0

:build_msi
set ARCH=%~1
set CC=%~2-w64-mingw32-gcc
set WIX_ARCH=%~3

"%~dp0\..\nhcolor.exe" 0E "[+] Creating directory for %ARCH%"
if not exist "%ARCH%" mkdir "%ARCH%"

"%~dp0\..\nhcolor.exe" 0E "[+] Compiling customactions.dll for %ARCH%"
%CC% %CFLAGS% %LDFLAGS% -o "%ARCH%\customactions.dll" customactions.c %LDLIBS% 
if errorlevel 1 (
    "%~dp0\..\nhcolor.exe" 0C "[-] Failed to compile customactions.dll for %ARCH%"
    exit /b 1
)

if not "%SigningProvider%"=="" if not "%TimestampServer%"=="" (
    "%~dp0\..\nhcolor.exe" 0E "[+] Signing customactions.dll for %ARCH%"
    signtool sign %SigningProvider% /fd sha256 /tr "%TimestampServer%" /td sha256 /d "AmneziaWG Setup Custom Actions" "%ARCH%\customactions.dll" 
    if errorlevel 1 (
        "%~dp0\..\nhcolor.exe" 0C "[-] Failed to sign customactions.dll for %ARCH%"
        exit /b 1
    )
)

"%~dp0\..\nhcolor.exe" 0E "[+] Compiling WiX objects for %ARCH%"
"%WIX%bin\candle.exe" %WIX_CANDLE_FLAGS% -dWIREGUARD_PLATFORM="%ARCH%" -out "%ARCH%\wireguard.wixobj" -arch %WIX_ARCH% wireguard.wxs
if errorlevel 1 (
    "%~dp0\..\nhcolor.exe" 0C "[-] Failed to compile WiX objects for %ARCH%"
    exit /b 1
)

"%~dp0\..\nhcolor.exe" 0E "[+] Linking MSI for %ARCH%"
"%WIX%bin\light.exe" %WIX_LIGHT_FLAGS% -out "dist\amneziawg-%ARCH%-%WIREGUARD_VERSION%.msi" "%ARCH%\wireguard.wixobj"
if errorlevel 1 (
    "%~dp0\..\nhcolor.exe" 0C "[-] Failed to link MSI for %ARCH%"
    exit /b 1
)

exit /b 0

:error
"%~dp0\..\nhcolor.exe" 0C "[-] Build failed with error #%errorlevel%"
exit /b 1