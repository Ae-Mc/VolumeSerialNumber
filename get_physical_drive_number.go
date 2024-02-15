package main

// #include <windows.h>
import "C"
import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func get_volume_disk_extents(
	logicalDrivePath string,
) (volumeExtents VOLUME_DISK_EXTENTS, err error) {
	driveHandle := open_drive_file(logicalDrivePath)
	defer windows.CloseHandle(driveHandle)
	var bytesReturned uint32
	err = windows.DeviceIoControl(
		driveHandle,
		C.IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS,
		nil,
		0,
		(*byte)(unsafe.Pointer(&volumeExtents)),
		uint32(unsafe.Sizeof(volumeExtents)),
		&bytesReturned,
		nil,
	)
	return volumeExtents, err
}
