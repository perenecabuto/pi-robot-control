package device

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/blackjack/webcam"
)

type frameSize []uint32

// WebCam controller for webcam device
type WebCam struct {
	address      string
	maxFrameSize frameSize
}

// NewWebCam creates a WebCam on device address
func NewWebCam(address string, maxFrameSize frameSize) *WebCam {
	return &WebCam{address, maxFrameSize}
}

// Listen to webcam frames
func (w *WebCam) Listen(fps int, onFrame func([]byte)) error {
	cam, err := webcam.Open(w.address) // Open webcam
	if err != nil {
		return err
	}
	defer cam.Close()

	format, err := setupCamImageFormat(cam, w.maxFrameSize[0], w.maxFrameSize[1])
	if err != nil {
		return err
	}

	log.Println("Starting cam input stream with format:", format)
	if err := cam.StartStreaming(); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	var jpegFrame []byte

	for {
		select {
		case <-ticker.C:
			go onFrame(jpegFrame)
		default:
			err := cam.WaitForFrame(1)
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
				if format.IsJPEG() {
					jpegFrame = frame
				} else {
					jpegFrame, err = CompressImageToJpeg(frame, format.Width, format.Height)
					if err != nil {
						log.Println("Skipping frame - reason:", err)
						return err
					}
				}
			}
		}
	}
}

type CamPixelFormat struct {
	Name   string
	Width  uint32
	Height uint32
}

func (f *CamPixelFormat) IsJPEG() bool {
	return strings.Contains(f.Name, "JPEG")
}

func (f *CamPixelFormat) String() string {
	return fmt.Sprintf("format: %s, dimensions: %dx%d", f.Name, f.Width, f.Height)
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

func setupCamImageFormat(cam *webcam.Webcam, maxWidth, maxHeight uint32) (*CamPixelFormat, error) {
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
	sort.Sort(supportedSizes)
	log.Println("Supported sizes:", supportedSizes)

	size := supportedSizes[0]

	found.Width, found.Height = size.MaxWidth, size.MaxHeight
	if found.Width > maxWidth {
		found.Width = maxWidth
	}
	if found.Height > maxHeight {
		found.Height = maxHeight
	}
	_, _, _, err := cam.SetImageFormat(format, found.Width, found.Height)

	return found, err
}
