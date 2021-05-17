package i2c

import (
	"time"

	goi2c "github.com/d2r2/go-i2c"
	"github.com/rpoisel/iot/internal/config"
)

type I2CBus struct {
	InDevs           []I2CInputDevice
	OutDevs          []I2COutputDevice
	ScheduledOutDevs chan I2COutputDevice
}

func NewI2CBus(busNum int, bus config.I2CBus) (*I2CBus, error) {
	result := &I2CBus{
		ScheduledOutDevs: make(chan I2COutputDevice),
		InDevs:           make([]I2CInputDevice, 0, len(bus.In)),
		OutDevs:          make([]I2COutputDevice, 0, len(bus.Out)),
	}

	for _, device := range bus.In {
		devHandle, err := goi2c.NewI2C(device.Address, busNum)
		if err != nil {
			return nil, err
		}
		result.InDevs = append(result.InDevs, &PCF8574{
			handle: devHandle,
		})
	}
	// start go routine for output aggregator
	for _, device := range bus.Out {
		devHandle, err := goi2c.NewI2C(device.Address, busNum)
		if err != nil {
			return nil, err
		}
		result.OutDevs = append(result.OutDevs, &PCF8574{
			handle: devHandle,
		})
		// start go routine for device
	}
	return result, nil
}

// worker mainly synchronizes bus access
// read: every 50ms
// write; event based
func (i *I2CBus) worker() {
	for {
		select {
		case <-time.After(50 * time.Millisecond):
			for _, device := range i.InDevs {
				device.Read()
			}
			/* write to InputValues channel */

		case device := <-i.ScheduledOutDevs:
			device.Write()
		}
	}
}

func (i *I2CBus) Start() {
	go i.worker()
}
