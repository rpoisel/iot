package comm

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MQTTReceive struct {
	logger *zap.SugaredLogger
	client mqtt.Client
	topic  string
	Out    chan<- string
}

func NewMQTTReceive(logger *zap.SugaredLogger, client mqtt.Client, topic string) *MQTTReceive {
	return &MQTTReceive{
		logger: logger,
		client: client,
		topic:  topic,
	}
}

func (m *MQTTReceive) Process() {
	if token := m.client.Subscribe(m.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		m.Out <- string(msg.Payload())
		m.logger.Debug("received message",
			zap.String("topic", msg.Topic()), zap.String("message", string(msg.Payload())))
	}); token.Wait() && token.Error() != nil {
		m.logger.Panicf("Cannot subscribe: %s", token.Error())
	}
	select {
	// wait forever
	}
}

type MQTTPublish struct {
	logger *zap.SugaredLogger
	client mqtt.Client
	topic  string
	In     <-chan string
}

func NewMQTTPublish(logger *zap.SugaredLogger, client mqtt.Client, topic string) *MQTTPublish {
	return &MQTTPublish{
		logger: logger,
		client: client,
		topic:  topic,
	}
}

func (m *MQTTPublish) Process() {
	for {
		msg := <-m.In
		token := m.client.Publish(m.topic, 0, true, msg)
		token.Wait()
		m.logger.Debug("published message",
			zap.String("topic", m.topic), zap.String("message", msg))
	}
}
