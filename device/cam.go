package device

import (
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/blackjack/webcam"
)

type WebcamCapture struct {
	timeout uint32
	address string
}

func NewWebcamCapture(timeout uint32, address string) *WebcamCapture {
	return &WebcamCapture{timeout, address}
}

func (w WebcamCapture) Listen(onFrame func([]byte)) error {
	cam, err := webcam.Open(w.address) // Open webcam
	if err != nil {
		return err
	}
	defer cam.Close()

	format, err := setupCamImageFormat(cam)
	if err != nil {
		return err
	}

	log.Println("Chosen format:", format)

	if err := cam.StartStreaming(); err != nil {
		return err
	}
	for {
		err := cam.WaitForFrame(w.timeout)
		if err != nil {
			switch err.(type) {
			case *webcam.Timeout:
				continue
			default:
				return err
			}
		}

		frame, err := cam.ReadFrame()
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

type SortedSizes []webcam.FrameSize

func (s SortedSizes) Len() int {
	return len(s)
}

func (s SortedSizes) Less(i, j int) bool {
	return s[i].MaxWidth < s[j].MaxWidth
}

func (s SortedSizes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
	if format == 0 {
		return nil, errors.New("No pixel format found")
	}

	supportedSizes := SortedSizes(cam.GetSupportedFrameSizes(format))
	log.Println("Supported sizes:", supportedSizes)
	sort.Sort(supportedSizes)
	size := supportedSizes[0]

	found.Width, found.Height = size.MaxWidth, size.MaxHeight
	_, _, _, err := cam.SetImageFormat(format, found.Width, found.Height)

	return found, err
}
