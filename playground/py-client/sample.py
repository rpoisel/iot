#!/usr/bin/env python
# -*- coding: utf-8 -*-

from pymodbus.client.sync import ModbusSerialClient
import ctypes
import sys
import time


def main():
    client = ModbusSerialClient(
        method='rtu', port='/dev/ttyUSB0',
        baudrate=19200, parity='E', stopbits=1, bytesize=8)
    client.connect()
    for u in range(1, len(sys.argv)):
        # unit ... address of slave on bus (default value: 0)
        response = client.read_holding_registers(
            address=0x5B00, count=66, unit=int(sys.argv[u]))
        voltage = ((response.registers[0] << 16) | response.registers[1]) / 10
        activePowerTotal = (
            (response.registers[20] << 16) | response.registers[21])
        print("Voltage L1-N: " + str(voltage) + "V")
        print("Active Power Total: "
              + str(ctypes.c_int32(activePowerTotal).value / 100)
              + "W (" + hex(activePowerTotal) + ")")
        response = client.read_holding_registers(
            address=0x5000, count=4, unit=int(sys.argv[1]))
        activeImport = (
            (response.registers[0] << 32)
            | (response.registers[1] << 24)
            | (response.registers[2] << 16)
            | response.registers[3]) / 100
        print("Imported: " + str(activeImport) + "kWh")
        time.sleep(1)
    client.close()


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: " + sys.argv[0] + " <units>")
        sys.exit(1)
    main()
