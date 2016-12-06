package device

import "testing"

var (
	i2cPort = byte(1)
	address = byte(0x1E)
)

func TestNewCompass(t *testing.T) {
	compass := NewCompass(i2cPort, address, 1)
	if compass == nil {
		t.Fatal("Can not create new compass")
	}
}
