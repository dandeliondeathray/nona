"""
Nonastaging is a WebSocket based interface to the Nona system. Its purpose is to simplify
acceptance testing of the Nona system deployed in a staging environment.
"""

import asyncio
import os
import threading
import json
from nonainterface import NonaInterface, AvroSchemas
import websockets


class NonaStaging:
    """NonaStaging implements the entire micro service."""
    def __init__(self, loop, bootstrap_servers, schema_path):
        self.loop = loop
        schemas = AvroSchemas(schema_path)
        self._nona = NonaInterface("staging", bootstrap_servers, schemas)
        self._clients = []
        self._chat_event_reader = None
        self._stop_semaphore = threading.Semaphore()

    def read_chat_events(self):
        while True:
            chat_event = self._nona.chat_events.get(block=True)
            print("Read chat event from nona:", chat_event)
            if chat_event is None:
                break
            self.loop.call_soon_threadsafe(asyncio.async, self.send_chat_event(chat_event))

    def wait_for_stop(self):
        self._stop_semaphore.acquire()

    def stop(self):
        self._nona.stop()
        self.stop_chat_event_reader()
        self.loop.stop()
        self._stop_semaphore.release()

    def stop_chat_event_reader(self):
        self._nona.chat_events.put(None)
        self._chat_event_reader.join()

    @asyncio.coroutine
    def send_chat_event(self, chat_event):
        clients = self._clients[:]
        chat_event_json = {'user_id': chat_event.user_id,
                           'team': chat_event.team,
                           'text': chat_event.text}
        for client in clients:
            yield from client.send(json.dumps(chat_event_json))

    @asyncio.coroutine
    def handle_command(self, msg):
        command = json.loads(msg)
        if command['name'] == 'user_requests_puzzle':
            user_id = command['user_id']
            self._nona.user_requests_puzzle(user_id)
        else:
            print("Invalid command:", command)

    @asyncio.coroutine
    def add_client(self, websocket, _path):
        print("New connection")
        self._clients.append(websocket)
        while True:
            try:
                msg = yield from websocket.recv()
            except websockets.exceptions.ConnectionClosed:
                print("Connection closed.")
                return
            print("Read command:", msg)
            yield from self.handle_command(msg)

    def run_forever(self):
        """Start the event loop and only return when it is stopped."""
        self._nona.start()
        self._chat_event_reader = threading.Thread(target=self.read_chat_events)
        self._chat_event_reader.start()

        start_server = websockets.serve(self.add_client, '0.0.0.0', 8765)
        self.loop.run_until_complete(start_server)
        self.loop.run_forever()

class NonaStagingService:
    def __init__(self, brokers, schema_path):
        self._brokers = brokers
        self._schema_path = schema_path
        self.loop = None
        self._nonastaging = None

    def stop(self):
        self.loop.call_soon_threadsafe(self._nonastaging.stop)
        self._nonastaging.wait_for_stop()

    def start(self):
        self.loop = asyncio.get_event_loop()
        self._nonastaging = NonaStaging(self.loop, self._brokers, self._schema_path)
        self._nonastaging.run_forever()

if __name__ == "__main__":
 #   import pymetamorph.metamorph
#    m = pymetamorph.metamorph.Metamorph()
#    m.connect()
#    m.request_kafka_reset(['nona_UserRequestsPuzzle'])
#    m.await_reset_complete()
#    print("Reset complete")

    my_brokers = os.environ['KAFKA_BROKERS']
    my_schema_path = os.environ['SCHEMA_PATH']
    service = NonaStagingService(my_brokers, my_schema_path)
    service.start()
