package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/rpoisel/iot/internal/config"
	"github.com/rpoisel/iot/internal/i2c"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := flag.String("c", "configuration.yml", "path to the configuration file")
	flag.Parse()

	config := config.Configuration{}
	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf(`Cannot read configuration file "%s": %v\n`, *configPath, err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf(`Cannot parse configuration file "%s": %v\n`, *configPath, err)
	}
	for busNum, devices := range config.I2C {
		i2cHandle, err := i2c.NewI2CBus(busNum, devices)
		if err != nil {
			log.Fatalf(`Cannot create bus handle for bus "%d": %v\n`, busNum, err)
		}
		i2cHandle.Start()
	}
}
