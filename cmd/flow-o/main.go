package main

import (
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
	n := graph.NewGraph()

	n.AddOrPanic("pcf8574in", io.NewPCF8574In(logger, 1, 0x38, componentsData.I2CBusMutex))
	n.AddOrPanic("pcf8574out", io.NewPCF8574Out(logger, 1, 0x20, componentsData.I2CBusMutex))
	n.AddOrPanic("MQTTin", comm.NewMQTTReceive(logger, componentsData.MQTTClient, "/homeautomation/lights/KuecheZeile/toggle"))
	n.AddOrPanic("MQTTout", comm.NewMQTTPublish(logger, componentsData.MQTTClient, "/homeautomation/lights/KuecheZeile/state"))
	n.AddOrPanic("convert", new(logic.StringToBool))
	n.AddOrPanic("triggerStiegeLicht", new(logic.R_Trig))
	n.AddOrPanic("triggerKuecheLichtZeile", new(logic.R_Trig))
	n.AddOrPanic("lightStiege", new(homeautomation.Light))
	n.AddOrPanic("lightKuecheZeile", new(homeautomation.Light))
	n.AddOrPanic("splitLightStiege", new(logic.Split2Bool))
	n.AddOrPanic("splitLightKuecheZeile", new(logic.Split2Bool))
	n.AddOrPanic("bool2string", new(logic.BoolToString))
	n.AddOrPanic("nop", new(logic.NopBool))

	n.ConnectOrPanic("pcf8574in", "Pin0", "nop", "In")
	n.ConnectOrPanic("pcf8574in", "Pin1", "triggerStiegeLicht", "In")
	n.ConnectOrPanic("pcf8574in", "Pin2", "triggerKuecheLichtZeile", "In")
	n.ConnectOrPanic("pcf8574in", "Pin3", "nop", "In")
	n.ConnectOrPanic("pcf8574in", "Pin4", "nop", "In")
	n.ConnectOrPanic("pcf8574in", "Pin5", "nop", "In")
	n.ConnectOrPanic("pcf8574in", "Pin6", "nop", "In")
	n.ConnectOrPanic("pcf8574in", "Pin7", "nop", "In")
	n.ConnectOrPanic("triggerStiegeLicht", "Out", "lightStiege", "In")
	n.ConnectOrPanic("lightStiege", "Out", "splitLightStiege", "In")
	n.ConnectOrPanic("splitLightStiege", "Out0", "pcf8574out", "Pin1")
	n.ConnectOrPanic("splitLightStiege", "Out1", "pcf8574out", "Pin2")
	n.ConnectOrPanic("triggerKuecheLichtZeile", "Out", "lightKuecheZeile", "In")
	n.ConnectOrPanic("MQTTin", "Out", "convert", "In")
	n.ConnectOrPanic("convert", "Out", "lightKuecheZeile", "In")
	n.ConnectOrPanic("lightKuecheZeile", "Out", "splitLightKuecheZeile", "In")
	n.ConnectOrPanic("splitLightKuecheZeile", "Out0", "pcf8574out", "Pin3")
	n.ConnectOrPanic("splitLightKuecheZeile", "Out1", "bool2string", "In")
	n.ConnectOrPanic("bool2string", "Out", "MQTTout", "In")

	return n
}

func main() {
	entrypoint.Entrypoint(newHomeautomationApp)
}
