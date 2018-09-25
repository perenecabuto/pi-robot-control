package handler

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// MJPEGStream streams jpeg images
type MJPEGStream struct {
	frame    []byte
	fps      int
	boundary string
	lock     sync.Mutex
}

// NewMJPEGStream create a MJPEGStream with FPS
func NewMJPEGStream(fps int) *MJPEGStream {
	return &MJPEGStream{
		frame:    nil,
		fps:      fps,
		boundary: strconv.Itoa(rand.Int()),
	}
}

func (s *MJPEGStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+s.boundary)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	fps, err := strconv.ParseInt(r.URL.Query().Get("fps"), 10, 32)
	if err != nil {
		fps = int64(s.fps)
	}

	log.Println("Stream:", r.RemoteAddr, "connected - video FPS:", fps)
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, err := w.Write(s.frame)
			if err != nil {
				break
			}
		case <-r.Context().Done():
			log.Println("Stream:", r.RemoteAddr, "disconnected")
			return
		}
	}
}

// UpdateJPEG set the current jpeg frame
func (s *MJPEGStream) UpdateJPEG(jpeg []byte) {
	if jpeg == nil {
		return
	}
	jpegLen := len(jpeg)
	header := s.buildHeader(jpegLen)
	if s.frame == nil || len(s.frame) < jpegLen+len(header) {
		s.frame = make([]byte, (jpegLen+len(header))*2)
	}

	s.lock.Lock()
	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)
	s.lock.Unlock()
}

func (s *MJPEGStream) buildHeader(length int) string {
	return "\r\n" +
		"--" + s.boundary + "\r\n" +
		"Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(length) + "\r\n" +
		"X-Timestamp: 0.000000\r\n" +
		"\r\n"
}
