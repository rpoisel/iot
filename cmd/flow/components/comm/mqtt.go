package comm

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTReceive struct {
	client mqtt.Client
	Out    chan<- string
}

func (m *MQTTReceive) Process() {
	m.Out <- "true"
	select {
	// just hanging out
	}
}

type MQTTPublish struct {
	client mqtt.Client
	In     <-chan string
}

func (m *MQTTPublish) Process() {

}
