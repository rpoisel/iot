package domain

type Light struct {
	state  bool
	Input  <-chan bool
	Output chan<- bool
	MQTT   <-chan string
}

func (l *Light) Run() {
	for {
		select {
		case input := <-l.Input:
			l.Output <- input
		case /* command := */ <-l.MQTT:
		}
	}
}
