package nodered

type MQTTIn struct {
	name   string
	broker string
}

func (n *MQTTIn) Name() string {
	return n.name
}

func (n *MQTTIn) Broker() string {
	return n.broker
}
