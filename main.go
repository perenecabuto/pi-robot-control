package main

import (
	"log"
	"net/http"

	"github.com/saljam/mjpeg"
)

var (
	FrameTimeout  = 1000
	CameraDevice  = "/dev/video0"
	ServerAddress = "0.0.0.0:8000"
)

func main() {
	robot := NewRobot(17, 27, 4, 22)
	robotHandler := NewRobotHandler(robot)
	stream := mjpeg.NewStream()
	capture := NewWebcamCapture(uint32(FrameTimeout), CameraDevice)

	http.Handle("/move/", robotHandler)
	http.HandleFunc("/", IndexHandler)

	err := capture.Initialize()
	if err == nil {
		http.Handle("/camera", stream)
		go capture.Listen(func(frame []byte) {
			stream.UpdateJPEG(frame)
		})
	} else {
		log.Println("Error starting camera:", err)
	}

	go robot.Initialize()
	log.Println("Starting server on:", ServerAddress)
	log.Panic(http.ListenAndServe(ServerAddress, nil))
}
