#!/usr/bin/env python
# -*- coding: utf-8 -*-

from paho.mqtt import client as mqtt
from pymodbus.client.sync import ModbusSerialClient
from threading import Thread
import ctypes
import time

quitCond = False
valuePusher = None


class ValuePusher(Thread):
    def __init__(self, connection, units):
        Thread.__init__(self)
        self.__connection = connection
        self.__units = units
        self.__client = None

    def __maintainModbusConnection(self):
        if self.__client == None:
            self.__client = ModbusSerialClient(
                method='rtu', port='/dev/ttyUSB0', baudrate=19200, parity='E', stopbits=1, bytesize=8)
            self.__client.connect()

    def __readFromModbus(self):
        result = {}
        for u in self.__units:
            try:
                response = self.__client.read_holding_registers(
                    address=0x5B00, count=66, unit=u)
                power = ctypes.c_int32(
                    (response.registers[20] << 16) | response.registers[21]).value / 100

                result[self.__units[u]] = int(power)
            except AttributeError:
                self.__client.close()
                self.__client = None
                time.sleep(2)
            time.sleep(.1)
        return result

    def __publishToMqtt(self, values):
        for v in values:
            powerStr = str(values[v])
            print(v + ": " + powerStr)
            self.__connection.publish("/homeautomation/power/" + v, powerStr)

    def run(self):
        while not quitCond:
            self.__maintainModbusConnection()
            values = self.__readFromModbus()

            values['total'] = (values['solar'] if values['solar'] > 0 else 0) \
                + values['obtained']

            self.__publishToMqtt(values)
            time.sleep(1)

    def __del__(self):
        self.__client.close()


def on_connect(client, userdata, flags, rc):
    global valuePusher
    print("Connected with result code "+str(rc))
    valuePusher = ValuePusher(client, {1: 'solar', 2: 'obtained'})
    valuePusher.start()


def main():
    global quitCond
    client = mqtt.Client()

    client.on_connect = on_connect
    client.tls_set()
    client.username_pw_set(username='abc', password='xyz')

    client.connect("hostname.tld", 8883, 60)
    try:
        client.loop_forever()
    except KeyboardInterrupt:
        quitCond = True


if __name__ == "__main__":
    main()
