#!/usr/bin/env python
# -*- coding: utf-8 -*-

from pymodbus.client.sync import ModbusSerialClient
import sys


def main():
    client = ModbusSerialClient(
        method='rtu', port='/dev/ttyUSB0',
        baudrate=19200, parity='E', stopbits=1, bytesize=8)
    client.connect()
    # unit ... address of slave on bus (default value: 0)
    response = client.read_holding_registers(
        address=0x5B00, count=66, unit=int(sys.argv[1]))
    voltage = ((response.registers[0] << 16) | response.registers[1]) / 10
    print("Voltage L1-N: " + str(voltage) + "V")
    response = client.read_holding_registers(
        address=0x5000, count=4, unit=int(sys.argv[1]))
    activeImport = (
        (response.registers[0] << 32)
        | (response.registers[1] << 24)
        | (response.registers[2] << 16)
        | response.registers[3]) / 100
    print("Imported: " + str(activeImport) + "kWh")
    client.close()


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: " + sys.argv[0] + " <unit>")
        sys.exit(1)
    main()
