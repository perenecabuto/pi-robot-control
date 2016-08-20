package main

/*
#cgo LDFLAGS: -ljpeg
#include "jpeg_compress.h"
*/
import "C"
import "unsafe"

func CompressImageToJpeg(frame []byte) ([]byte, error) {
	jpeg := make([]byte, len(frame), len(frame))
	C.compressYUYVtoJPEG((*C.char)(unsafe.Pointer(&frame[0])), (*C.char)(unsafe.Pointer(&jpeg[0])), 480, 480)
	return jpeg, nil
}
