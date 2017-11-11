import asyncio
import websockets
import json
import threading


class NoSuchChatMessage(Exception):
    def __init__(self, user_id, team, text):
        self.user_id = user_id
        self.team = team
        self.text = text

    def __str__(self):
        return "No such chat message: user_id {}, team {}, text '{}'".format(
            self.user_id, self.team, self.text)

class NonaStagingClient:
    def __init__(self, address, loop=None):
        self._address = address
        self._loop = loop if loop is not None else asyncio.new_event_loop()
        self._ws = None
        self._receive_queue = asyncio.Queue(loop=self._loop)
        self._event_buffer = []
        self._thread = None

    def start(self):
        self._thread = threading.Thread(target=self._connect_and_run)
        self._thread.start()

    def stop(self):
        self._loop.call_soon_threadsafe(self._loop.stop)
        self._thread.join()

    def user_requests_puzzle(self, user_id):
        print("NonaStagingClient: User requests puzzle for user", user_id)
        command = {'name': 'user_requests_puzzle', 'user_id': user_id}
        self._loop.call_soon_threadsafe(self._send_command, command)

    def await_chat(self, user_id, team, text):
        try:
            asyncio.wait_for(self._find_chat_event(user_id, team, text),
                             timeout=2.0,
                             loop=self._loop)
        except asyncio.TimeoutError:
            raise NoSuchChatMessage(user_id, team, text)

    def _connect_and_run(self):
        self._ws = self._loop.run_until_complete(websockets.connect(self._address))
        self._loop.run_forever()

    @asyncio.coroutine
    def _find_chat_event(self, user_id, team, text):
        for i, chat in enumerate(self._event_buffer):
            if chat.user_id == user_id and chat.team == team and chat.text == text:
                self._event_buffer = self._event_buffer[:i] + self._event_buffer[i+1:]
                return chat
        while True:
            chat = yield from self._receive_queue.get()
            if chat.user_id == user_id and chat.team == team and chat.text == text:
                return chat
            self._event_buffer.append(chat)

    @asyncio.coroutine
    def _send_command(self, command):
        command_json = json.dumps(command)
        yield from self._ws.send(command_json)

    @asyncio.coroutine
    def _recv_chat_events(self):
        while True:
            chat_event = yield from self._ws.recv()
            decoded = json.loads(chat_event)
            yield from self._receive_queue.put(decoded)
