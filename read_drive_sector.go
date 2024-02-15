package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"
)

func read_drive_sector(
	drive string,
	offset int64,
) (result [512]byte, err error) {
	fileHandle := open_drive_file(drive)
	defer windows.CloseHandle(fileHandle)
	file, err := os.OpenFile(drive, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)
	readBytesCount, err := file.Read(result[:])
	if err != nil {
		return
	}
	if readBytesCount != 512 {
		err = fmt.Errorf(
			"error reading sector, read %d bytes instead of 512",
			readBytesCount,
		)
		return
	}
	return
}

func open_drive_file(drivePath string) windows.Handle {
	driveUtf16, err := windows.UTF16FromString(drivePath)
	noErr(err)
	fileHandle, err := windows.CreateFile(
		&driveUtf16[0],
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		windows.Handle(windows.GetShellWindow()),
	)
	noErr(err)
	return fileHandle
}
