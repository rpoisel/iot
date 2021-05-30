package i2c

import (
	"fmt"
	"time"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/rpoisel/iot/internal/config"
	"github.com/rpoisel/iot/internal/io"
)

type i2cBus struct {
	inDevs           []io.InputDevice8Bit
	scheduledOutDevs chan io.OutputDevice8Bit
}

func SetupI2C(busses config.I2C, ioDevices *io.Devices) error {
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	for busNum, busDevices := range busses {
		i2cBus, err := newI2CBus(busNum, busDevices, ioDevices)
		if err != nil {
			return fmt.Errorf("cannot create bus handle for bus \"%d\": %w", busNum, err)
		}
		go i2cBus.worker()
	}
	return nil
}

func newI2CBus(busNum int, bus config.I2CBus, ioDevices *io.Devices) (*i2cBus, error) {
	result := &i2cBus{
		inDevs:           make([]io.InputDevice8Bit, 0, len(bus.In)),
		scheduledOutDevs: make(chan io.OutputDevice8Bit),
	}

	for _, inputDevice := range bus.In {
		devHandle, err := goi2c.NewI2C(inputDevice.Address, busNum)
		if err != nil {
			return nil, err
		}
		chip := NewPCF8574Input(devHandle)
		ioDevices.AddInputDevice(inputDevice.Name, chip)
		result.inDevs = append(result.inDevs, chip)
	}
	for _, outputDevice := range bus.Out {
		devHandle, err := goi2c.NewI2C(outputDevice.Address, busNum)
		if err != nil {
			return nil, err
		}
		chip := NewPCF8574Output(devHandle)
		ioDevices.AddOutputDevice(outputDevice.Name, chip)
		go chip.Run(result.scheduledOutDevs)
	}
	return result, nil
}

// worker mainly synchronizes bus access
// read: every 50ms
// write; event-based
func (i *i2cBus) worker() {
	for {
		select {
		case <-time.After(50 * time.Millisecond):
			for _, device := range i.inDevs {
				device.Read()
			}

		case device := <-i.scheduledOutDevs:
			device.Write()
		}
	}
}
