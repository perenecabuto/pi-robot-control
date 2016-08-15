package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stianeikeland/go-rpio"
)

var (
	ServerAddress = "0.0.0.0:8000"
)

func main() {
	robot, err := NewRobot(17, 27, 4, 22)
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/move/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				http.Error(w, "Unexpcted error", 500)
			}
		}()

		direction := r.URL.Query().Get("direction")
		switch direction {
		case "right":
			robot.Right()
		case "left":
			robot.Left()
		case "forward":
			robot.Forward()
		case "backward":
			robot.Backward()
		case "stop":
			robot.Stop()
		default:
			http.Error(w, "action "+direction+" is not allowed", 500)
			return
		}

		log.Println("Got move action:", direction)
	})

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

	log.Println("Starting server on:", ServerAddress)
	log.Panic(http.ListenAndServe(ServerAddress, nil))
}

type Robot struct {
	pin1F rpio.Pin
	pin2F rpio.Pin
	pin1B rpio.Pin
	pin2B rpio.Pin
}

func NewRobot(gpio1F, gpio2F, gpio1B, gpio2B uint8) (*Robot, error) {
	pin1F := rpio.Pin(gpio1F)
	pin2F := rpio.Pin(gpio2F)
	pin1B := rpio.Pin(gpio1B)
	pin2B := rpio.Pin(gpio2B)

	return &Robot{pin1F, pin2F, pin1B, pin2B}, nil
}

func (r Robot) Initialize() error {
	err := rpio.Open()
	if err != nil {
		return err
	}

	r.pin1F.Output()
	r.pin2F.Output()
	r.pin1B.Output()
	r.pin2B.Output()

	r.Stop()
	return nil
}

func (r Robot) Forward() {
	r.Stop()
	r.pin1F.High()
	r.pin2F.High()
}

func (r Robot) Backward() {
	r.Stop()
	r.pin1B.High()
	r.pin2B.High()
}

func (r Robot) Left() {
	r.Stop()
	r.pin1F.High()
	r.pin2B.High()
}

func (r Robot) Right() {
	r.Stop()
	r.pin1B.High()
	r.pin2F.High()
}
func (r Robot) Stop() {
	r.pin1F.Low()
	r.pin2F.Low()
	r.pin1B.Low()
	r.pin2B.Low()
}
