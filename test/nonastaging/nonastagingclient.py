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
        self._event_buffer = []

    def start(self):
        self._ws = self._loop.run_until_complete(websockets.connect(self._address, loop=self._loop))

    def stop(self):
        self._loop.call_soon_threadsafe(self._loop.stop)

    def user_requests_puzzle(self, user_id):
        print("NonaStagingClient: User requests puzzle for user", user_id)
        command = {'name': 'user_requests_puzzle', 'user_id': user_id}
        self._loop.run_until_complete(self._send_command(command))

    def await_chat(self, user_id, team, text):
        self._loop.run_until_complete(self._await_chat(user_id, team, text))

    @asyncio.coroutine
    def _await_chat(self, user_id, team, text):
        try:
            message = yield from asyncio.wait_for(self._find_chat_event(user_id, team, text),
                                                  timeout=2.0,
                                                  loop=self._loop)
            return message
        except asyncio.TimeoutError:
            raise NoSuchChatMessage(user_id, team, text)

    @asyncio.coroutine
    def _find_chat_event(self, user_id, team, text):
        for i, chat in enumerate(self._event_buffer):
            if chat.user_id == user_id and chat.team == team and chat.text == text:
                self._event_buffer = self._event_buffer[:i] + self._event_buffer[i+1:]
                return chat
        while True:
            chat_event = yield from self._ws.recv()
            print("Received chat event in client:", chat_event)
            chat = json.loads(chat_event)
            if chat.user_id == user_id and chat.team == team and chat.text == text:
                return chat
            self._event_buffer.append(chat)

    @asyncio.coroutine
    def _send_command(self, command):
        command_json = json.dumps(command)
        yield from self._ws.send(command_json)
