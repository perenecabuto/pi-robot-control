package main

import (
	"errors"
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
	cam     *webcam.Webcam
}

func NewWebcamCapture(timeout uint32, address string) *WebcamCapture {
	return &WebcamCapture{timeout, address, nil}
}

func (w *WebcamCapture) Initialize() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	cam, err := webcam.Open(w.address) // Open webcam
	if err != nil {
		log.Panic(err.Error())
	}
	err = setupCamImageFormat(cam)
	if err != nil {
		log.Panic(err.Error())
	}
	err = cam.StartStreaming()
	if err != nil {
		log.Panic(err.Error())
	}
	w.cam = cam
	return nil
}

func (w WebcamCapture) Listen(onFrame func([]byte)) error {
	if w.cam == nil {
		if err := w.Initialize(); err != nil {
			return err
		}
	}
	defer w.cam.Close()

	for {
		err := w.cam.WaitForFrame(w.timeout)
		if err != nil {
			switch err.(type) {
			case *webcam.Timeout:
				continue
			default:
				return err
			}
		}

		frame, err := w.cam.ReadFrame()
		if len(frame) != 0 {
			onFrame(frame)
		} else if err != nil {
			return err
		}
	}
}

func setupCamImageFormat(cam *webcam.Webcam) error {
	var format webcam.PixelFormat
	log.Println("Supported formats:", cam.GetSupportedFormats())

	for f, name := range cam.GetSupportedFormats() {
		if strings.Contains(name, "JPEG") {
			log.Println("Camera JPEG format found:", name)
			format = f
			break
		}
	}
	if format == 0 {
		return errors.New("No format found")
	}

	log.Println("Camera dimensions:", cam.GetSupportedFrameSizes(format))

	_, _, _, err := cam.SetImageFormat(format, IMAGE_WIDTH, IMAGE_HEIGHT)
	return err
}
