package logic

type CompareString struct {
	cmp string
	In  <-chan string
	Out chan<- bool
}

func NewCompareString(cmp string) *CompareString {
	return &CompareString{
		cmp: cmp,
	}
}

func (c *CompareString) Process() {
	for {
		value := <-c.In
		c.Out <- value == c.cmp
	}
}
