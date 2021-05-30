package io

import (
	"regexp"
	"strconv"
)

const (
	pathPattern = `([^\[]+)\[(\d+)\]`
)

var (
	pathRegex = regexp.MustCompile(pathPattern)
)

type InputDevice8Bit interface {
	// Read is called by the containing I2CBus
	// in order to transfer the hardware state
	// to the representation
	Read()

	GetInputChannel(idx uint) chan bool
}

type OutputDevice8Bit interface {
	Run(outputDevices chan<- OutputDevice8Bit)
	// Write is called by the containing I2CBus
	// in order to transfer the current state of this
	// device to the hardware
	Write()

	GetOutputChannel(idx uint) chan bool
}

type Devices struct {
	in  map[string]InputDevice8Bit
	out map[string]OutputDevice8Bit
}

func NewDevices() *Devices {
	return &Devices{
		in:  map[string]InputDevice8Bit{},
		out: map[string]OutputDevice8Bit{},
	}
}

func (d *Devices) AddInputDevice(path string, device InputDevice8Bit) {
	d.in[path] = device
}

func (d *Devices) AddOutputDevice(path string, device OutputDevice8Bit) {
	d.out[path] = device
}

func (d *Devices) GetInputChannel(path string) <-chan bool {
	parts := pathRegex.FindStringSubmatch(path)
	idx, _ := strconv.ParseUint(parts[2], 10, 16) // TODO error handling
	return d.in[parts[1]].GetInputChannel(uint(idx))
}

func (d *Devices) GetOutputChannel(path string) chan<- bool {
	parts := pathRegex.FindStringSubmatch(path)
	idx, _ := strconv.ParseUint(parts[2], 10, 16) // TODO error handling
	return d.out[parts[1]].GetOutputChannel(uint(idx))
}
