package main

// #include <windows.h>
import "C"
import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func get_partition_information(
	drivePath string,
) (partitionInformation PARTITION_INFORMATION_EX) {
	driveHandle := open_drive_file(drivePath)
	defer windows.CloseHandle(driveHandle)
	var bytes_returned uint32
	err := windows.DeviceIoControl(
		driveHandle,
		C.IOCTL_DISK_GET_PARTITION_INFO_EX,
		nil,
		0,
		(*byte)(unsafe.Pointer(&partitionInformation)),
		uint32(unsafe.Sizeof(PARTITION_INFORMATION_EX{})),
		&bytes_returned,
		nil,
	)
	noErr(err)
	return partitionInformation
}
