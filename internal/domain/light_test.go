package domain_test

import (
	"testing"

	"github.com/rpoisel/iot/internal/domain"
	"github.com/stretchr/testify/suite"
)

type LightTestSuite struct {
	suite.Suite
}

func (l *LightTestSuite) TestSwitchCycle() {
	inputCh := make(chan bool)
	outputCh := make(chan bool)
	light := domain.Light{
		Input:  inputCh,
		Output: outputCh,
	}
	go light.Run()

	inputCh <- true
	l.True(<-outputCh)
	inputCh <- false
	l.False(<-outputCh)
}

func TestLightTestSuite(t *testing.T) {
	suite.Run(t, new(LightTestSuite))
}
