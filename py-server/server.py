#!/usr/bin/env python
# -*- coding: utf-8 -*-

import ctypes
import socket
import time
from pymodbus.client.sync import ModbusSerialClient
from threading import Thread
from socketserver import ThreadingMixIn

runCond = True


class ValuePusher(Thread):
    def __init__(self, worker):
        Thread.__init__(self)
        self.__worker = worker
        self.__connections = []

    def newConn(self, conn):
        self.__connections.append(conn)

    def run(self):
        while runCond:
            value = self.__worker()
            connsToDelete = []
            for conn in self.__connections:
                try:
                    conn.send(value)
                except BrokenPipeError:
                    connsToDelete.append(conn)
            for conn in connsToDelete:
                self.__connections.remove(conn)
            time.sleep(1)
        for conn in self.__connections:
            conn.close()


class ModbusReader(object):
    def __init__(self, unitId):
        self.__unitId = unitId
        self.__client = ModbusSerialClient(
            method='rtu', port='/dev/ttyUSB0',
            baudrate=19200, parity='E', stopbits=1, bytesize=8)
        self.__client.connect()

    def __call__(self):
        response = self.__client.read_holding_registers(
            address=0x5B00, count=66, unit=self.__unitId)
        activePowerTotal = ctypes.c_int32(
            (response.registers[20] << 16) | response.registers[21]).value / 100
        return (str(activePowerTotal) + " W\n").encode("UTF-8")

    def __del__(self):
        self.__client.close()


def main():
    global runCond

    TCP_IP = '0.0.0.0'
    TCP_PORT = 2004
    TCP_BACKLOG = 4

    tcpServer = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcpServer.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    tcpServer.bind((TCP_IP, TCP_PORT))

    valuePusher = ValuePusher(ModbusReader(1))
    valuePusher.start()

    threads = []
    threads.append(valuePusher)

    try:
        while runCond:
            tcpServer.listen(TCP_BACKLOG)
            print("Waiting for connections from TCP clients...")
            (conn, (ip, port)) = tcpServer.accept()
            valuePusher.newConn(conn)
    except KeyboardInterrupt:
        runCond = False

    for t in threads:
        t.join()


if __name__ == "__main__":
    main()
