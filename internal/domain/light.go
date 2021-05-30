package domain

type Light struct {
	state  bool
	Input  <-chan bool
	Output chan<- bool
}

func (l *Light) Run() {
	for {
		l.Output <- <-l.Input

		// input := <-l.Input
		// if !input {
		// 	continue
		// }
		// l.state = !l.state
		// l.Output <- l.state
	}
}
