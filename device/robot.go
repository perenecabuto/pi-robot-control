package device

import (
	"log"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

const ACTION_TIMEOUT_SEC = 5

type Robot struct {
	pin1F *gpio.DirectPinDriver
	pin2F *gpio.DirectPinDriver
	pin1B *gpio.DirectPinDriver
	pin2B *gpio.DirectPinDriver
}

var (
	gbot = gobot.NewGobot()
	rasp = raspi.NewRaspiAdaptor("raspi")
)

func NewRobot(gpio1F, gpio2F, gpio1B, gpio2B int) *Robot {
	rasp.Connect()
	gpios := []int{gpio1F, gpio2F, gpio1B, gpio2B}
	pins := make([]*gpio.DirectPinDriver, 4)
	for i, g := range gpios {
		name := strconv.Itoa(g)
		pins[i] = gpio.NewDirectPinDriver(rasp, "gpio"+name, name)
	}
	return &Robot{pins[0], pins[1], pins[2], pins[3]}
}

func (r Robot) Pins() []*gpio.DirectPinDriver {
	return []*gpio.DirectPinDriver{r.pin1F, r.pin2F, r.pin1B, r.pin2B}
}

func (r Robot) Initialize() error {
	robot := gobot.NewRobot("robot", []gobot.Connection{rasp})
	gbot.AddRobot(robot)
	for _, pin := range r.Pins() {
		if err := pin.Start(); err != nil {
			log.Panicf("%#v\n", err)
		}
		robot.AddDevice(pin)
	}

	go gbot.Start()
	r.Stop()

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

func toggle(pin *gpio.DirectPinDriver, on bool) {
	pin.Start()
	if on {
		if err := pin.On(); err != nil {
			log.Println(err.Error())
		}
	} else {
		if err := pin.Off(); err != nil {
			log.Println(err.Error())
		}
	}
}
