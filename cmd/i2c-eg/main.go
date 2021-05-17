package main

import (
	"log"
	"time"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

func main() {
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	i2ha, err := goi2c.NewI2C(0x20, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2ha.Close()
	i2he, err := goi2c.NewI2C(0x38, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2he.Close()

	inputs := []byte{0x00}
	outputs := []byte{0xff}
	firstTime := true
	for {
		select {
		case <-time.After(time.Millisecond * 50):
			curInputs := []byte{0x00}
			if _, err := i2he.ReadBytes(curInputs); err != nil {
				log.Printf("Could read from I2C device: %w\n", err)
			}
			if firstTime {
				firstTime = false
			} else {
				if curInputs[0]&0x02 == 0x00 && inputs[0]&0x02 == 0x02 {
					outputs[0] ^= 0x06
				}
				_, err := i2ha.WriteBytes(outputs)
				if err != nil {
					log.Printf("Could not write to I2C device: %w\n", err)
				}
			}
			inputs = curInputs
		}
	}
}
