package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"
)

const (
	DRIVE_TYPE_UNKNOWN    = iota
	DRIVE_TYPE_WRONG_PATH = iota
	DRIVE_TYPE_REMOVABLE  = iota
	DRIVE_TYPE_FIXED      = iota
	DRIVE_TYPE_REMOTE     = iota
	DRIVE_TYPE_CDROM      = iota
	DRIVE_TYPE_RAMDISK    = iota
)

func getDriveType(drivePath string) int {
	utf16DrivePath, err := windows.UTF16FromString(drivePath)
	noErr(err)
	driveType := windows.GetDriveType(&utf16DrivePath[0])
	return int(driveType)
}

func readDriveSector(
	drive string,
	offset int64,
	sectorSize uint64,
) (result []byte, err error) {
	file, err := os.OpenFile(drive, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)
	result = make([]byte, sectorSize)
	readBytesCount, err := file.Read(result)
	if err != nil {
		return
	}
	if readBytesCount != int(sectorSize) {
		err = fmt.Errorf(
			"error reading sector, read %d bytes instead of 512",
			readBytesCount,
		)
		return
	}
	return
}
