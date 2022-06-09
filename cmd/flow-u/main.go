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

	n.Add("pcf8574in", io.NewPCF8574In(logger, 1, 0x38, componentsData.I2CBusMutex))
	n.Add("pcf8574out", io.NewPCF8574Out(logger, 1, 0x20, componentsData.I2CBusMutex))
	n.Add("MQTTout", comm.NewMQTTPublish(logger, componentsData.MQTTClient, "/homeautomation/lights/KellerStiege/toggle"))
	n.Add("trigger1", new(logic.R_Trig))
	n.Add("trigger5", new(logic.R_Trig))
	n.Add("convert1", new(logic.BoolToString))
	n.Add("MultiClick1", logic.NewMultiClick(300*time.Millisecond))
	n.Add("MultiClick5", logic.NewMultiClick(300*time.Millisecond))

	for pin := range pinsDefaultHandling {
		n.Add(fmt.Sprintf("trigger%d", pin), new(logic.R_Trig))
		n.Add(fmt.Sprintf("light%d", pin), new(homeautomation.Light))

		n.Connect("pcf8574in", fmt.Sprintf("Pin%d", pin),
			fmt.Sprintf("trigger%d", pin), "In")
		n.Connect(FMT.Sprintf("trigger%d", pin), "Out",
			fmt.Sprintf("light%d", pin), "In")
		n.Connect(fmt.Sprintf("light%d", pin), "Out",
			"pcf8574out", fmt.Sprintf("Pin%d", pin))
	}

	n.Connect("pcf8574in", "Pin1", "trigger1", "In")
	n.Connect("trigger1", "Out", "MultiClick0", "In")
	n.Connect("MultiClick1", "Out", "convert1", "In")
	n.Connect("convert1", "Out", "MQTTout", "In")

	n.Connect("pcf8574in", "Pin5", "trigger5", "In")
	n.Connect("trigger5", "Out", "MultiClick1", "In")

	return n
}

func main() {
	entrypoint.Entrypoint(newHomeautomationApp)
}
