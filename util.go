package main

import (
	"fmt"
)

func toHex(data []byte) (repr []string) {
	for _, b := range data {
		repr = append(repr, fmt.Sprintf("0x%X", b))
	}
	return repr
}

func filter[Type any](arr []Type, test func(item Type) bool) (result []Type) {
	for _, item := range arr {
		if test(item) {
			result = append(result, item)
		}
	}
	return
}

func noErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
