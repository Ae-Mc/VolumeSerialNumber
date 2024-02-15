package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	FAT_TYPE_UNKNOWN = iota
	FAT_TYPE_FAT12   = iota
	FAT_TYPE_FAT16   = iota
	FAT_TYPE_FAT32   = iota
	FAT_TYPE_EXFAT   = iota
)

func getFATType(firstDiskSector [512]byte) int {
	reader := bytes.NewReader(firstDiskSector[:])
	var BS_JmpBoot [3]byte
	var BPB_BytsPerSec uint16
	var BPB_SecPerClus uint8
	var BPB_ResvdSecCnt uint16
	var BPB_NumFATs uint8
	var BPB_RootEntCnt uint16
	var BPB_FATSz16 uint16
	var BPB_TotSec16 uint16
	var BPB_Media uint8
	var BPB_TotSec32 uint32
	var BPB_BootSign uint16

	binary.Read(reader, binary.LittleEndian, BS_JmpBoot[:])
	reader.Seek(510, io.SeekStart)
	binary.Read(reader, binary.LittleEndian, &BPB_BootSign)

	if bytes.Equal(BS_JmpBoot[:], []byte{0xEB, 0x76, 0x90}) &&
		BPB_BootSign == 0xAA55 {
		var fileSystemName [8]byte
		reader.Seek(3, io.SeekStart)
		binary.Read(reader, binary.LittleEndian, fileSystemName[:])
		var mustBeZero, zero [53]byte
		binary.Read(reader, binary.LittleEndian, mustBeZero)
		if bytes.Equal(
			fileSystemName[:],
			[]byte{'E', 'X', 'F', 'A', 'T', ' ', ' ', ' '},
		) &&
			bytes.Equal(mustBeZero[:], zero[:]) {
			return FAT_TYPE_EXFAT
		}
	}

	reader.Seek(11, io.SeekStart)
	binary.Read(reader, binary.LittleEndian, &BPB_BytsPerSec)
	binary.Read(reader, binary.LittleEndian, &BPB_SecPerClus)
	binary.Read(reader, binary.LittleEndian, &BPB_ResvdSecCnt)
	binary.Read(reader, binary.LittleEndian, &BPB_NumFATs)
	binary.Read(reader, binary.LittleEndian, &BPB_RootEntCnt)
	binary.Read(reader, binary.LittleEndian, &BPB_TotSec16)
	binary.Read(reader, binary.LittleEndian, &BPB_Media)
	binary.Read(reader, binary.LittleEndian, &BPB_FATSz16)
	reader.Seek(32, io.SeekStart)
	binary.Read(reader, binary.LittleEndian, &BPB_TotSec32)

	if (BS_JmpBoot[0] != 0xE9 &&
		!(BS_JmpBoot[0] == 0xEB && BS_JmpBoot[2] == 0x90)) ||
		BPB_NumFATs < 1 ||
		BPB_BootSign != 0xAA55 ||
		BPB_SecPerClus == 0 ||
		(BPB_FATSz16 == 0 && BPB_RootEntCnt != 0) ||
		bytes.IndexByte(
			[]byte{0xF0, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF},
			BPB_Media,
		) == -1 ||
		((BPB_TotSec16 == 0) == (BPB_TotSec32 == 0)) {
		return FAT_TYPE_UNKNOWN
	}

	var TotSec uint32 = uint32(BPB_TotSec16)
	if TotSec == 0 {
		TotSec = BPB_TotSec32
	}
	var FATSz uint32
	if BPB_FATSz16 == 0 {
		// Probably FAT32
		reader.Seek(36, io.SeekStart)
		binary.Read(reader, binary.LittleEndian, &FATSz)
	} else {
		// Probably FAT12 or FAT16
		FATSz = uint32(BPB_FATSz16)
	}

	FatStartSector := uint32(BPB_ResvdSecCnt)
	FatSectors := FATSz * uint32(BPB_NumFATs)
	RootDirStartSector := FatStartSector + FatSectors
	RootDirSectors := (32*uint32(BPB_RootEntCnt) + uint32(BPB_BytsPerSec) - 1) / uint32(
		BPB_BytsPerSec,
	)
	DataStartSector := RootDirStartSector + RootDirSectors
	DataSectors := TotSec - DataStartSector
	CountofClusters := DataSectors / uint32(BPB_SecPerClus)

	switch {
	case CountofClusters < 4085:
		return FAT_TYPE_FAT12
	case CountofClusters < 65525:
		return FAT_TYPE_FAT16
	default:
		return FAT_TYPE_FAT32
	}
}

func isFAT12(firstDiskSector [512]byte) bool {
	return getFATType(firstDiskSector) == FAT_TYPE_FAT12
}
func isFAT16(firstDiskSector [512]byte) bool {
	return getFATType(firstDiskSector) == FAT_TYPE_FAT16
}
func isFAT32(firstDiskSector [512]byte) bool {
	return getFATType(firstDiskSector) == FAT_TYPE_FAT32
}
func isEXFAT(firstDiskSector [512]byte) bool {
	return getFATType(firstDiskSector) == FAT_TYPE_EXFAT
}
