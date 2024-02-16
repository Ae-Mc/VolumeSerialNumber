package volumeID

import (
	"slices"
	"strings"
	"unicode/utf16"

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

func GetDriveType(drivePath string) int {
	utf16DrivePath, err := windows.UTF16FromString(drivePath)
	if err != nil {
		panic(err)
	}
	driveType := windows.GetDriveType(&utf16DrivePath[0])
	return int(driveType)
}

func GetDrivesList(should_sort bool) (drives_list []string, err error) {
	var logical_drives [256]uint16
	_, err = windows.GetLogicalDriveStrings(
		uint32(cap(logical_drives)),
		&logical_drives[0],
	)
	if err == nil {
		drives_list = strings.FieldsFunc(
			string(utf16.Decode(logical_drives[:])),
			func(r rune) bool { return r == 0 },
		)
		drives_list = drives_list[:len(drives_list)-1]
		if should_sort {
			slices.Sort(drives_list)
		}
	}
	return
}
