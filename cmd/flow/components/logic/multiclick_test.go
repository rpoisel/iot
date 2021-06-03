package logic_test

import (
	"sync"
	"testing"
	"time"

	"github.com/rpoisel/iot/cmd/flow/components/logic"
	"github.com/stretchr/testify/suite"
)

type MultiClickTestSuite struct {
	suite.Suite
}

func (m *MultiClickTestSuite) TestSingleClick() {
	mc := logic.NewMultiClick(100 * time.Microsecond)
	in := make(chan interface{})
	out := make(chan uint)
	mc.In = in
	mc.Out = out
	go mc.Process()
	in <- true
	m.Equal(uint(1), <-out)
}

func (m *MultiClickTestSuite) TestDoubleClick() {
	mc := logic.NewMultiClick(100 * time.Microsecond)
	in := make(chan interface{})
	out := make(chan uint)
	mc.In = in
	mc.Out = out
	go mc.Process()
	in <- true
	in <- true
	m.Equal(uint(2), <-out)
}

func (m *MultiClickTestSuite) TestSingleThenDoubleClick() {
	mc := logic.NewMultiClick(100 * time.Microsecond)
	in := make(chan interface{})
	out := make(chan uint)
	mc.In = in
	mc.Out = out
	go mc.Process()
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		m.Equal(uint(1), <-out)
		m.Equal(uint(2), <-out)
	}()
	in <- true
	time.Sleep(150 * time.Millisecond)
	in <- true
	in <- true
	wg.Wait()
}

func (m *MultiClickTestSuite) TestDoubleThenSingleClick() {
	mc := logic.NewMultiClick(100 * time.Microsecond)
	in := make(chan interface{})
	out := make(chan uint)
	mc.In = in
	mc.Out = out
	go mc.Process()
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		m.Equal(uint(2), <-out)
		m.Equal(uint(1), <-out)
	}()
	in <- true
	in <- true
	time.Sleep(150 * time.Millisecond)
	in <- true
	wg.Wait()
}

func TestMultiClickTestSuite(t *testing.T) {
	suite.Run(t, new(MultiClickTestSuite))
}
