#!/usr/bin/env python
# -*- coding: utf-8 -*-

import asyncio
import paho.mqtt.client as mqtt
import socket
import websockets


class AsyncioHelper:
    def __init__(self, loop, client):
        self.loop = loop
        self.client = client
        self.client.on_socket_open = self.on_socket_open
        self.client.on_socket_close = self.on_socket_close
        self.client.on_socket_register_write = self.on_socket_register_write
        self.client.on_socket_unregister_write = self.on_socket_unregister_write

    def on_socket_open(self, client, userdata, sock):
        print("Socket opened")

        def cb():
            print("Socket is readable, calling loop_read")
            client.loop_read()

        self.loop.add_reader(sock, cb)
        self.misc = self.loop.create_task(self.misc_loop())

    def on_socket_close(self, client, userdata, sock):
        print("Socket closed")
        self.loop.remove_reader(sock)
        self.misc.cancel()

    def on_socket_register_write(self, client, userdata, sock):
        print("Watching socket for writability.")

        def cb():
            print("Socket is writable, calling loop_write")
            client.loop_write()

        self.loop.add_writer(sock, cb)

    def on_socket_unregister_write(self, client, userdata, sock):
        print("Stop watching socket for writability.")
        self.loop.remove_writer(sock)

    async def misc_loop(self):
        print("misc_loop started")
        while self.client.loop_misc() == mqtt.MQTT_ERR_SUCCESS:
            try:
                await asyncio.sleep(1)
            except asyncio.CancelledError:
                break
        print("misc_loop finished")


class AsyncMqttClient:
    def __init__(self, loop):
        self.loop = loop
        self.got_message = self.loop.create_future()

        self.client = mqtt.Client()
        self.client.on_connect = self.on_connect
        self.client.on_message = self.on_message
        self.client.tls_set()
        self.client.username_pw_set(
            username='abc', password='xyz')

        AsyncioHelper(asyncio.get_event_loop(), self.client)
        self.client.connect("hostname.tld", 8883, 60)
        self.client.socket().setsockopt(socket.SOL_SOCKET,
                                        socket.SO_SNDBUF,
                                        2048)

    def on_connect(self, client, userdata, flags, rc):
        client.subscribe('/homeautomation/power/active')
        client.subscribe('/homeautomation/power/solar')

    def on_message(self, client, userdata, msg):
        self.got_message.set_result(msg)


async def mqtt_subscribe(websocket, path):
    client = AsyncMqttClient(asyncio.get_event_loop())
    while True:
        msg = await client.got_message
        client.got_message = asyncio.get_event_loop().create_future()
        try:
            await websocket.send(
                msg.topic + " = " + msg.payload.decode('UTF-8'))
        except websockets.exceptions.ConnectionClosedError:
            return


def main():
    start_server = websockets.serve(mqtt_subscribe, "127.0.0.1", 5678)

    asyncio.get_event_loop().run_until_complete(start_server)
    asyncio.get_event_loop().run_forever()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass
