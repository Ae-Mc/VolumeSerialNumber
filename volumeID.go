package volumeID

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	FILE_SYSTEM_NTFS  = iota
	FILE_SYSTEM_FAT12 = FAT_TYPE_FAT12
	FILE_SYSTEM_FAT16 = FAT_TYPE_FAT16
	FILE_SYSTEM_FAT32 = FAT_TYPE_FAT32
	FILE_SYSTEM_EXFAT = FAT_TYPE_EXFAT
)

func GetFileSystem(firstDiskSector [512]byte) (file_system int, err error) {
	fat_type := getFATType(firstDiskSector)
	if fat_type == FAT_TYPE_UNKNOWN {
		if isNTFS(firstDiskSector[:]) {
			file_system = FILE_SYSTEM_NTFS
		} else {
			err = fmt.Errorf("unknown file system")
		}
	} else {
		file_system = fat_type
	}
	return
}

func GetVolumeSerialNumberSize(
	firstDiskSector [512]byte,
) (size int64, err error) {
	sizes := map[int]int64{
		FILE_SYSTEM_NTFS: 8,
		FAT_TYPE_FAT12:   4,
		FAT_TYPE_FAT16:   4,
		FAT_TYPE_FAT32:   4,
		FAT_TYPE_EXFAT:   4,
	}
	if file_system, err := GetFileSystem(firstDiskSector); err == nil {
		size = sizes[file_system]
	}
	return
}

func GetVolumeSerialNumberAddr(
	firstDiskSector [512]byte,
) (addr int64, err error) {
	addrs := map[int]int64{
		FILE_SYSTEM_NTFS: 0x48,
		FAT_TYPE_FAT12:   0x27,
		FAT_TYPE_FAT16:   0x27,
		FAT_TYPE_FAT32:   0x43,
		FAT_TYPE_EXFAT:   0x64,
	}
	if file_system, err := GetFileSystem(firstDiskSector); err == nil {
		addr = addrs[file_system]
	}
	return
}

func getVolumeSerialNumberAddrAndSize(
	firstDiskSector [512]byte,
) (addr, size int64, err error) {
	if addr, err = GetVolumeSerialNumberAddr([512]byte(firstDiskSector)); err != nil {
		return
	}
	size, err = GetVolumeSerialNumberSize([512]byte(firstDiskSector))
	return
}

func GetVolumeSerialNumber(drive string) (volume_sn uint64, err error) {
	var firstSector []byte
	if firstSector, err = ReadDriveSector(drive, 0, 512); err != nil {
		return
	}
	addr, size, err := getVolumeSerialNumberAddrAndSize([512]byte(firstSector))
	if err != nil {
		return
	}

	reader := bytes.NewReader(firstSector[:])
	reader.Seek(addr, io.SeekStart)
	if size == 4 {
		var volume_sn32 uint32
		binary.Read(reader, binary.LittleEndian, &volume_sn32)
		volume_sn = uint64(volume_sn32)
	} else {
		binary.Read(reader, binary.LittleEndian, &volume_sn)
	}

	return
}

func fillExFATAdditionalSectors(drive string, sector_size uint64) (err error) {
	var sectors []byte
	sectors, err = ReadDriveSector(drive, 0, uint64(sector_size*11))
	if err != nil {
		return
	}
	var checksum uint32
	checksum, err = exFatChecksum(sectors, uint16(sector_size))
	if err != nil {
		return
	}
	checksum_sector := make([]byte, sector_size)
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, checksum)
	checksum_sector = bytes.ReplaceAll(
		checksum_sector,
		[]byte{0, 0, 0, 0},
		buf.Bytes(),
	)
	WriteDriveSector(drive, int64(sector_size)*12, sectors)
	WriteDriveSector(drive, int64(sector_size)*11, checksum_sector)
	WriteDriveSector(drive, int64(sector_size)*23, checksum_sector)
	return
}

func SetVolumeSerialNumber(drive string, volume_sn uint64) (err error) {
	var firstSector []byte
	if firstSector, err = ReadDriveSector(drive, 0, 512); err != nil {
		return
	}
	addr, size, err := getVolumeSerialNumberAddrAndSize([512]byte(firstSector))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	if size == 4 {
		binary.Write(&buf, binary.LittleEndian, uint32(volume_sn))
	} else {
		binary.Write(&buf, binary.LittleEndian, volume_sn)
	}
	for i := addr; i < addr + size; i++ {
		firstSector[i] = buf.Bytes()[i - addr]
	}
	file_system, err := GetFileSystem([512]byte(firstSector))
	if err != nil {
		return
	}

	err = WriteDriveSector(drive, 0, firstSector)

	if file_system == FAT_TYPE_EXFAT {
		sectorSize := uint64(1) << firstSector[108]
		err = WriteDriveSector(
			drive,
			int64(sectorSize)*12,
			firstSector,
		) // Backup Boot Sector
		if err != nil {
			return
		}
		// Change checksum and backups
		err = fillExFATAdditionalSectors(drive, sectorSize)
		if err != nil {
			return
		}
	}

	return
}
