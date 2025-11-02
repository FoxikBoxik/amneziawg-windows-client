/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package version

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
)

var (
	versionapi  = windows.NewLazyDLL("version.dll")
	getFileInfo = versionapi.NewProc("GetFileVersionInfoW")
	queryValue  = versionapi.NewProc("VerQueryValueW")
)

type vsFixedFileInfo struct {
	dwSignature        uint32
	dwStrucVersion     uint32
	dwFileVersionMS    uint32
	dwFileVersionLS    uint32
	dwProductVersionMS uint32
	dwProductVersionLS uint32
	dwFileFlagsMask    uint32
	dwFileFlags        uint32
	dwFileOS           uint32
	dwFileType         uint32
	dwFileSubtype      uint32
	dwFileDateMS       uint32
	dwFileDateLS       uint32
}

func WintunVersion() string {
	// First, try to get version from running wintun
	wintunVersion, err := wintun.RunningVersion()
	if err == nil {
		return fmt.Sprintf("%d.%d", (wintunVersion>>16)&0xffff, wintunVersion&0xffff)
	}

	// If that fails, try to read version from DLL file
	exePath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	exeDir := filepath.Dir(exePath)
	dllPath := filepath.Join(exeDir, "wintun.dll")

	// Try alternative locations
	paths := []string{
		dllPath,
		filepath.Join(exeDir, "..", "wintun.dll"),
		filepath.Join(exeDir, "amd64", "wintun.dll"),
		filepath.Join(exeDir, "x86", "wintun.dll"),
		filepath.Join(exeDir, "arm64", "wintun.dll"),
		"C:\\Windows\\System32\\wintun.dll",
		"C:\\Windows\\SysWOW64\\wintun.dll",
	}

	for _, path := range paths {
		if version := getDllVersion(path); version != "" {
			return version
		}
	}

	return "unknown"
}

func getDllVersion(dllPath string) string {
	if _, err := os.Stat(dllPath); err != nil {
		return ""
	}

	pathPtr, err := syscall.UTF16PtrFromString(dllPath)
	if err != nil {
		return ""
	}

	// Get the size needed for the version info buffer
	dwLen, _, _ := getFileInfo.Call(uintptr(unsafe.Pointer(pathPtr)), 0, 0, 0)
	if dwLen == 0 {
		return ""
	}

	buf := make([]byte, dwLen)
	ret, _, _ := getFileInfo.Call(uintptr(unsafe.Pointer(pathPtr)), 0, dwLen, uintptr(unsafe.Pointer(&buf[0])))
	if ret == 0 {
		return ""
	}

	subBlock, err := syscall.UTF16FromString("\\")
	if err != nil {
		return ""
	}

	var lpffi *vsFixedFileInfo
	var uLen uintptr
	ret, _, _ = queryValue.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&subBlock[0])),
		uintptr(unsafe.Pointer(&lpffi)),
		uintptr(unsafe.Pointer(&uLen)),
	)
	if ret == 0 || lpffi == nil || uLen == 0 {
		return ""
	}

	major := (lpffi.dwFileVersionMS >> 16) & 0xffff
	minor := lpffi.dwFileVersionMS & 0xffff
	return fmt.Sprintf("%d.%d", major, minor)
}