package main

import (
	"flag"
	"log"
	"net/http"

	"./device"
	"./handler"
)

var (
	CameraDevice  = flag.String("d", "/dev/video0", "Video dev path")
	ServerAddress = flag.String("a", "0.0.0.0:8000", "Server address")
	FrameTimeout  = flag.Int("frametimeout", 100, "Frame timeout")
	FPS           = flag.Int("fps", 5, "Frames per second")
)

func main() {
	flag.Parse()
	robot := device.NewRobot(17, 27, 4, 22)
	if err := robot.Initialize(); err != nil {
		log.Println(err.Error())
	}

	stream := handler.NewMJPEGStream(*FPS)
	robotHandler := handler.NewRobotHandler(robot)
	http.Handle("/move/", robotHandler)
	http.Handle("/control/", robotHandler.ListenWS())
	http.Handle("/camera", stream)
	http.HandleFunc("/", handler.IndexHandler)

	capture := device.NewWebcamCapture(uint32(*FrameTimeout), *CameraDevice)
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
