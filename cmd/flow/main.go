package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rpoisel/iot/cmd/flow/graph"
	"github.com/rpoisel/iot/cmd/flow/logic"
)

type Emitter struct {
	Out1 chan<- bool
	Out2 chan<- bool
}

func (e *Emitter) Process() {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			val := rand.Int()%2 == 0
			fmt.Printf("Sender: %v\n", val)
			e.Out1 <- val
		}
	}
}

type Receiver struct {
	In <-chan bool
}

func (r *Receiver) Process() {
	for in := range r.In {
		fmt.Printf("Receiver: %v\n", in)
	}
}

func newBinaryApp() *graph.Graph {
	n := graph.NewGraph()

	n.AddPanic("emitter", new(Emitter))
	n.AddPanic("trigger", new(logic.R_Trig))
	n.AddPanic("receiver", new(Receiver))
	n.AddPanic("nop", new(logic.NopBool))

	n.ConnectPanic("emitter", "Out1", "trigger", "In")
	n.ConnectPanic("emitter", "Out2", "nop", "In")
	n.ConnectPanic("trigger", "Out", "receiver", "In")

	return n
}

func main() {
	net := newBinaryApp()

	wait := net.Run()

	<-wait
}
