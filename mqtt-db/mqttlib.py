#!/usr/bin/env python
# -*- coding: utf-8 -*-

import paho.mqtt.client as mqtt
import struct


class MeasurementBroker(mqtt.Client):
    TOPIC = '/homeautomation/power/cumulative'

    def __init__(self, handler):
        super().__init__()
        self.handler = handler

    def on_connect(self, client, user_data, flags, rc):
        print("Connected with rc = " + str(rc))

    def on_message(self, client, user_data, msg):
        if msg.topic != MeasurementBroker.TOPIC:
            return

        solar = struct.unpack("<i", msg.payload[0:4])[0]
        total = struct.unpack("<i", msg.payload[8:12])[0]

        self.handler.on_measurement(solar, total)

    def run(self):
        self.tls_set()
        self.username_pw_set(username='abc', password='xyz')
        self.connect("hostname.tld", 8883, 60)

        self.subscribe(MeasurementBroker.TOPIC, 0)
        rc = 0
        while rc == 0:
            rc = self.loop()
        return rc
