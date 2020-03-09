#!/usr/bin/env python
# -*- coding: utf-8 -*-

from influxdb import InfluxDBClient
import paho.mqtt.client as mqtt


class ValueWriter(mqtt.Client):
    DEFAULT_PATH = '/homeautomation/power/'

    def __init__(self, influx):
        super().__init__()
        self.influx = influx

    def on_connect(self, client, user_data, flags, rc):
        print("Connected with rc = " + str(rc))

    def on_message(self, client, user_data, msg):
        source = msg.topic.replace(ValueWriter.DEFAULT_PATH, '')
        dbData = [
            {
                "measurement": "power",
                "tags": {
                    "source": source,
                },
                "fields": {
                    "value": float(msg.payload)
                }
            }
        ]
        self.influx.write_points(dbData)

    def run(self, topics):
        self.tls_set()
        self.username_pw_set(username='abc', password='xyz')
        self.connect("hostname.tld", 8883, 60)

        for topic in topics:
            self.subscribe(ValueWriter.DEFAULT_PATH + topic, 0)
        rc = 0
        while rc == 0:
            rc = self.loop()
        return rc


def main():
    influx = InfluxDBClient(host='localhost', port=8086, database='power')
    valueWriter = ValueWriter(influx)
    valueWriter.run(['obtained', 'solar', 'total'])


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass
