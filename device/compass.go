package device

import "github.com/davecheney/i2c"

type Direction string

const (
	North Direction = "North"
	West  Direction = "East"
	East  Direction = "East"
	South Direction = "South"

	ConfigurationRegisterA  = 0x00
	ConfigurationRegisterB  = 0x01
	ModeRegister            = 0x02
	AxisXDataRegisterMSB    = 0x03
	AxisXDataRegisterLSB    = 0x04
	AxisZDataRegisterMSB    = 0x05
	AxisZDataRegisterLSB    = 0x06
	AxisYDataRegisterMSB    = 0x07
	AxisYDataRegisterLSB    = 0x08
	StatusRegister          = 0x09
	IdentificationRegisterA = 0x10
	IdentificationRegisterB = 0x11
	IdentificationRegisterC = 0x12

	MeasurementContinuous = 0x00
	MeasurementSingleShot = 0x01
	MeasurementIdle       = 0x03
)

type Gauss float32

type Compass interface {
	Read() (*CompassData, error)
}

type RaspiCompass struct {
	I2CPort byte
	Address byte
	Scale   Gauss

	bus *i2c.I2C
}

type CompassData struct {
	X, Y, Z int
}

func NewCompass(i2cPort byte, address byte, scale Gauss) Compass {
	if scale == 0 {
		scale = 1.3
	}
	return &RaspiCompass{i2cPort, address, scale, nil}
}

func (c RaspiCompass) Read() (*CompassData, error) {
	if err := c.initialize(); err != nil {
		return nil, err
	}

	x, _ := c.readRegUInt16(AxisXDataRegisterMSB)
	z, _ := c.readRegUInt16(AxisZDataRegisterMSB)
	y, _ := c.readRegUInt16(AxisYDataRegisterMSB)

	return &CompassData{int(x), int(y), int(z)}, nil
}

func (c RaspiCompass) readRegUInt16(reg byte) (int16, error) {
	c.bus.WriteByte(reg)
	buf := make([]byte, 2)
	if _, err := c.bus.Read(buf); err != nil {
		return 0, err
	}
	value := int16(buf[0])<<8 + int16(buf[1])
	return value, nil
}

func (c *RaspiCompass) initialize() error {
	if c.bus == nil {
		var err error
		if c.bus, err = i2c.New(c.Address, int(c.I2CPort)); err != nil {
			return err
		}
		if _, err := c.bus.Write([]byte{ModeRegister, MeasurementContinuous}); err != nil {
			return err
		}
	}
	return nil
}
