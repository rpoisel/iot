all: modbus-mqtt joyblind/joyblind loxone-proxy/loxone-proxy

modbus-mqtt: main.go go.mod
	env GOOS=linux GOARCH=arm GOARM=7 go build

joyblind/joyblind: joyblind/joyblind.go
	(cd joyblind && go build)

loxone-proxy/loxone-proxy: loxone-proxy/main.go
	(cd loxone-proxy && env GOOS=linux GOARCH=arm GOARM=7 go build)

.PHONY: clean
clean:
	-rm modbus-mqtt joyblind/joyblind loxone-proxy/loxone-proxy
