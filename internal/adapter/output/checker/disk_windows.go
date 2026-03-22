//go:build windows

package checker

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// CheckDiskSpace uses GetDiskFreeSpaceExW to check available disk space on Windows.
func (c *SystemChecker) CheckDiskSpace(_ context.Context, minGB int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.MustFindProc("GetDiskFreeSpaceExW")

	dirPtr, err := syscall.UTF16PtrFromString(home)
	if err != nil {
		return fmt.Errorf("converting path: %w", err)
	}

	var freeBytesAvailable, totalBytes, totalFreeBytes uint64
	ret, _, callErr := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(dirPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)
	if ret == 0 {
		return fmt.Errorf("GetDiskFreeSpaceExW failed: %w", callErr)
	}

	freeGB := freeBytesAvailable / (1024 * 1024 * 1024)
	if int(freeGB) < minGB {
		return fmt.Errorf("insufficient disk space: %dGB free, need %dGB", freeGB, minGB)
	}
	return nil
}
