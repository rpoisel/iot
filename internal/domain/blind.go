package domain

type Blind struct {
	Input1  <-chan bool
	Input2  <-chan bool
	Output1 chan<- bool
	Output2 chan<- bool
}

func (b *Blind) Run() {
	for {
		select {
		case /* in := */ <-b.Input1:
		case /* in := */ <-b.Input2:
		}
		b.Output1 <- false // write output according to change in Inputs
		b.Output2 <- false
	}
}
