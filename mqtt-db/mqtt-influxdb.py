#!/usr/bin/env python
# -*- coding: utf-8 -*-

from influxdb import InfluxDBClient
from mqttlib import MeasurementBroker


class InfluxWriter(object):
    def __init__(self, **kwargs):
        self.__influx = InfluxDBClient(**kwargs)

    def on_measurement(self, solar, total):
        dbData = [
            {
                "measurement": "power",
                "fields": {
                    "solar": solar,
                    "total": total,
                }
            }
        ]
        self.__influx.write_points(dbData)


def main():
    influx = InfluxWriter(host='localhost', port=8086, database='power')
    measurementBroker = MeasurementBroker(influx)
    measurementBroker.run()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass
