package main

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/d2r2/go-logger"
	i2clogger "github.com/d2r2/go-logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rpoisel/iot/cmd/flow/components/comm"
	"github.com/rpoisel/iot/cmd/flow/components/homeautomation"
	"github.com/rpoisel/iot/cmd/flow/components/io"
	"github.com/rpoisel/iot/cmd/flow/components/logic"
	"github.com/rpoisel/iot/cmd/flow/graph"
	"go.uber.org/zap"
)

type config struct {
	MQTTUser   string `env:"MQTT_USER,required"`
	MQTTPass   string `env:"MQTT_PASS,required"`
	MQTTClient string `env:"MQTT_CLIENTID,required"`
	MQTTBroker string `env:"MQTT_BROKER,required"`
}

func newHomeautomationApp(cfg *config, logger *zap.SugaredLogger) *graph.Graph {
	n := graph.NewGraph()

	m := &sync.Mutex{}
	opts := mqtt.NewClientOptions().AddBroker(cfg.MQTTBroker)
	opts.SetUsername(cfg.MQTTUser)
	opts.SetPassword(cfg.MQTTPass)
	opts.SetClientID(cfg.MQTTClient)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {})
	opts.SetPingTimeout(1 * time.Second)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logger.Panic(token.Error())
	}

	n.AddOrPanic("pcf8574in", io.NewPCF8574In(logger, 1, 0x38, m))
	n.AddOrPanic("pcf8574out", io.NewPCF8574Out(logger, 1, 0x20, m))
	n.AddOrPanic("MQTTin", comm.NewMQTTReceive(logger, mqttClient, "/homeautomation/lights/KuecheZeile/toggle"))
	n.AddOrPanic("MQTTout", comm.NewMQTTPublish(logger, mqttClient, "/homeautomation/lights/KuecheZeile/state"))
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
	i2clogger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panicf("%+v\n", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Could not instantiate logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	net := newHomeautomationApp(&cfg, sugar)

	wait := net.Run()

	<-wait
}
