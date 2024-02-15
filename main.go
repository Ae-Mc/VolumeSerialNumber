package main

import (
	"fmt"
)

func main() {
	drive := "\\\\.\\C:"
	volumeExtents, err := get_volume_disk_extents(drive)
	noErr(err)
	if volumeExtents.NumberOfDiskExtents == 0 {
		panic("number of disk extents is 0")
	}
	physicalDrive := "\\\\.\\PHYSICALDRIVE" + string(
		volumeExtents.Extents[0].DiskNumber,
	)
	content, err := read_drive_sector(
		physicalDrive,
		int64(volumeExtents.Extents[0].StartingOffset),
	)
	noErr(err)
	fmt.Println(content)
	// var readBytesCount uint32
	// err = windows.ReadFile(fileHandle, buffer[:], &readBytesCount, nil)
	// noErr(err)
	// fmt.Println(readBytesCount)
}

func noErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
