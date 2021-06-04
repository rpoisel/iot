package config

type Config struct {
	MQTTUser   string `env:"MQTT_USER,required"`
	MQTTPass   string `env:"MQTT_PASS,required"`
	MQTTClient string `env:"MQTT_CLIENTID,required"`
	MQTTBroker string `env:"MQTT_BROKER,required"`
}
