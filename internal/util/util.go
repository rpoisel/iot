package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfiguration struct {
	Broker   string
	Username string
	Password string
}

func ReadConfig(path string, config interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func SetupMqtt(config MqttConfiguration, defaultHandler MQTT.MessageHandler) (client MQTT.Client) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	if defaultHandler != nil {
		opts.SetDefaultPublishHandler(defaultHandler)
	}
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

type Readings struct {
	Solar    int32
	Obtained int32
	Total    int32
}

func NewReadings(buf []byte) (r *Readings, err error) {
	if len(buf) != 12 {
		return nil, errors.New("Given buffer has size != 12")
	}
	r = new(Readings)
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, r)
	return
}

func (r Readings) ToBuf() (buf []byte) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, r)
	return b.Bytes()
}
