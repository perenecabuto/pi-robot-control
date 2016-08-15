package main

import (
	"errors"
	"log"
	"reflect"
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

	setupCamImageFormat(cam)
	if err != nil {
		log.Panic(err.Error())
	}

	err = cam.StartStreaming()
	if err != nil {
		log.Panic(err.Error())
	}

	for {
		err = cam.WaitForFrame(w.timeout)
		if err != nil {
			if reflect.TypeOf(err).Name() == "*webcam.Timeout" {
				continue
			}
			log.Panic(err.Error())
		}

		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			onFrame(frame)
		} else if err != nil {
			log.Panic(err.Error())
		}
	}
}

func setupCamImageFormat(cam *webcam.Webcam) error {
	var format webcam.PixelFormat
	for f, name := range cam.GetSupportedFormats() {
		if strings.Contains(name, "JPEG") {
			format = f
			break
		}
	}
	if format == 0 {
		return errors.New("No format found")
	}

	_, _, _, err := cam.SetImageFormat(format, IMAGE_WIDTH, IMAGE_HEIGHT)
	return err
}
