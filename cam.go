package main

import (
	"errors"
	"log"
	"strings"

	"github.com/blackjack/webcam"
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
		if len(frame) != 0 {
			var jpegFrame []byte
			if format.IsJPEG() {
				jpegFrame = frame
			} else {
				jpegFrame, err = CompressImageToJpeg(frame, format.Width, format.Height)
				if err != nil {
					return err
				}
			}
			onFrame(jpegFrame)
		}
	}
}

type CamPixelFormat struct {
	Name   string
	Width  uint32
	Height uint32
}

func (f CamPixelFormat) IsJPEG() bool {
	return strings.Contains(f.Name, "JPEG")
}

func setupCamImageFormat(cam *webcam.Webcam) (*CamPixelFormat, error) {
	log.Println("Supported formats:", cam.GetSupportedFormats())

	var format webcam.PixelFormat
	var found *CamPixelFormat
	for f, n := range cam.GetSupportedFormats() {
		format, found = f, &CamPixelFormat{Name: n}
		if found.IsJPEG() {
			log.Println("Camera JPEG format found:", found)
			break
		}
	}

	supportedSizes := cam.GetSupportedFrameSizes(format)
	size := supportedSizes[0]

	found.Width, found.Height = size.MaxWidth, size.MaxHeight
	_, _, _, err := cam.SetImageFormat(format, found.Width, found.Height)

	return found, err
}
