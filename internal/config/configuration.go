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

type PhysicalInput struct {
	Name   string
	Modify string
}

type Configuration struct {
	Version    string
	I2C        map[int]I2CBus
	Automation struct {
		Blinds []struct {
			Name   string
			Input1 string
			Input2 string
			Output string
		}
		Lights []struct {
			Name   string
			Inputs struct {
				Physical PhysicalInput
				MQTT     string
			}
		}
		MQTT []struct {
			Topic  string
			Inputs struct {
				Physical PhysicalInput
			}
		}
	}
}
