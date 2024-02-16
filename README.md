# Volume serial number

This package provides many functions, but most important are 
`GetVolumeSerialNumber` and `SetVolumeSerialNumber`.
These functions use uint64 for storing serial number.
If serial number is 4 bytes long then it uses only 4 low bytes of uint64.
Serial number contains parsed value in native byte order.
To print it in often used format use `binary.Write` with `binary.BigEndian` byte order.
Checkout source code for more explanations.

## Supported file systems
- FAT12
- FAT16
- FAT32
- exFAT
- NTFS

## Supported OS

### Windows
Tested on XP 32 bit and Windows 10 32 and 64 bit.
> [!Warning]
> To use these functions drive name must be in format `\\\\.\\C:` (this is escaped string).

### Linux
Must be supported if you can somehow get right drive path without offset
(with access to FAT/exFAT/NTFS's Boot Sector'). No OS specific functions were
used — only `os.Read`, `os.File.Seek` and `os.File.Write` for disk
manipulations. But untested.