package i2c

func SetBit(data *byte, idx uint8, state bool) {
	if state {
		(*data) |= byte(0x01 << idx)
	} else {
		(*data) &= ^byte(0x01 << idx)
	}
}
