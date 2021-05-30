package i2c_test

import (
	"testing"

	"github.com/rpoisel/iot/internal/io/i2c"
	"github.com/stretchr/testify/suite"
)

type I2CTestSuite struct {
	suite.Suite
}

func (i *I2CTestSuite) TestBitSets() {
	val := byte(0xff)
	i2c.SetBit(&val, 0, false)

	i.Equal(byte(0xfe), val)
}

func TestLightTestSuite(t *testing.T) {
	suite.Run(t, new(I2CTestSuite))
}
