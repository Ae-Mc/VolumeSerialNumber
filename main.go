package main

import (
	"fmt"
	"slices"
)

func main() {
	drives := get_drives_list()

	supportedDrives := Filter(drives, func(drive string) bool {
		return slices.Contains(
			[]int{DRIVE_TYPE_FIXED, DRIVE_TYPE_REMOVABLE},
			get_drive_type(drive),
		)
	})

	for _, drive := range supportedDrives {
		drive = fmt.Sprintf("\\\\.\\%c:", drive[0])
		content, err := read_drive_sector(
			drive,
			0,
			512,
		)
		noErr(err)
		fileSystemName := "unknown"
		switch {
		case isNTFS(content[:]):
			fileSystemName = "NTFS"
		case isFAT12([512]byte(content)):
			fileSystemName = "FAT12"
		case isFAT16([512]byte(content)):
			fileSystemName = "FAT16"
		case isFAT32([512]byte(content)):
			fileSystemName = "FAT32"
		case isEXFAT([512]byte(content)):
			fileSystemName = "EXFAT"
		}
		fmt.Println("Drive", drive, "file system is", fileSystemName)
	}
}
