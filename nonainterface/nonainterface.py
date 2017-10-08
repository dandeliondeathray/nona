import queue
from collections import namedtuple
import io
import avro.io
import avro.schema


ChatMessage = namedtuple('ChatMessage', 'user_id team text')


class NonaInterface:
    def __init__(self, team=None):
        self.team = team if team else 'konsulatet'
        self.chat_events = queue.Queue(maxsize=1000)
        with open('../schema/Chat.avsc', 'r') as schema_file:
            schema = avro.schema.Parse(schema_file.read())
        self._reader = avro.io.DatumReader(writer_schema=schema)

    def _decode_chat_event(self, data):
        """
        Decode a chat event to a ChatMessage.

        :param data: a binary blob of Avro encoded data.
        :return: a ChatMessage object
        """
        decoder = avro.io.BinaryDecoder(io.BytesIO(data))
        message = self._reader.read(decoder)
        return ChatMessage(user_id=message['user_id'], team=message['team'], text=message['text'])

    def user_requests_puzzle(self, user_id):
        pass
