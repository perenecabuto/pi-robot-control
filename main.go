package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/perenecabuto/robot-control/device"
	"github.com/perenecabuto/robot-control/handler"
)

var (
	cameraDevice   = flag.String("d", "/dev/video0", "Video dev path")
	serverAddress  = flag.String("a", "0.0.0.0:8000", "Server address")
	maxFrameWidth  = flag.Uint("frame-width", 320, "frame max width")
	maxFrameHeight = flag.Uint("frame-height", 240, "frame max height")
	fps            = flag.Int("fps", 30, "Frames per second")
	wheelPins      = flag.String("pins", "25,27,17,22", "Wheel gpios as int separated with by comma."+
		"The order is : <left-forward>,<right-forward>,<left back>,<right back>")
)

func main() {
	flag.Parse()

	pins := make([]uint8, 4, 4)
	for i, pin := range strings.SplitN(*wheelPins, ",", 4) {
		if ipin, err := strconv.Atoi(pin); err == nil {
			pins[i] = uint8(ipin)
		} else {
			log.Fatal("ERROR - Pins must be a list of four ints separated by comma: " + err.Error())
		}
	}
	log.Println("WheelPins", pins)

	motorC := device.NewMotorController(pins[0], pins[1], pins[2], pins[3])
	camPositionC := device.NewCamPositionController(24, 23)
	robot := device.NewRobot(motorC, camPositionC)
	if err := robot.Initialize(); err != nil {
		log.Println(err.Error())
	}

	robotH := handler.NewRobotHandler(robot)
	http.Handle("/control/", robotH.ListenWS())
	http.Handle("/robot/look-to/", robotH.LookTo())

	compass := device.NewCompass(1, 0x1e, 1.3)
	compassH := handler.NewCompassHandler(compass)

	cam := device.NewWebCam(*cameraDevice, []uint32{uint32(*maxFrameWidth), uint32(*maxFrameHeight)})
	stream := handler.NewMJPEGStream(*fps)
	go func() {
		err := cam.Listen(*fps, stream.UpdateJPEG)
		if err != nil {
			log.Fatal("Could not start stream:", err.Error())
		}
	}()
	http.Handle("/camera", stream)
	http.Handle("/compass/", compassH)
	http.HandleFunc("/", handler.IndexHandler)

	robot.Look.To(90, 90)
	log.Println("Starting server on:", *serverAddress)
	log.Fatal(http.ListenAndServe(*serverAddress, nil))
}
