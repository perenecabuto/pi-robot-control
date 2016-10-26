package device

import "testing"

var (
	i2cPort = 1
	address = 0x1E

	compass *Compass
)

func TestNewCompass(t *testing.T) {
	compass = NewCompass(i2cPort, address, 1)
	if compass == nil {
		t.Fatal("Can not create new compass")
	}
}
