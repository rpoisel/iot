package logic

type StringToBool struct {
	In  <-chan string
	Out chan<- bool
}

func (s *StringToBool) Process() {
	for {
		value := <-s.In
		if value == "true" {
			s.Out <- true
		}
		if value == "false" {
			s.Out <- false
		}
	}
}

type BoolToString struct {
	In  <-chan bool
	Out chan<- string
}

func (b *BoolToString) Process() {
	for {
		value := <-b.In
		if value {
			b.Out <- "true"
		} else {
			b.Out <- "false"
		}
	}
}
