package volumeID

import "golang.org/x/sys/windows"

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
