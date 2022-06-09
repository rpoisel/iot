package graph

import (
	"log"

	"github.com/trustmaster/goflow"
)

type Graph struct {
	n *goflow.Graph
}

func NewGraph() *Graph {
	g := &Graph{
		n: goflow.NewGraph(),
	}

	return g
}

func (g *Graph) Add(name string, c interface{}) {
	if err := g.n.Add(name, c); err != nil {
		log.Panicf("Coult not add: %s", err)
	}
}

func (g *Graph) Connect(senderName, senderPort, receiverName, receiverPort string) {
	if err := g.n.Connect(senderName, senderPort, receiverName, receiverPort); err != nil {
		log.Panicf("Could not connect: %s", err)
	}
}

func (g *Graph) Run() goflow.Wait {
	return goflow.Run(g.n)
}
