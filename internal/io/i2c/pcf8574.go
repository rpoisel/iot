package i2c

import (
	goi2c "github.com/d2r2/go-i2c"
	"github.com/rpoisel/iot/internal/io"
)

const (
	numInputs  = 8
	numOutputs = 8
)

type pcf8574Input struct {
	handle *goi2c.I2C
	state  byte
	inputs []chan bool
}

func NewPCF8574Input(handle *goi2c.I2C) *pcf8574Input {
	result := &pcf8574Input{
		handle: handle,
		inputs: make([]chan bool, 0, numInputs),
	}
	for idx := 0; idx < numInputs; idx++ {
		result.inputs = append(result.inputs, make(chan bool))
	}
	return result
}

type pcf8574Output struct {
	handle  *goi2c.I2C
	state   byte
	outputs []chan bool
}

func NewPCF8574Output(handle *goi2c.I2C) *pcf8574Output {
	result := &pcf8574Output{
		handle:  handle,
		outputs: make([]chan bool, 0, numInputs),
	}
	for idx := 0; idx < numInputs; idx++ {
		result.outputs = append(result.outputs, make(chan bool))
	}
	return result
}

func (p *pcf8574Input) Read() {
	newState := []byte{p.state}
	p.handle.ReadBytes(newState) // TODO err handling
	p.state = ^newState[0]
	for idx := 0; idx < numInputs; idx++ {
		val := p.state&(0x01<<idx) == (0x01 << idx)
		// non-blocking write
		select {
		case p.inputs[idx] <- val:
		default:
		}
	}
}

func (p *pcf8574Output) Run(outputDevices chan<- io.OutputDevice8Bit) {
	for {
		// blocking read
		select {
		case val := <-p.outputs[0]:
			SetBit(&p.state, 0, val)
		case val := <-p.outputs[1]:
			SetBit(&p.state, 1, val)
		case val := <-p.outputs[2]:
			SetBit(&p.state, 2, val)
		case val := <-p.outputs[3]:
			SetBit(&p.state, 3, val)
		case val := <-p.outputs[4]:
			SetBit(&p.state, 4, val)
		case val := <-p.outputs[5]:
			SetBit(&p.state, 5, val)
		case val := <-p.outputs[6]:
			SetBit(&p.state, 6, val)
		case val := <-p.outputs[7]:
			SetBit(&p.state, 7, val)
		}
		outputDevices <- p
	}
}

// TODO err handling
func (p *pcf8574Output) Write() {
	p.handle.WriteBytes([]byte{^p.state})
}

func (p *pcf8574Input) GetInputChannel(idx uint) chan bool {
	return p.inputs[idx]
}

func (p *pcf8574Output) GetOutputChannel(idx uint) chan bool {
	return p.outputs[idx]
}
