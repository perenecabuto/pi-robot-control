package handler

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type MJPEGStream struct {
	frame    []byte
	fps      int
	boundary string
	lock     sync.Mutex
}

func NewMJPEGStream(fps int) *MJPEGStream {
	return &MJPEGStream{
		frame:    nil,
		fps:      fps,
		boundary: strconv.Itoa(rand.Int()),
	}
}

func (s *MJPEGStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Stream:", r.RemoteAddr, "connected")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+s.boundary)

	c := time.Tick(time.Second / time.Duration(s.fps))
	for range c {
		_, err := w.Write(s.frame)
		if err != nil {
			break
		}
	}

	log.Println("Stream:", r.RemoteAddr, "disconnected")
}

func (s *MJPEGStream) UpdateJPEG(jpeg []byte) {
	header := s.buildHeader(len(jpeg))
	if s.frame == nil || len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}

	s.lock.Lock()
	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)
	s.lock.Unlock()
}

func (s MJPEGStream) buildHeader(length int) string {
	return "\r\n" +
		"--" + s.boundary + "\r\n" +
		"Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(length) + "\r\n" +
		"X-Timestamp: 0.000000\r\n" +
		"\r\n"
}
