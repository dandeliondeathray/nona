"""
nonainterface.testing contains code for helping with the testing of the Nona system through nonainterface, as well as
nonainterface itself.
"""

import queue
from nonainterface import ChatMessage


class ChatQueue:
    """
    A testing interface against the arriving chat events from Nonainterface.
    """
    def __init__(self, queue):
        self._queue = queue
        self._got_from_queue = []

    def has(self, fields, timeout=2.0):
        """
        Check if the chat queue has a message matching specific fields.

        :param fields: A dict of all fields and their expected value.
        :param timeout: Wait for roughly this long before giving up.
        :return: the message.
        """
        for i, message in enumerate(self._got_from_queue):
            if self._matches(message, fields):
                self._got_from_queue = self._got_from_queue[:i] + self._got_from_queue[i+1:]
                return message
        no_of_waits = int(timeout / 0.5)
        for i in range(no_of_waits):
            try:
                new_message = self._queue.get(timeout=0.5)
                if self._matches(new_message, fields):
                    return new_message
                self._got_from_queue.append(new_message)
            except queue.Empty:
                pass
        raise RuntimeError("No such message")

    def _matches(self, message, fields):
        for field in fields:
            if field not in message or message[field] != fields[field]:
                return False
        return True
