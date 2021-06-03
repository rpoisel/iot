package logic

import "time"

type MultiClick struct {
	d   time.Duration
	cnt uint
	In  <-chan interface{}
	Out chan<- uint
}

func NewMultiClick(d time.Duration) *MultiClick {
	return &MultiClick{
		d: d,
	}
}

func (m *MultiClick) Process() {
	for {
		<-m.In
		m.cnt = 1
		timer := time.NewTimer(m.d)
	countLoop:
		for {
			select {
			case <-m.In:
				m.cnt++
			case <-timer.C:
				break countLoop
			}
		}
		m.Out <- m.cnt
		timer.Stop()
	}
}
