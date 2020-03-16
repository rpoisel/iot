#!/usr/bin/env python
# -*- coding: utf-8 -*-

from mqttlib import MeasurementBroker
import psycopg2


class PostgresWriter(object):
    def __init__(self, connect_arg):
        self.__conn = psycopg2.connect(connect_arg)

    def on_measurement(self, solar, total):
        cur = self.__conn.cursor()
        cur.execute(
            "INSERT INTO public.power (solar, total) VALUES (%s, %s)", (solar, total))
        self.__conn.commit()
        cur.close()

    def __del__(self):
        self.__conn.close()


def main():
    postgres = PostgresWriter(
        "dbname=power user=power_rw password=xxx host=localhost")
    measurementBroker = MeasurementBroker(postgres)
    measurementBroker.run()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass
