package volumeID

import (
	"fmt"
	"io"
	"os"
)

func ReadDriveSector(
	drive string,
	offset int64,
	sectorSize uint64,
) (result []byte, err error) {
	file, err := os.Open(drive)
	if err != nil {
		return
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)
	result = make([]byte, sectorSize)
	readBytesCount, err := file.Read(result)
	if err == nil {
		if readBytesCount != int(sectorSize) {
			err = fmt.Errorf(
				"error reading sector, read %d bytes instead of 512",
				readBytesCount,
			)
		}
	}
	return
}

func WriteDriveSector(drive string, offset int64, sector []byte) (err error) {
	file, err := os.OpenFile(drive, os.O_WRONLY, os.ModeDevice)
	if err != nil {
		return
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)
	_, err = file.Write(sector)
	return
}
