package device

import (
	"math"

	"github.com/davecheney/i2c"
)

// Direction is the compass orientation
type Direction string

// Direction types
const (
	North     Direction = "North"
	NorthEast Direction = "NorthEast"
	East      Direction = "East"
	SouthEast Direction = "SouthEast"
	South     Direction = "South"
	SouthWest Direction = "SouthWest"
	West      Direction = "West"
	NorthWest Direction = "NorthWest"
)

/*
	HMC5883L Compass Default values
	for reference read: page 11 of
		https://dlnmh9ip6v2uc.cloudfront.net/datasheets/Sensors/Magneto/HMC5883L-FDS.pdf
*/
const (
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

// Gauss measurement type
type Gauss float32

// Compass read and return orientation data
type Compass interface {
	Read() (*CompassData, error)
}

// CompassData stores diraction and x, y, z coordinates
type CompassData struct {
	Direction Direction
	Degress   int
	X, Y, Z   int
}

type hmc5883LCompass struct {
	I2CPort byte
	Address byte
	Scale   Gauss

	bus *i2c.I2C
}

// NewCompass create an hmc5883l compass instance
func NewCompass(i2cPort byte, address byte, scale Gauss) Compass {
	if scale == 0 {
		scale = 1.3
	}
	return &hmc5883LCompass{i2cPort, address, scale, nil}
}

func (c hmc5883LCompass) Read() (*CompassData, error) {
	if err := c.initialize(); err != nil {
		return nil, err
	}

	x, _ := c.readRegUInt16(AxisXDataRegisterMSB)
	z, _ := c.readRegUInt16(AxisZDataRegisterMSB)
	y, _ := c.readRegUInt16(AxisYDataRegisterMSB)
	degress := degressByXY(x, y)
	direction := directionByDegress(degress)

	return &CompassData{direction, degress, int(x), int(y), int(z)}, nil
}

func (c hmc5883LCompass) readRegUInt16(reg byte) (int16, error) {
	c.bus.WriteByte(reg)
	buf := make([]byte, 2)
	if _, err := c.bus.Read(buf); err != nil {
		return 0, err
	}
	value := int16(buf[0])<<8 + int16(buf[1])
	return value, nil
}

func (c *hmc5883LCompass) initialize() error {
	var err error
	if c.bus == nil {
		c.bus, err = i2c.New(c.Address, int(c.I2CPort))
		if err != nil {
			return err
		}
		_, err = c.bus.Write([]byte{ModeRegister, MeasurementContinuous})
	}
	return err
}

func degressByXY(x, y int16) int {
	heading := math.Atan2(float64(y), float64(x))
	// heading += 233.9 / 1000
	if heading < 0 {
		heading += 2 * math.Pi
	}
	if heading > 2*math.Pi {
		heading -= 2 * math.Pi
	}
	return int(heading * 180 / math.Pi)
}

func directionByDegress(deg int) Direction {
	if deg < 0 {
		deg = deg + 360
	} else if deg >= 360 {
		deg = deg - 360
	}
	directions := []Direction{North, NorthEast, East, SouthEast, South, SouthWest, West, NorthWest}
	directionsAngle := 360 / len(directions)
	return directions[deg/directionsAngle]
}
