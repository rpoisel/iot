package i2c

import (
	goi2c "github.com/d2r2/go-i2c"
)

type PCF8574 struct {
	handle *goi2c.I2C
	state  byte
}

// TODO err handling
func (p *PCF8574) Read() {
	newState := []byte{p.state}
	p.handle.ReadBytes(newState)
	p.state = newState[0]
}

func (p *PCF8574) NumInputs() uint16 { return 8 }

// TODO err handling
func (p *PCF8574) Write() {
	p.handle.WriteBytes([]byte{p.state})
}

func (p *PCF8574) NumOutputs() uint16 { return 8 }
