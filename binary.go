package dlt645

import (
	"fmt"
	"strconv"
)

func Hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}
func Byte2Hex(bs []byte) string {
	s := ""
	for _, b := range bs {
		s += fmt.Sprintf("%02x", b)
	}
	return s
}

func ByteSub(bs []byte, sub byte) []byte {
	r := []byte{}
	for _, b := range bs {
		r = append(r, b-sub)
	}
	return r
}
func ByteAdd(bs []byte, add byte) []byte {
	r := []byte{}
	for _, b := range bs {
		r = append(r, b+add)
	}
	return r
}

// 反转 byte
func ByteRev(bs []byte) []byte {
	r := make([]byte, len(bs))
	for i, b := range bs {
		r[len(bs)-i-1] = b
	}
	return r
}
