all: cmd/modbus-mqtt/modbus-mqtt \
	cmd/joyblind/joyblind \
	cmd/loxone-proxy/loxone-proxy \
	cmd/mqtt-db-postgres/mqtt-db-postgres \
	cmd/i2c-eg/i2c-eg

cmd/modbus-mqtt/modbus-mqtt: cmd/modbus-mqtt/main.go internal/util/util.go
	(cd cmd/modbus-mqtt && env GOOS=linux GOARCH=arm GOARM=7 go build -v -mod=vendor)

cmd/joyblind/joyblind: cmd/joyblind/main.go internal/util/util.go
	(cd cmd/joyblind && go build -v -mod=vendor)

cmd/loxone-proxy/loxone-proxy: cmd/loxone-proxy/main.go internal/util/util.go
	(cd cmd/loxone-proxy && env GOOS=linux GOARCH=arm GOARM=7 go build -v -mod=vendor)

cmd/mqtt-db-postgres/mqtt-db-postgres: cmd/mqtt-db-postgres/main.go internal/util/util.go
	(cd cmd/mqtt-db-postgres && env GOOS=linux GOARCH=arm GOARM=7 go build -v -mod=vendor)

cmd/i2c-eg/i2c-eg: cmd/i2c-eg/main.go
	(cd cmd/i2c-eg && env GOOS=linux GOARCH=arm GOARM=7 go build -v -mod=vendor)

.PHONY: clean
clean:
	-rm \
		cmd/modbus-mqtt/modbus-mqtt \
		cmd/joyblind/joyblind \
		cmd/loxone-proxy/loxone-proxy \
		cmd/mqtt-db-postgres/mqtt-db-postgres \
		cmd/i2c-eg/i2c-eg
