package io

import (
	"log"
	"sync"
	"time"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/rpoisel/iot/cmd/flow/util"
	"go.uber.org/zap"
)

type PCF8574In struct {
	logger *zap.SugaredLogger
	mutex  *sync.Mutex
	handle *goi2c.I2C
	ticker *time.Ticker
	Pin0   chan<- bool
	Pin1   chan<- bool
	Pin2   chan<- bool
	Pin3   chan<- bool
	Pin4   chan<- bool
	Pin5   chan<- bool
	Pin6   chan<- bool
	Pin7   chan<- bool
}

func NewPCF8574In(logger *zap.SugaredLogger, bus int, addr uint8, mutex *sync.Mutex) *PCF8574In {
	handle, err := goi2c.NewI2C(addr, bus)
	if err != nil {
		log.Panicf("Could not create PCF8574 instance: %s\n", err)
	}
	return &PCF8574In{
		mutex:  mutex,
		handle: handle,
		logger: logger,
		ticker: time.NewTicker(50 * time.Millisecond),
	}
}

func (p *PCF8574In) Process() {
	for {
		select {
		case <-p.ticker.C:
			inputData := []byte{0x00}
			p.mutex.Lock()
			_, err := p.handle.ReadBytes(inputData)
			p.mutex.Unlock()
			if err != nil {
				log.Printf("Could not read from I2C device: %s\n", err)
			}
			p.Pin0 <- !(inputData[0]&(0x01<<0) == (0x01 << 0))
			p.Pin1 <- !(inputData[0]&(0x01<<1) == (0x01 << 1))
			p.Pin2 <- !(inputData[0]&(0x01<<2) == (0x01 << 2))
			p.Pin3 <- !(inputData[0]&(0x01<<3) == (0x01 << 3))
			p.Pin4 <- !(inputData[0]&(0x01<<4) == (0x01 << 4))
			p.Pin5 <- !(inputData[0]&(0x01<<5) == (0x01 << 5))
			p.Pin6 <- !(inputData[0]&(0x01<<6) == (0x01 << 6))
			p.Pin7 <- !(inputData[0]&(0x01<<7) == (0x01 << 7))
		}
	}
}

type PCF8574Out struct {
	mutex  *sync.Mutex
	handle *goi2c.I2C
	logger *zap.SugaredLogger
	state  byte
	Pin0   <-chan bool
	Pin1   <-chan bool
	Pin2   <-chan bool
	Pin3   <-chan bool
	Pin4   <-chan bool
	Pin5   <-chan bool
	Pin6   <-chan bool
	Pin7   <-chan bool
}

func NewPCF8574Out(logger *zap.SugaredLogger, bus int, addr uint8, mutex *sync.Mutex) *PCF8574Out {
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

func (p *PCF8574Out) Process() {
	for {
		select {
		case value := <-p.Pin0:
			util.SetBit(&p.state, 0, !value)
		case value := <-p.Pin1:
			util.SetBit(&p.state, 1, !value)
		case value := <-p.Pin2:
			util.SetBit(&p.state, 2, !value)
		case value := <-p.Pin3:
			util.SetBit(&p.state, 3, !value)
		case value := <-p.Pin4:
			util.SetBit(&p.state, 4, !value)
		case value := <-p.Pin5:
			util.SetBit(&p.state, 5, !value)
		case value := <-p.Pin6:
			util.SetBit(&p.state, 6, !value)
		case value := <-p.Pin7:
			util.SetBit(&p.state, 7, !value)
		}
		p.mutex.Lock()
		_, err := p.handle.WriteBytes([]byte{p.state})
		p.mutex.Unlock()
		if err != nil {
			log.Printf("Could not read from I2C device: %s\n", err)
		}
	}
}
