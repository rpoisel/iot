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
	Out1 chan<- bool
	Out2 chan<- bool
}

func (s *Split2Bool) Process() {
	for in := range s.In {
		s.Out1 <- in
		s.Out2 <- in
	}
}

type NopBool struct {
	In <-chan bool
}

func (n *NopBool) Process() {
	for range n.In {
	}
}
