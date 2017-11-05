"""
nonaspec.interface contains code for helping with the testing of the Nona system through
nonainterface, as well as nonainterface itself.
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

    def has(self, user_id, team, text, timeout=5.0):
        """
        Check if the chat queue has a message matching specific fields.

        :param fields: A dict of all fields and their expected value.
        :param timeout: Wait for roughly this long before giving up.
        :return: the message.
        """
        for i, message in enumerate(self._got_from_queue):
            if self._matches(message, user_id, team, text):
                self._got_from_queue = self._got_from_queue[:i] + self._got_from_queue[i+1:]
                return message
        no_of_waits = int(timeout / 0.5)
        for i in range(no_of_waits):
            try:
                new_message = self._queue.get(timeout=0.5)
                if self._matches(new_message, user_id, team, text):
                    return new_message
                self._got_from_queue.append(new_message)
            except queue.Empty:
                pass
        raise RuntimeError("No such message")

    def _matches(self, message, user_id, team, text):
        return message.user_id == user_id and message.team == team and message.text == text