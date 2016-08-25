package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/saljam/mjpeg"
)

var (
	FrameTimeout  = flag.Int("frametimeout", 100, "Frame timeout")
	CameraDevice  = flag.String("d", "/dev/video0", "Video dev path")
	ServerAddress = flag.String("a", "0.0.0.0:8000", "Server address")
)

func main() {
	flag.Parse()
	robot := NewRobot(17, 27, 4, 22)
	if err := robot.Initialize(); err != nil {
		log.Println(err.Error())
	}
	robotHandler := NewRobotHandler(robot)
	stream := mjpeg.NewStream()
	capture := NewWebcamCapture(uint32(*FrameTimeout), *CameraDevice)

	http.Handle("/move/", robotHandler)
	http.HandleFunc("/", IndexHandler)

	http.Handle("/camera", stream)
	go func() {
		err := capture.Listen(func(frame []byte) {
			stream.UpdateJPEG(frame)
		})
		if err != nil {
			log.Println("Error starting camera:", err)
		}
	}()

	log.Println("Starting server on:", *ServerAddress)
	log.Panic(http.ListenAndServe(*ServerAddress, nil))
}
