package domain

import (
	"github.com/rpoisel/iot/internal/config"
	"github.com/rpoisel/iot/internal/io"
)

// StartAutomation performs the following steps:
// - reads from the configuration
// - starts workers for automation components
func StartAutomation(automation *config.Automation, devices *io.Devices) {
	for _, blindConfig := range automation.Blinds {
		blind := &Blind{
			Input1:  devices.GetInputChannel(blindConfig.Input1),
			Input2:  devices.GetInputChannel(blindConfig.Input2),
			Output1: devices.GetOutputChannel(blindConfig.Output1),
			Output2: devices.GetOutputChannel(blindConfig.Output2),
		}
		go blind.Run()
	}
	for _, lightConfig := range automation.Lights {
		light := &Light{
			Input:  devices.GetInputChannel(lightConfig.Inputs.Local.Name),
			Output: devices.GetOutputChannel(lightConfig.Output),
		}
		go light.Run()
	}
}
