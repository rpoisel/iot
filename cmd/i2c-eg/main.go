package main

import (
	"log"
	"time"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

const (
	inputBadLicht = 1 << iota
	inputStiegeLicht
	inputKuecheLichtZeile
	input3
	input4
	input5
	input6
	input7
)

const (
	outputKlingel = 1 << iota
	outputStiegeLichtOben
	outputStiegeLichtUnten
	outputKuecheLichtZeile
	unused0
	unused1
	unused2
	unused3
)

const (
	i2haAddr = 0x20
	i2heAddr = 0x38
)

func toggle(data *byte, mask byte) {
	(*data) ^= mask
}

func isBitsSet(data byte, mask byte) bool {
	return data&mask == mask
}

func isBitsUnset(data byte, mask byte) bool {
	return data&mask == 0x00
}

func main() {
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	i2ha, err := goi2c.NewI2C(i2haAddr, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2ha.Close()
	i2he, err := goi2c.NewI2C(i2heAddr, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2he.Close()

	inputs := byte(0x00)
	outputs := byte(0xff)
	firstTime := true
	for {
		select {
		case <-time.After(time.Millisecond * 50):
			inputData := []byte{0x00}
			if _, err := i2he.ReadBytes(inputData); err != nil {
				log.Printf("Could read from I2C device: %s\n", err)
			}
			curInputs := inputData[0]

			if firstTime {
				firstTime = false
			} else {
				if isBitsUnset(curInputs, inputStiegeLicht) && isBitsSet(inputs, inputStiegeLicht) {
					toggle(&outputs, outputStiegeLichtOben)
					toggle(&outputs, outputStiegeLichtUnten)
				}
				if isBitsUnset(curInputs, inputKuecheLichtZeile) && isBitsSet(inputs, inputKuecheLichtZeile) {
					toggle(&outputs, outputKuecheLichtZeile)
				}
				_, err := i2ha.WriteBytes([]byte{outputs})
				if err != nil {
					log.Printf("Could not write to I2C device: %s\n", err)
				}
			}
			inputs = curInputs
		}
	}
}
