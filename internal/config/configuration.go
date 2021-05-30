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
	Name   string
	Modify string
}

type Automation struct {
	Blinds []struct {
		Name    string
		Input1  string
		Input2  string
		Output1 string
		Output2 string
	}
	Lights []struct {
		Name   string
		Inputs struct {
			Local LocalInput
			MQTT  string
		}
		Output string
	}
	MQTT []struct {
		Topic  string
		Inputs struct {
			Local LocalInput
		}
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
