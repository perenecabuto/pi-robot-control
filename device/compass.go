package device

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

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
type Compass struct {
	I2CPort byte
	Address byte
	Scale   Gauss

	bus embd.I2CBus
}

type CompassData struct {
	X, Y, Z int
}

func NewCompass(i2cPort byte, address byte, scale Gauss) *Compass {
	if scale == 0 {
		scale = 1.3
	}
	return &Compass{i2cPort, address, scale, nil}
}

func (c *Compass) Read() (*CompassData, error) {
	if c.bus == nil {
		c.bus = embd.NewI2CBus(c.I2CPort)
		// INITIALIZE
		// Wire.send(0x02); //select mode register
		// Wire.send(0x00); //continuous measurement mode
		if err := c.bus.WriteByteToReg(c.Address, ModeRegister, MeasurementContinuous); err != nil {
			return nil, err
		}
	}

	// GET DATA
	//  int x,y,z; //triple axis data
	// //Tell the HMC5883L where to begin reading data
	// Wire.beginTransmission(address);
	// Wire.send(0x03); //select register 3, X MSB register
	// Wire.endTransmission();
	// //Read data from each axis, 2 registers per axis
	// Wire.requestFrom(address, 6);
	// if(6<=Wire.available()){
	//   x = Wire.receive()<<8; //X msb
	//   x |= Wire.receive(); //X lsb
	//   z = Wire.receive()<<8; //Z msb
	//   z |= Wire.receive(); //Z lsb
	//   y = Wire.receive()<<8; //Y msb
	//   y |= Wire.receive(); //Y lsb
	// }
	x, _ := c.bus.ReadByteFromReg(c.Address, AxisXDataRegisterMSB)
	xL, _ := c.bus.ReadByteFromReg(c.Address, AxisXDataRegisterLSB)
	x <<= 8
	x |= xL
	z, _ := c.bus.ReadByteFromReg(c.Address, AxisZDataRegisterMSB)
	zL, _ := c.bus.ReadByteFromReg(c.Address, AxisZDataRegisterLSB)
	z <<= 8
	z |= zL
	y, _ := c.bus.ReadByteFromReg(c.Address, AxisYDataRegisterMSB)
	yL, _ := c.bus.ReadByteFromReg(c.Address, AxisYDataRegisterLSB)
	y <<= 8
	y |= yL

	return &CompassData{int(x), int(y), int(z)}, nil
}

func (c Compass) send(value byte) error {
	return nil
}
