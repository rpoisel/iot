package comm

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTReceive struct {
	client mqtt.Client
	topic  string
	Out    chan<- string
}

func NewMQTTReceive(client mqtt.Client, topic string) *MQTTReceive {
	return &MQTTReceive{
		client: client,
		topic:  topic,
	}
}

func (m *MQTTReceive) Process() {
	if token := m.client.Subscribe(m.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		m.Out <- string(msg.Payload())
	}); token.Wait() && token.Error() != nil {
		log.Panicf("Cannot subscribe: %s", token.Error())
	}
	select {
	// wait forever
	}
}

type MQTTPublish struct {
	client mqtt.Client
	topic  string
	In     <-chan string
}

func NewMQTTPublish(client mqtt.Client, topic string) *MQTTPublish {
	return &MQTTPublish{
		client: client,
		topic:  topic,
	}
}

func (m *MQTTPublish) Process() {
	for {
		msg := <-m.In
		token := m.client.Publish(m.topic, 0, true, msg)
		token.Wait()
	}
}
