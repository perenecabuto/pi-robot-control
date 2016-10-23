package device

import (
	"fmt"
	"log"
	"os"

	"github.com/stianeikeland/go-rpio"
)

type Robot struct {
	Move *MotorController
}

type MotorController struct {
	pin1F rpio.Pin
	pin2F rpio.Pin
	pin1B rpio.Pin
	pin2B rpio.Pin
}

func NewMotorController(gpio1F, gpio2F, gpio1B, gpio2B uint8) *MotorController {
	pin1F := rpio.Pin(gpio1F)
	pin2F := rpio.Pin(gpio2F)
	pin1B := rpio.Pin(gpio1B)
	pin2B := rpio.Pin(gpio2B)

	return &MotorController{pin1F, pin2F, pin1B, pin2B}
}

func (c MotorController) Initialize() error {
	err := rpio.Open()
	if err != nil {
		return err
	}

	c.pin1F.Output()
	c.pin2F.Output()
	c.pin1B.Output()
	c.pin2B.Output()
	c.Stop()

	return nil
}

func (c MotorController) Forward() {
	c.move(true, true, false, false)
}

func (c MotorController) Backward() {
	c.move(false, false, true, true)
}

func (c MotorController) Left() {
	c.move(true, false, false, true)
}

func (c MotorController) Right() {
	c.move(false, true, true, false)
}

func (c MotorController) Stop() {
	c.move(false, false, false, false)
}

func (c MotorController) move(val1F, val2F, val1B, val2B bool) {
	toggle(c.pin1F, val1F)
	toggle(c.pin2F, val2F)
	toggle(c.pin1B, val1B)
	toggle(c.pin2B, val2B)
}

func toggle(pin rpio.Pin, on bool) {
	if on {
		pin.High()
	} else {
		pin.Low()
	}
}
