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
		volumeExtents, err := get_volume_disk_extents(drive)
		noErr(err)
		if volumeExtents.NumberOfDiskExtents == 0 {
			panic("number of disk extents is 0")
		}
		physicalDrive := fmt.Sprintf(
			"\\\\.\\PHYSICALDRIVE%d", volumeExtents.Extents[0].DiskNumber,
		)
		content, err := read_drive_sector(
			physicalDrive,
			int64(volumeExtents.Extents[0].StartingOffset),
			512,
		)
		noErr(err)
		fileSystemName := "unknown"
		switch {
		case isNTFS(content[:12]):
			fileSystemName = "NTFS"
		case isFAT12(content):
			fileSystemName = "FAT12"
		case isFAT16(content):
			fileSystemName = "FAT16"
		case isFAT32(content):
			fileSystemName = "FAT32"
		}
		fmt.Println("Drive", drive, "file system is", fileSystemName)
	}
}
