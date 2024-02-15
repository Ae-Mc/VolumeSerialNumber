package main

import (
	"sort"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func get_drives_list() (drivesList []string) {
	var logicalDrives [256]uint16
	_, err := windows.GetLogicalDriveStrings(
		uint32(cap(logicalDrives)),
		&logicalDrives[0],
	)
	noErr(err)
	drivesList = strings.FieldsFunc(
		strings.Trim(string(utf16.Decode(logicalDrives[:])), string(rune(0))),
		func(r rune) bool { return r == 0 },
	)
	sort.Slice(
		drivesList,
		func(i, j int) bool { return drivesList[i] < drivesList[j] },
	)
	return
}
