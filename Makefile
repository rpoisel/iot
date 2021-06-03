ALL_BINARIES := cmd/modbus-mqtt/modbus-mqtt \
	cmd/joyblind/joyblind \
	cmd/loxone-proxy/loxone-proxy \
	cmd/mqtt-db-postgres/mqtt-db-postgres \
	cmd/i2c-eg/i2c-eg \
	cmd/homeautomation/homeautomation \
	cmd/flow/flow
GO ?= /usr/lib/go-1.16/bin/go

all: $(ALL_BINARIES)

.PHONY: clean \
	version \
	mod_update \
	test \
	$(ALL_BINARIES)

cmd/modbus-mqtt/modbus-mqtt:
	(cd cmd/modbus-mqtt && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

cmd/joyblind/joyblind:
	(cd cmd/joyblind && $(GO) build -v -mod=vendor)

cmd/loxone-proxy/loxone-proxy:
	(cd cmd/loxone-proxy && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

cmd/mqtt-db-postgres/mqtt-db-postgres:
	(cd cmd/mqtt-db-postgres && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

cmd/i2c-eg/i2c-eg:
	(cd cmd/i2c-eg && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

cmd/homeautomation/homeautomation:
	(cd cmd/homeautomation && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

cmd/flow/flow:
	(cd cmd/flow && env GOOS=linux GOARCH=arm GOARM=7 $(GO) build -v -mod=vendor)

clean:
	-rm $(ALL_BINARIES)

version:
	$(GO) version

mod_update:
	$(GO) get -u all

test:
	go test -mod=vendor ./...
