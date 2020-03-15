#!/usr/bin/env python
# -*- coding: utf-8 -*-

from paho.mqtt import client as mqtt
from pymodbus.client.sync import ModbusSerialClient
from threading import Thread
import ctypes
import struct
import time

quitCond = False
valuePusher = None


class ValueRetriever(object):
    def __init__(self, units):
        super().__init__()
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

    def next(self):
        self.__maintainModbusConnection()
        return self.__readFromModbus()

    def __del__(self):
        if self.__client is not None:
            self.__client.close()


class ValuePusher(Thread):
    def __init__(self, client):
        super().__init__()
        self.__client = client
        self.valueRetriever = ValueRetriever({1: 'solar', 2: 'obtained'})

    def __publishToMqtt(self, values):
        for v in values:
            pubVal = values[v]
            if v in ['solar', 'obtained', 'total']:
                pubVal = str(values[v])
                print(v + ": " + pubVal)
            self.__client.publish("/homeautomation/power/" + v, pubVal)

    def run(self):
        while not quitCond:
            values = self.valueRetriever.next()
            values['total'] = (values['solar'] if values['solar']
                               > 0 else 0) + values['obtained']
            values['cumulative'] = struct.pack('<iii',
                                               values['solar'],
                                               values['obtained'],
                                               values['total'])

            self.__publishToMqtt(values)

            time.sleep(1)


def on_connect(client, userdata, flags, rc):
    global valuePusher
    print("Connected with result code "+str(rc))
    valuePusher = ValuePusher(client)
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
