package main

import (
	"fmt"
	"time"

	"github.com/rpoisel/iot/internal/flow/components/comm"
	"github.com/rpoisel/iot/internal/flow/components/homeautomation"
	"github.com/rpoisel/iot/internal/flow/components/io"
	"github.com/rpoisel/iot/internal/flow/components/logic"
	"github.com/rpoisel/iot/internal/flow/config"
	"github.com/rpoisel/iot/internal/flow/entrypoint"
	"github.com/rpoisel/iot/internal/flow/graph"
	"go.uber.org/zap"
)

func newHomeautomationApp(logger *zap.SugaredLogger, cfg *config.Config, componentsData *entrypoint.ComponentsData) *graph.Graph {
	var pinsDefaultHandling = [...]uint{0, 2, 3, 4, 6, 7}
	n := graph.NewGraph()

	n.AddOrPanic("pcf8574in", io.NewPCF8574In(logger, 1, 0x38, componentsData.I2CBusMutex))
	n.AddOrPanic("pcf8574out", io.NewPCF8574Out(logger, 1, 0x20, componentsData.I2CBusMutex))
	n.AddOrPanic("MQTTout", comm.NewMQTTPublish(logger, componentsData.MQTTClient, "/homeautomation/lights/KellerStiege/toggle"))
	n.AddOrPanic("trigger1", new(logic.R_Trig))
	n.AddOrPanic("trigger5", new(logic.R_Trig))
	n.AddOrPanic("convert1", new(logic.BoolToString))
	n.AddOrPanic("MultiClick1", logic.NewMultiClick(300*time.Millisecond))
	n.AddOrPanic("MultiClick5", logic.NewMultiClick(300*time.Millisecond))

	for pin := range pinsDefaultHandling {
		n.AddOrPanic(fmt.Sprintf("trigger%d", pin), new(logic.R_Trig))
		n.AddOrPanic(fmt.Sprintf("light%d", pin), new(homeautomation.Light))

		n.ConnectOrPanic("pcf8574in", fmt.Sprintf("Pin%d", pin),
			fmt.Sprintf("trigger%d", pin), "In")
		n.ConnectOrPanic(fmt.Sprintf("trigger%d", pin), "Out",
			fmt.Sprintf("light%d", pin), "In")
		n.ConnectOrPanic(fmt.Sprintf("light%d", pin), "Out",
			"pcf8574out", fmt.Sprintf("Pin%d", pin))
	}

	n.ConnectOrPanic("pcf8574in", "Pin1", "trigger1", "In")
	n.ConnectOrPanic("trigger1", "Out", "MultiClick0", "In")
	n.ConnectOrPanic("MultiClick1", "Out", "convert1", "In")
	n.ConnectOrPanic("convert1", "Out", "MQTTout", "In")

	n.ConnectOrPanic("pcf8574in", "Pin5", "trigger5", "In")
	n.ConnectOrPanic("trigger5", "Out", "MultiClick1", "In")

	return n
}

func main() {
	entrypoint.Entrypoint(newHomeautomationApp)
}
