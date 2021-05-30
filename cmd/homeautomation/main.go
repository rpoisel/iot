package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/rpoisel/iot/internal/config"
	"github.com/rpoisel/iot/internal/domain"
	"github.com/rpoisel/iot/internal/io"
	"github.com/rpoisel/iot/internal/io/i2c"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := flag.String("c", "configuration.yml", "path to the configuration file")
	flag.Parse()

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Cannot read configuration file \"%s\": %s\n", *configPath, err)
	}

	config := config.Configuration{}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Cannot parse configuration file \"%s\": %s\n", *configPath, err)
	}

	ioDevices := io.NewDevices()
	if err := i2c.SetupI2C(config.IO.I2C, ioDevices); err != nil {
		log.Fatalf("Cannot setup IOs: %s\n", err)
	}

	domain.StartAutomation(&config.Automation, ioDevices)

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
}
