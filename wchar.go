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

func NewWcharString(length int) WcharString {
	return make(WcharString, length)
}

func (out *Wchar) FromStr(s string) {
	ConvertGoStringToWcharString(s, out)
}

func ToStr(in *Wchar) string {
	return ConvertWcharStringToGoString(in)
}

func FromWcharStringPtr() WcharString {
	if uintptr(first) == 0x0 {
		return NewWcharString(0)
	}

	wcharPtr := uintptr(first)

	ws := make(WcharString, 0)

	var w Wchar
	for {
		w = *((*Wchar)(unsafe.Pointer(wcharPtr)))
		if w == 0 {
			break
		}

		ws = append(ws, w)
		wcharPtr += Wsize
	}

	return ws
}

func ConvertGoStringToWcharString(input string, out *Wchar) {
	if input == "" {
		zs := NewWcharString(0)
		out = &zs[0]
	}

	outLen := len(input) * Wsize

	ret := make(WcharString, 0, outLen)

	for _, char := range input {
		ret = append(ret, Wchar(char), Wchar(0))
	}

	C.memcpy(unsafe.Pointer(out), unsafe.Pointer(&ret[0]), C.size_t(outLen))
}

func ConvertWcharStringToGoString(in *Wchar) (output string) {
	first := unsafe.Pointer(in)
	if uintptr(first) == 0x0 {
		return ""
	}

	wcharPtr := uintptr(first)
	ws := make(WcharString, 0)

	var w Wchar
	for {
		w = *((*Wchar)(unsafe.Pointer(wcharPtr)))
		if w == 0 {
			break
		}

		ws = append(ws, w)
		wcharPtr += Wsize
	}

	inputAsCChars := make([]C.char, 0, len(ws)*4)
	wcharAsBytes := make([]byte, 4)
	for _, nextWchar := range ws {
		binary.LittleEndian.PutUint32(wcharAsBytes, uint32(nextWchar))
		for i := 0; i < 4; i++ {
			inputAsCChars = append(inputAsCChars, C.char(wcharAsBytes[i]))
		}
	}

	output = C.GoStringN((*C.char)(&inputAsCChars[0]), (C.int)(len(inputAsCChars)))

	return output
}
