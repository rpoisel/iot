package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
)

// MqttConfiguration is a standard configuration that should be used in all IoT modules.
type MqttConfiguration struct {
	Broker   string
	Username string
	Password string
	BasePath string
}

// ReadConfig reads configuration from a configuration file.
func ReadConfig(path string, config interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

// SetupMqtt intializes the MQTT connection.
func SetupMqtt(config MqttConfiguration, defaultMsgHandler MQTT.MessageHandler, onConnectHandler MQTT.OnConnectHandler) (client MQTT.Client) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	if defaultMsgHandler != nil {
		opts.SetDefaultPublishHandler(defaultMsgHandler)
	}
	if onConnectHandler != nil {
		opts.SetOnConnectHandler(onConnectHandler)
	}
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

// Readings contains the values to be read from the power meters
type Readings struct {
	Solar    int32
	Obtained int32
	Total    int32
}

// NewReadings initializes a Readings data structure from a byte buffer.
func NewReadings(buf []byte) (r *Readings, err error) {
	if len(buf) != 12 {
		return nil, errors.New("Given buffer has size != 12")
	}
	r = new(Readings)
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, r)
	return
}

// ToBuf converts a Readings instance to a bytes buffer.
func (r Readings) ToBuf() (buf []byte) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, r)
	return b.Bytes()
}
