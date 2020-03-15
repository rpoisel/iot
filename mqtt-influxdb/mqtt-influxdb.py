#!/usr/bin/env python
# -*- coding: utf-8 -*-

from influxdb import InfluxDBClient
import paho.mqtt.client as mqtt
import struct


class ValueWriter(mqtt.Client):
    def __init__(self, influx):
        super().__init__()
        self.influx = influx

    def on_connect(self, client, user_data, flags, rc):
        print("Connected with rc = " + str(rc))

    def on_message(self, client, user_data, msg):
        solar = struct.unpack("<i", msg.payload[0:4])[0]
        total = struct.unpack("<i", msg.payload[8:12])[0]

        dbData = [
            {
                "measurement": "power",
                "fields": {
                    "solar": solar,
                    "total": total,
                }
            }
        ]
        self.influx.write_points(dbData)

    def run(self):
        self.tls_set()
        self.username_pw_set(username='abc', password='xyz')
        self.connect("hostname.tld", 8883, 60)

        self.subscribe('/homeautomation/power/cumulative', 0)
        rc = 0
        while rc == 0:
            rc = self.loop()
        return rc


def main():
    influx = InfluxDBClient(host='localhost', port=8086, database='power')
    valueWriter = ValueWriter(influx)
    valueWriter.run()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass
