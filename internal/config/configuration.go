package config

type I2CDevice struct {
	Name    string
	Address uint8
	Chip    string
}

type I2CBus struct {
	In  []I2CDevice
	Out []I2CDevice
}

type LocalInput struct {
	Path   string
	Modify string
}

type Automation struct {
	Blinds []struct {
		Name    string
		Input1  LocalInput
		Input2  LocalInput
		Output1 string
		Output2 string
		MQTT    string
	}
	Lights []struct {
		Name   string
		Input  LocalInput
		Output string
		MQTT   string
	}
	MQTT []struct {
		Topic   string
		Message string
		Input   LocalInput
	}
}

type I2C map[int]I2CBus

type IO struct {
	I2C I2C
}

type Configuration struct {
	Version    string
	Automation Automation
	IO         IO
}
