package wchar

/*
#include <stdlib.h>
#include <wchar.h>
*/
import "C"

import (
	"encoding/binary"
	"unsafe"
)

const Wsize = 2

type Wchar uint8

type WcharString []Wchar

func (out *Wchar) FromStr(s string) {
	ConvertGoStringToWcharString(s, out)
}

func ToStr(in *Wchar) string {
	return ConvertWcharStringToGoString(in)
}

func ConvertGoStringToWchar(input string, out *Wchar) {
	if input == "" {
		zs := make(WcharString, 0)
		out = &zs[0]
	}

	outLen := len(input) * Wsize

	ret := make(WcharString, 0, outLen)

	for _, char := range input {
		ret = append(ret, Wchar(char), Wchar(0))
	}

	C.memcpy(unsafe.Pointer(out), unsafe.Pointer(&ret[0]), C.size_t(outLen))
}

func ConvertWcharToGoString(in *Wchar) string {
	out := ""

	wcharPtr := uintptr(unsafe.Pointer(in))
	for {
		w := *((*Wchar)(unsafe.Pointer(wcharPtr)))
		if w == 0 {
			break
		}

		ws = append(ws, w)
		out += C.GoString(C.char(ws))
		wcharPtr += Wsize
	}

	return out
}
