package logic

type R_Trig struct {
	lastIn bool
	In     <-chan bool
	Out    chan<- bool
}

func (r *R_Trig) Process() {
	for in := range r.In {
		if !r.lastIn && in {
			r.Out <- true
		}
		r.lastIn = in
	}
}

type Split2Bool struct {
	In   <-chan bool
	Out0 chan<- bool
	Out1 chan<- bool
}

func (s *Split2Bool) Process() {
	for in := range s.In {
		s.Out0 <- in
		s.Out1 <- in
	}
}

type NopBool struct {
	In <-chan interface{}
}

func (n *NopBool) Process() {
	for range n.In {
		// discard all received messages
	}
}
