package device

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

const ACTION_TIMEOUT_SEC = 5

type Robot struct {
	pin1F     embd.DigitalPin
	pin2F     embd.DigitalPin
	pin1B     embd.DigitalPin
	pin2B     embd.DigitalPin
	stopTimer *time.Timer
}

func NewRobot(gpio1F, gpio2F, gpio1B, gpio2B int) *Robot {
	pins := make([]embd.DigitalPin, 4, 4)
	for i, gpio := range []int{gpio1F, gpio2F, gpio1B, gpio2B} {
		pins[i], _ = embd.NewDigitalPin(gpio)
	}

	timer := time.NewTimer(ACTION_TIMEOUT_SEC)
	return &Robot{pins[0], pins[1], pins[2], pins[3], timer}
}

func (r Robot) Initialize() error {
	r.CleanUP()
	if err := embd.InitGPIO(); err != nil {
		return err
	}

	for _, pin := range r.Pins() {
		if err := pin.SetDirection(embd.Out); err != nil {
			return err
		}
	}

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

func (r Robot) CleanUP() {
	log.Println("Deactivating robot...")
	for _, pin := range r.Pins() {
		if err := pin.Close(); err != nil {
			log.Println("CleanUP ERROR - close:", err.Error())
		}
	}
	embd.CloseGPIO()
}

func (r Robot) Pins() []embd.DigitalPin {
	return []embd.DigitalPin{r.pin1F, r.pin1B, r.pin2F, r.pin2B}
}

func (r Robot) move(val1F, val2F, val1B, val2B bool) {
	go r.autoStop()
	toggle(r.pin1F, val1F)
	toggle(r.pin2F, val2F)
	toggle(r.pin1B, val1B)
	toggle(r.pin2B, val2B)
}

func (r Robot) autoStop() {
	if !r.stopTimer.Stop() {
		<-r.stopTimer.C
	}
	r.stopTimer.Reset(ACTION_TIMEOUT_SEC)
}

func toggle(pin embd.DigitalPin, on bool) {
	io := embd.Low
	if on {
		io = embd.High
	}
	if err := pin.Write(io); err != nil {
		log.Println("Error:", err.Error())
	}
}
