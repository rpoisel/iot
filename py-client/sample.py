#!/usr/bin/env python
# -*- coding: utf-8 -*-

from pymodbus.client.sync import ModbusSerialClient


def main():
    client = ModbusSerialClient(
        method='rtu', port='/dev/ttyUSB0',
        baudrate=19200, parity='N', stopbits=1, bytesize=8)
    client.connect()
    response = client.read_holding_registers(0x5B00, 66, unit=1)
    voltage = ((response.registers[0] << 16) | response.registers[1]) / 10
    print("Voltage L1-N: " + str(voltage) + "V")
    client.close()


if __name__ == "__main__":
    main()
