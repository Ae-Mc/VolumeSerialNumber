package main

import "bytes"

func isNTFS(firstDiskSector []byte) bool {
	return bytes.Equal(
		firstDiskSector[3:11],
		[]byte{'N', 'T', 'F', 'S', ' ', ' ', ' ', ' '},
	)
}
