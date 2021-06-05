package logic

type EqualString struct {
	cmp string
	In  <-chan string
	Out chan<- bool
}

func NewCompareString(cmp string) *EqualString {
	return &EqualString{
		cmp: cmp,
	}
}

func (e *EqualString) Process() {
	for {
		value := <-e.In
		e.Out <- value == e.cmp
	}
}

type EqualUint struct {
	cmp uint
	In  <-chan uint
	Out chan<- bool
}

func NewEqualUint(cmp uint) *EqualUint {
	return &EqualUint{
		cmp: cmp,
		In:  make(<-chan uint),
		Out: make(chan<- bool),
	}
}

func (e *EqualUint) Process() {
	for {
		value := <-e.In
		e.Out <- value == e.cmp
	}
}
