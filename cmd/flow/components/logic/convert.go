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
