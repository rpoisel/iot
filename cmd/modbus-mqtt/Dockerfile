FROM golang:1.15 as builder
WORKDIR /go/src/IoT
COPY . .
WORKDIR /go/src/IoT/cmd/modbus-mqtt
RUN go build -o modbus-mqtt


FROM debian:bullseye-slim
WORKDIR /
COPY --from=builder /go/src/IoT/cmd/modbus-mqtt/modbus-mqtt .
RUN apt-get update && apt-get install -yy ca-certificates

CMD ["/modbus-mqtt"]

