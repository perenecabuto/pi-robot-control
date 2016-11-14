package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"./device"
	"./handler"
)

var (
	CameraDevice  = flag.String("d", "/dev/video0", "Video dev path")
	ServerAddress = flag.String("a", "0.0.0.0:8000", "Server address")
	FrameTimeout  = flag.Int("frametimeout", 1000, "Frame timeout")
	FPS           = flag.Int("fps", 5, "Frames per second")
	WheelPins     = flag.String("pins", "25,27,17,22", "Wheel gpios as int separated with by comma."+
		"The order is : <left-forward>,<right-forward>,<left back>,<right back>")
)

func main() {
	flag.Parse()

	pins := make([]uint8, 4, 4)
	for i, pin := range strings.SplitN(*WheelPins, ",", 4) {
		if ipin, err := strconv.Atoi(pin); err == nil {
			pins[i] = uint8(ipin)
		} else {
			panic("ERROR - Pins must be a list of four ints separated by comma: " + err.Error())
		}
	}
	log.Println("WheelPins", pins)

	motorC := device.NewMotorController(pins[0], pins[1], pins[2], pins[3])
	camPositionC := device.NewCamPositionController(24, 23)
	robot := device.NewRobot(motorC, camPositionC)
	if err := robot.Initialize(); err != nil {
		log.Println(err.Error())
	}

	robotHandler := handler.NewRobotHandler(robot)
	http.Handle("/control/", robotHandler.ListenWS())
	http.Handle("/robot/look-to/", robotHandler.LookToHandler())
	http.HandleFunc("/compass/", handler.CompassHandler)
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
