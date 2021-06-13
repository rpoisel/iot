package model

type Node interface {
	Name() string
	GetWires() [][]Node
}

type NodeMQTTIn interface {
	Name() string
	Broker() string
}

type Visitor interface {
	MQTTIn(n NodeMQTTIn)
	Connection(from Node /* need more info - which output pin */, to Node)
}

type Model interface {
	Vist(visitor Visitor)
}
