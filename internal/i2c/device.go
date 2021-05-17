package i2c

type I2CInputDevice interface {
	// Read is called by the containing I2CBus
	// in order to transfer the hardware state
	// to the representation
	Read()
	// NumInputs is mainly used to perform configuration validation
	NumInputs() uint16
}

type I2COutputDevice interface {
	// Write is called by the containing I2CBus
	// in order to transfer the current state of this
	// device to the hardware
	Write()
	// NumOutputs is mainly used to perform configuration validation
	NumOutputs() uint16
}
