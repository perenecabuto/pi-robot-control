package main

import "github.com/stianeikeland/go-rpio"

type Robot struct {
	pin1F rpio.Pin
	pin2F rpio.Pin
	pin1B rpio.Pin
	pin2B rpio.Pin
}

func NewRobot(gpio1F, gpio2F, gpio1B, gpio2B uint8) *Robot {
	pin1F := rpio.Pin(gpio1F)
	pin2F := rpio.Pin(gpio2F)
	pin1B := rpio.Pin(gpio1B)
	pin2B := rpio.Pin(gpio2B)

	return &Robot{pin1F, pin2F, pin1B, pin2B}
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
