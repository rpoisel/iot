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

	n.Add("pcf8574in", io.NewPCF8574In(logger, 1, 0x38, componentsData.I2CBusMutex))
	n.Add("pcf8574out", io.NewPCF8574Out(logger, 1, 0x20, componentsData.I2CBusMutex))
	n.Add("MQTTin", comm.NewMQTTReceive(logger, componentsData.MQTTClient, "/homeautomation/lights/KuecheZeile/toggle"))
	n.Add("MQTTout", comm.NewMQTTPublish(logger, componentsData.MQTTClient, "/homeautomation/lights/KuecheZeile/state"))
	n.Add("convert", new(logic.StringToBool))
	n.Add("triggerStiegeLicht", new(logic.R_Trig))
	n.Add("triggerKuecheLichtZeile", new(logic.R_Trig))
	n.Add("triggerCharger", new(logic.R_Trig))
	n.Add("lightStiege", new(homeautomation.Light))
	n.Add("lightKuecheZeile", new(homeautomation.Light))
	n.Add("charger", new(homeautomation.Light))
	n.Add("splitLightStiege", new(logic.Split2Bool))
	n.Add("splitLightKuecheZeile", new(logic.Split2Bool))
	n.Add("bool2string", new(logic.BoolToString))
	n.Add("nop", new(logic.NopBool))

	n.Connect("pcf8574in", "Pin0", "nop", "In")
	n.Connect("pcf8574in", "Pin1", "triggerStiegeLicht", "In")
	n.Connect("pcf8574in", "Pin2", "triggerKuecheLichtZeile", "In")
	n.Connect("pcf8574in", "Pin3", "triggerCharger", "In")
	n.Connect("pcf8574in", "Pin4", "nop", "In")
	n.Connect("pcf8574in", "Pin5", "nop", "In")
	n.Connect("pcf8574in", "Pin6", "nop", "In")
	n.Connect("pcf8574in", "Pin7", "nop", "In")
	n.Connect("triggerStiegeLicht", "Out", "lightStiege", "In")
	n.Connect("lightStiege", "Out", "splitLightStiege", "In")
	n.Connect("splitLightStiege", "Out0", "pcf8574out", "Pin1")
	n.Connect("splitLightStiege", "Out1", "pcf8574out", "Pin2")
	n.Connect("triggerKuecheLichtZeile", "Out", "lightKuecheZeile", "In")
	n.Connect("MQTTin", "Out", "convert", "In")
	n.Connect("convert", "Out", "lightKuecheZeile", "In")
	n.Connect("lightKuecheZeile", "Out", "splitLightKuecheZeile", "In")
	n.Connect("splitLightKuecheZeile", "Out0", "pcf8574out", "Pin3")
	n.Connect("splitLightKuecheZeile", "Out1", "bool2string", "In")
	n.Connect("bool2string", "Out", "MQTTout", "In")
	n.Connect("triggerCharger", "Out", "charger", "In")
	n.Connect("charger", "Out", "pcf8574out", "Pin0")

	return n
}

func main() {
	entrypoint.Entrypoint(newHomeautomationApp)
}
