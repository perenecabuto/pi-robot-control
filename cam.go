package main

import (
	"log"
	"strings"

	"github.com/blackjack/webcam"
)

const (
	IMAGE_WIDTH  = 480
	IMAGE_HEIGHT = 480
)

type WebcamCapture struct {
	timeout uint32
	address string
}

func NewWebcamCapture(timeout uint32, address string) *WebcamCapture {
	return &WebcamCapture{timeout, address}
}

func (w WebcamCapture) Listen(onFrame func([]byte)) {
	cam, err := webcam.Open(w.address) // Open webcam
	if err != nil {
		log.Panic(err.Error())
	}
	defer cam.Close()

	format, err := setupCamImageFormat(cam)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("Chosen format:", format)

	err = cam.StartStreaming()
	if err != nil {
		log.Panic(err.Error())
	}

	for {
		err = cam.WaitForFrame(w.timeout)
		if err != nil {
			switch err.(type) {
			case *webcam.Timeout:
				continue
			default:
				log.Panic(err.Error())
			}
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Panic(err.Error())
		}
		if len(frame) != 0 {
			var jpegFrame []byte
			if format.IsJPEG() {
				jpegFrame = frame
			} else {
				jpegFrame, err = CompressImageToJpeg(frame)
				if err != nil {
					log.Panic(err.Error())
				}
			}
			onFrame(jpegFrame)
		}
	}
}

type PixelFormatName string

func (n PixelFormatName) IsJPEG() bool {
	return strings.Contains(string(n), "JPEG")
}

func setupCamImageFormat(cam *webcam.Webcam) (PixelFormatName, error) {
	log.Println("Supported formats:", cam.GetSupportedFormats())

	var format webcam.PixelFormat
	var name PixelFormatName
	for f, n := range cam.GetSupportedFormats() {
		format, name = f, PixelFormatName(n)
		if name.IsJPEG() {
			log.Println("Camera JPEG format found:", name)
			break
		}
	}

	log.Println("Camera dimensions:", cam.GetSupportedFrameSizes(format))
	_, _, _, err := cam.SetImageFormat(format, IMAGE_WIDTH, IMAGE_HEIGHT)
	return name, err
}
