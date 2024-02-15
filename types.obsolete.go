package main

import (
	"fmt"
	"unsafe"
)

// Partition styles
const (
	PARTITION_STYLE_MBR = iota
	PARTITION_STYLE_GPT = iota
	PARTITION_STYLE_RAW = iota
)

// GPT Partition types
var (
	PARTITION_BASIC_DATA_GUID = GUID{
		0xebd0a0a2, 0xb9e5, 0x4433, [8]byte{0x68, 0xb6, 0xb7, 0x26, 0x99, 0xc7},
	}
	PARTITION_ENTRY_UNUSED_GUID = GUID{}
	PARTITION_SYSTEM_GUID       = GUID{
		0xc12a7328, 0xf81f, 0x11d2, [8]byte{0xba, 0x4b, 0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b},
	}
	PARTITION_MSFT_RESERVED_GUID = GUID{
		0xe3c9e316, 0x0b5c, 0x4db8, [8]byte{0x81, 0x7d, 0xf9, 0x2d, 0xf0, 0x02, 0x15, 0xae},
	}
	PARTITION_LDM_METADATA_GUID = GUID{
		0x5808c8aa, 0x7e8f, 0x42e0, [8]byte{0x85, 0xd2, 0xe1, 0xe9, 0x04, 0x34, 0xcf, 0xb3},
	}
	PARTITION_LDM_DATA_GUID = GUID{
		0xaf9b60a0, 0x1431, 0x4f62, [8]byte{0xbc, 0x68, 0x33, 0x11, 0x71, 0x4a, 0x69, 0xad},
	}
	PARTITION_MSFT_RECOVERY_GUID = GUID{
		0xde94bba4, 0x06d1, 0x4d40, [8]byte{0xa1, 0x6a, 0xbf, 0xd5, 0x01, 0x79, 0xd6, 0xac},
	}
)

// MBR Partition types
const (
	PARTITION_ENTRY_UNUSED = 0x00 // Неиспользуемая секция записи.
	PARTITION_EXTENDED     = 0x05 // Расширенная секция.
	PARTITION_FAT_12       = 0x01 // Раздел файловой системы FAT12.
	PARTITION_FAT_16       = 0x04 // Раздел файловой системы FAT16.
	PARTITION_FAT32        = 0x0B // Раздел файловой системы FAT32.
	PARTITION_IFS          = 0x07 // Секция IFS.
	PARTITION_LDM          = 0x42 // Раздел диспетчера логических дисков (LDM).
	PARTITION_NTFT         = 0x80 // Раздел NTFT.
	VALID_NTFT             = 0xC0 // Допустимый раздел NTFT.
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

func (guid GUID) String() string {
	result := fmt.Sprintf(
		"%04x-%02x-%02x-%x%x-%x%x%x%x%x%x",
		guid.Data1,
		guid.Data2,
		guid.Data3,
		guid.Data4[0],
		guid.Data4[1],
		guid.Data4[2],
		guid.Data4[3],
		guid.Data4[4],
		guid.Data4[5],
		guid.Data4[6],
		guid.Data4[7],
	)
	return result
}

type PARTITION_INFORMATION_MBR struct {
	PartitionType       byte
	BootIndicator       int32
	RecognizedPartition int32
	HiddenSectors       uint32
	PartitionId         GUID
}

type PARTITION_INFORMATION_GPT struct {
	PartitionType GUID
	PartitionId   GUID
	Attributes    uint64
	Name          [36]uint16
}

func (dummy PARTITION_INFORMATION_RAW) GPT() PARTITION_INFORMATION_GPT {
	return *(*PARTITION_INFORMATION_GPT)(unsafe.Pointer(&dummy))
}

type PARTITION_INFORMATION_RAW [unsafe.Sizeof(PARTITION_INFORMATION_GPT{})]byte

type PARTITION_INFORMATION_EX struct {
	PartitionStyle   uint32
	StartingOffset   uint64
	PartitionLength  uint64
	PartitionNumber  uint32
	RewritePartition int32 // bool
	// RewritePartition int16 // bool
	// IsServicePartition int16                     // bool
	DUMMYUNIONNAME PARTITION_INFORMATION_RAW // Can be GPT or MBR
}
