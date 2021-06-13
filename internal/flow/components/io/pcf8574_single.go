package io

import (
	"log"
	"sync"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/rpoisel/iot/internal/flow/util"
	"go.uber.org/zap"
)

type PCF8574OutMsg struct {
	Pin   uint8
	State bool
}

type PCF8574OutPin struct {
	pin uint8
	In  <-chan bool
	Out chan<- PCF8574OutMsg
}

func (p *PCF8574OutPin) Process() {
	for {
		p.Out <- PCF8574OutMsg{
			Pin:   p.pin,
			State: <-p.In,
		}

	}
}

func NewPCF8574OutPin(pin uint8) *PCF8574OutPin {
	return &PCF8574OutPin{
		pin: pin,
	}
}

type PCF8574OutSingle struct {
	mutex  *sync.Mutex
	handle *goi2c.I2C
	logger *zap.SugaredLogger
	state  byte
	In     <-chan PCF8574OutMsg
}

func NewPCF8574OutSingle(logger *zap.SugaredLogger, bus int, addr uint8, mutex *sync.Mutex) *PCF8574Out {
	handle, err := goi2c.NewI2C(addr, bus)
	if err != nil {
		log.Panicf("Could not create PCF8574 instance: %s\n", err)
	}
	return &PCF8574Out{
		mutex:  mutex,
		handle: handle,
		logger: logger,
		state:  0xff,
	}
}

func (p *PCF8574OutSingle) Process() {
	for {
		outMsg := <-p.In
		util.SetBit(&p.state, outMsg.Pin, !outMsg.State)
		p.mutex.Lock()
		_, err := p.handle.WriteBytes([]byte{p.state})
		p.mutex.Unlock()
		if err != nil {
			log.Printf("Could not read from I2C device: %s\n", err)
		}
	}
}
