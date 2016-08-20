package main

import (
	"log"
	"net/http"

	"github.com/saljam/mjpeg"
)

var (
	ServerAddress = "0.0.0.0:8000"
)

func main() {
	robot := NewRobot(17, 27, 4, 22)
	robotHandler := NewRobotHandler(robot)
	stream := mjpeg.NewStream()
	capture := NewWebcamCapture(1000, "/dev/video0")

	http.Handle("/camera", stream)
	http.Handle("/move/", robotHandler)
	http.HandleFunc("/", IndexHandler)

	go func() {
		err := capture.Listen(func(frame []byte) {
			stream.UpdateJPEG(frame)
		})
		if err != nil {
			log.Println(err)
		}
	}()

	go robot.Initialize()
	log.Println("Starting server on:", ServerAddress)
	log.Panic(http.ListenAndServe(ServerAddress, nil))
}
