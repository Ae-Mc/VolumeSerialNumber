package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Ae-Mc/volumeID"
)

func help() {
	fmt.Println(`Examples of usage to change volume serial number:
	` + os.Args[0] + ` C: 5678-9ABC
Above command changes volume serial number of drive C to 0000-0000-5678-9ABC if C is NTFS volume or to 5678-9ABC if C is FAT12/FAT16/FAT32/EXFAT volume
It's also possible to change NTFS volume serial number to full 8 byte value:
	` + os.Args[0] + ` C: 1234-CDEF-5678-9ABC
If C is FAT32 (or similar file system) volume, than it's serial will be changed to 5678-9ABC
Examples of usage to get volume serial number:
	` + os.Args[0] + ` C:
Possible output of above command is 1234-AB12`,
	)
}

func serialNumberToString(volume_sn uint64) string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, volume_sn)
	var volume_sn_bytes [8]byte = [8]byte(buf.Bytes())
	return fmt.Sprintf(
		"%02X%02X-%02X%02X-%02X%02X-%02X%02X",
		volume_sn_bytes[0],
		volume_sn_bytes[1],
		volume_sn_bytes[2],
		volume_sn_bytes[3],
		volume_sn_bytes[4],
		volume_sn_bytes[5],
		volume_sn_bytes[6],
		volume_sn_bytes[7],
	)

}

func printVolumeSerialNumber(drive string) {
	volume_sn, err := volumeID.GetVolumeSerialNumber(drive)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println(serialNumberToString(volume_sn))
}

func setVolumeSerialNumber(drive string, volume_sn uint64) {
	err := volumeID.SetVolumeSerialNumber(drive, volume_sn)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func userInputToSerialNumber(user_input string) (volume_sn uint64, err error) {
	volume_sn_str := strings.ToUpper(strings.ReplaceAll(user_input, "-", ""))
	regular_expr, err := regexp.Compile("[0-9A-Fa-f]{8,16}")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	matched := regular_expr.FindStringIndex(volume_sn_str)[0] == 0
	if !matched {
		err = fmt.Errorf("wrong serial number format")
		return
	}
	volume_sn_bytes, err := hex.DecodeString(volume_sn_str)
	if err != nil {
		err = fmt.Errorf("wrong serial number format")
		return
	}
	if len(volume_sn_bytes) == 4 {
		volume_sn_bytes = bytes.Join(
			[][]byte{make([]byte, 4), volume_sn_bytes},
			[]byte{},
		)
	}
	err = binary.Read(
		bytes.NewReader(volume_sn_bytes),
		binary.BigEndian,
		&volume_sn,
	)
	if err != nil {
		err = fmt.Errorf("wrong serial number format")
		return
	}
	return
}

func main() {
	if len(os.Args) != 3 && len(os.Args) != 2 {
		help()
		return
	}
	drive := fmt.Sprintf("\\\\.\\%c:", os.Args[1][0])
	drive_type := volumeID.GetDriveType(drive + "\\")

	is_supported := drive_type == volumeID.DRIVE_TYPE_FIXED ||
		drive_type == volumeID.DRIVE_TYPE_REMOVABLE
	if !is_supported {
		fmt.Println("Drive type", drive_type, "is unsupported")
		fmt.Println(
			"Check https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getdrivetypew for more information",
		)
		return
	}
	if len(os.Args) == 2 {
		printVolumeSerialNumber(drive)
	} else {
		volume_sn, err := userInputToSerialNumber(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		setVolumeSerialNumber(drive, volume_sn)
	}
}
