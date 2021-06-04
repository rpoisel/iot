package entrypoint

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/d2r2/go-logger"
	i2clogger "github.com/d2r2/go-logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rpoisel/iot/internal/flow/config"
	"github.com/rpoisel/iot/internal/flow/graph"
	"go.uber.org/zap"
)

type GraphGenerator func(logger *zap.SugaredLogger, cfg *config.Config, componentsData *ComponentsData) *graph.Graph

type ComponentsData struct {
	I2CBusMutex *sync.Mutex
	MQTTClient  mqtt.Client
}

func Entrypoint(generator GraphGenerator) {
	i2clogger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panicf("%+v\n", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Could not instantiate logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	net := generator(sugar, &cfg, initializeComponentsData(sugar, &cfg))

	wait := net.Run()

	<-wait
}

func initializeComponentsData(logger *zap.SugaredLogger, cfg *config.Config) *ComponentsData {
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
	return &ComponentsData{
		I2CBusMutex: &sync.Mutex{},
		MQTTClient:  mqttClient,
	}
}
