package main

import (
	"fmt"
	"slices"
)

func main() {
	drives, err := getDrivesList(true)
	noErr(err)

	supported_drives := filter(drives, func(drive string) bool {
		return slices.Contains(
			[]int{DRIVE_TYPE_FIXED, DRIVE_TYPE_REMOVABLE},
			getDriveType(drive),
		)
	})

	for _, drive := range supported_drives {
		drive = fmt.Sprintf("\\\\.\\%c:", drive[0])
		content, err := readDriveSector(
			drive,
			0,
			512,
		)
		noErr(err)
		file_system_name := "unknown"
		switch {
		case isNTFS(content[:]):
			file_system_name = "NTFS"
		case isFAT12([512]byte(content)):
			file_system_name = "FAT12"
		case isFAT16([512]byte(content)):
			file_system_name = "FAT16"
		case isFAT32([512]byte(content)):
			file_system_name = "FAT32"
		case isEXFAT([512]byte(content)):
			file_system_name = "EXFAT"
		}
		fmt.Println("Drive", drive, "file system is", file_system_name)
	}
}
