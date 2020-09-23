all: modbus-mqtt joyblind/joyblind

modbus-mqtt: main.go go.mod
	env GOOS=linux GOARCH=arm GOARM=7 go build

joyblind/joyblind: joyblind/joyblind.go
	(cd joyblind && go build)

.PHONY: clean
clean:
	-rm modbus-mqtt joyblind/joyblind
