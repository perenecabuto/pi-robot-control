package main

/*
#cgo LDFLAGS: -ljpeg
#include "jpeg_compress.h"
*/
import "C"
import "unsafe"

func CompressImageToJpeg(frame []byte, width uint32, height uint32) ([]byte, error) {
	jpeg := make([]byte, len(frame), len(frame))
	C.compressYUYVtoJPEG(
		(*C.uint8_t)(unsafe.Pointer(&frame[0])), (*C.uint8_t)(unsafe.Pointer(&jpeg[0])),
		(C.uint32_t)(width), (C.uint32_t)(height),
	)
	return jpeg, nil
}
