package device

import (
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio"
)

type Robot struct {
	Move *MotorController
	Look *CamPositionController
}

func NewRobot(motor *MotorController, camPosition *CamPositionController) *Robot {
	return &Robot{motor, camPosition}
}

func (r Robot) Initialize() error {
	return r.Move.Initialize()
}

type CamPositionController struct {
	xAxisGPIO uint8
	yAxisGPIO uint8
	MinPulse  float32
	MaxPulse  float32

	xMinAngle           uint8
	xMaxAngle           uint8
	yMinAngle           uint8
	yMaxAngle           uint8
	yCenterCompensation uint8

	fd *os.File
}

func NewCamPositionController(xAxisGPIO, yAxisGPIO uint8) *CamPositionController {
	return &CamPositionController{xAxisGPIO, yAxisGPIO, 0.05, 0.25,
		50, 140,
		30, 110,
		30,
		nil}
}

func (c CamPositionController) To(xAngle, yAngle uint8) error {

	xAngle = limit(xAngle, c.xMinAngle, c.xMaxAngle)
	if err := c.moveServoAngle(c.xAxisGPIO, xAngle); err != nil {
		return err
	}
	if yAngle >= c.yCenterCompensation {
		yAngle -= c.yCenterCompensation
	}
	yAngle = limit(yAngle, c.yMinAngle, c.yMaxAngle)
	return c.moveServoAngle(c.yAxisGPIO, yAngle)
}

func limit(value, min, max uint8) uint8 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func (c CamPositionController) moveServoAngle(gpio, angle uint8) error {
	pulse := c.MinPulse + (c.MaxPulse * (float32(angle) / 180.0))
	cmd := fmt.Sprintf("%d=%f\n", gpio, pulse)
	var err error
	if c.fd == nil {
		if c.fd, err = os.OpenFile("/dev/pi-blaster", os.O_WRONLY, os.ModeExclusive); err != nil {
			return err
		}
	}
	_, err = c.fd.WriteString(cmd)
	c.fd.Sync()
	return err
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
