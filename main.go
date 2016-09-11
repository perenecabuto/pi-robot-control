package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./device"
	"./handler"
)

var (
	CameraDevice  = flag.String("d", "/dev/video0", "Video dev path")
	ServerAddress = flag.String("a", "0.0.0.0:8000", "Server address")
	FrameTimeout  = flag.Int("frametimeout", 1000, "Frame timeout")
	FPS           = flag.Int("fps", 5, "Frames per second")
)

func main() {
	flag.Parse()
	robot := device.NewRobot(17, 27, 4, 22)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL|syscall.SIGTERM)
	go func() {
		<-c
		robot.CleanUP()
		os.Exit(1)
	}()

	if err := robot.Initialize(); err != nil {
		log.Println("Error - Robot init", err.Error())
	}

	robotHandler := handler.NewRobotHandler(robot)
	http.Handle("/control/", robotHandler.ListenWS())
	http.HandleFunc("/", handler.IndexHandler)

	cam := device.NewWebCam(uint32(*FrameTimeout), *CameraDevice)
	stream := handler.NewMJPEGStream(*FPS)
	endpointOpened := false
	go cam.Listen(*FPS, func(frame *[]byte) {
		if !endpointOpened {
			log.Println("Open camera endpoint")
			http.Handle("/camera", stream)
			endpointOpened = true
		}
		stream.UpdateJPEG(frame)
	})

	log.Println("Starting server on:", *ServerAddress)
	log.Panic(http.ListenAndServe(*ServerAddress, nil))
}
