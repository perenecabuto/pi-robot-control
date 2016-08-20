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
	robot, err := NewRobot(17, 27, 4, 22)
	if err != nil {
		log.Panic(err)
	}

	robotHandler := NewRobotHandler(robot)
	stream := mjpeg.NewStream()
	capture := NewWebcamCapture(1000, "/dev/video0")

	http.Handle("/camera", stream)
	http.Handle("/move", robotHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := ioutil.ReadFile("webapp/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(file)
	})

	go robot.Initialize()
	go capture.Listen(func(frame []byte) {
		stream.UpdateJPEG(frame)
	})

	log.Println("Starting server on:", ServerAddress)
	log.Panic(http.ListenAndServe(ServerAddress, nil))
}
