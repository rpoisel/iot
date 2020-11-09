FROM golang:1.15 as builder
WORKDIR /go/src/IoT
COPY . .
WORKDIR /go/src/IoT/cmd/mqtt-db-postgres
RUN go build -o mqtt-db-postgres


FROM debian:bullseye-slim
WORKDIR /
COPY --from=builder /go/src/IoT/cmd/mqtt-db-postgres/mqtt-db-postgres .

CMD ["/mqtt-db-postgres"]
