all: cmd/modbus-mqtt/modbus-mqtt cmd/joyblind/joyblind cmd/loxone-proxy/loxone-proxy cmd/mqtt-db-postgres/mqtt-db-postgres

cmd/modbus-mqtt/modbus-mqtt: cmd/modbus-mqtt/main.go
	(cd cmd/modbus-mqtt && env GOOS=linux GOARCH=arm GOARM=7 go build)

cmd/joyblind/joyblind: cmd/joyblind/main.go
	(cd cmd/joyblind && go build)

cmd/loxone-proxy/loxone-proxy: cmd/loxone-proxy/main.go
	(cd cmd/loxone-proxy && env GOOS=linux GOARCH=arm GOARM=7 go build)

cmd/mqtt-db-postgres/mqtt-db-postgres: cmd/mqtt-db-postgres/main.go
	(cd cmd/mqtt-db-postgres && env GOOS=linux GOARCH=arm GOARM=7 go build)

.PHONY: clean
clean:
	-rm \
		cmd/modbus-mqtt/modbus-mqtt \
		cmd/joyblind/joyblind \
		cmd/loxone-proxy/loxone-proxy \
		cmd/mqtt-db-postgres/mqtt-db-postgres
