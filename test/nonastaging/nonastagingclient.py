import asyncio
import websockets
import json
import threading


class NoSuchChatMessage(Exception):
    def __init__(self, chat_matcher):
        self._matcher = chat_matcher

    def __str__(self):
        return "No such chat message matching: {}".format(self._matcher)


class Matcher:
    def matches(self, value):
        raise NotImplementedError("Matcher.matches()")


class AnyMatcher(Matcher):
    def matches(self, _value):
        return True

    def __str__(self):
        return "<any value>"


class ValueMatcher(Matcher):
    def __init__(self, value):
        self._value = value

    def matches(self, value):
        return self._value == value

    def __str__(self):
        return str(self._value)


class ChatMatcher:
    """Matches against chat messages from Nona to find the one expected by a test."""
    def __init__(self, user_id, team, text_matcher):
        self._user_id = user_id
        self._team = team
        self._text_matcher = text_matcher

    def matches(self, chat_event):
        """Does the chat event match the user id, team, and text?

        Returns True if so, False otherwise.abs
        """
        m = self._user_id == chat_event['user_id'] and self._team == chat_event['team']
        return m and self._text_matcher.matches(chat_event['text'])

    def __str__(self):
        return "ChatMessage(user_id = {}, team = {}, text = {})".format(self._user_id,
                                                                        self._team,
                                                                        self._text_matcher)


class NonaStagingClient:
    def __init__(self, address, loop=None):
        self._address = address
        self._loop = loop if loop is not None else asyncio.new_event_loop()
        self._ws = None
        self._event_buffer = []

    def start(self):
        self._ws = self._loop.run_until_complete(websockets.connect(self._address, loop=self._loop))

    def stop(self):
        self._loop.run_until_complete(self._ws.close())

    def user_requests_puzzle(self, user_id):
        print("NonaStagingClient: User requests puzzle for user", user_id)
        command = {'name': 'user_requests_puzzle', 'user_id': user_id}
        self._loop.run_until_complete(self._send_command(command))

    def await_chat(self, user_id, team, text=None):
        if text is None:
            text_matcher = AnyMatcher()
        elif isinstance(text, Matcher):
            text_matcher = text
        else:
            text_matcher = ValueMatcher(text)
        chat_matcher = ChatMatcher(user_id, team, text_matcher)
        return self._loop.run_until_complete(self._await_chat(chat_matcher))

    @asyncio.coroutine
    def _await_chat(self, chat_matcher):
        try:
            message = yield from asyncio.wait_for(self._find_chat_event(chat_matcher),
                                                  timeout=2.0,
                                                  loop=self._loop)
            return message
        except asyncio.TimeoutError:
            raise NoSuchChatMessage(chat_matcher)

    @asyncio.coroutine
    def _find_chat_event(self, chat_matcher):
        for i, chat in enumerate(self._event_buffer):
            if chat_matcher.matches(chat):
                self._event_buffer = self._event_buffer[:i] + self._event_buffer[i+1:]
                return chat
        while True:
            chat_event = yield from self._ws.recv()
            print("Received chat event in client:", chat_event)
            chat = json.loads(chat_event)
            if chat_matcher.matches(chat):
                return chat
            self._event_buffer.append(chat)

    @asyncio.coroutine
    def _send_command(self, command):
        command_json = json.dumps(command)
        yield from self._ws.send(command_json)
