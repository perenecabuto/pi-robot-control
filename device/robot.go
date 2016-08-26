package device

import (
	"github.com/stianeikeland/go-rpio"
)

type Robot struct {
	pin1F     rpio.Pin
	pin2F     rpio.Pin
	pin1B     rpio.Pin
	pin2B     rpio.Pin
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

	go func() {
		for range r.stopTimer.C {
			r.stopTimer.Stop()
			r.Stop()
		}
	}()

	return nil
}

func (r Robot) Forward() {
	r.move(true, true, false, false)
}

func (r Robot) Backward() {
	r.move(false, false, true, true)
}

func (r Robot) Left() {
	r.move(true, false, false, true)
}

func (r Robot) Right() {
	r.move(false, true, true, false)
}

func (r Robot) Stop() {
	r.move(false, false, false, false)
}

func (r Robot) move(val1F, val2F, val1B, val2B bool) {
	toggle(r.pin1F, val1F)
	toggle(r.pin2F, val2F)
	toggle(r.pin1B, val1B)
	toggle(r.pin2B, val2B)
}

func toggle(pin rpio.Pin, on bool) {
	if on {
		pin.High()
	} else {
		pin.Low()
	}
}
