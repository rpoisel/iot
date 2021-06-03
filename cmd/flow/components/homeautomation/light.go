package homeautomation

type Light struct {
	state bool
	In    <-chan interface{}
	Out   chan<- bool
}

func (l *Light) Process() {
	for {
		<-l.In
		l.state = !l.state
		l.Out <- l.state
	}
}
