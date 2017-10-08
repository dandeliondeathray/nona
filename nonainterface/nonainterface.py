import queue
from collections import namedtuple
import io
import avro.io
import avro.schema
import confluent_kafka
import threading
import base64
import io


ChatMessage = namedtuple('ChatMessage', 'user_id team text')


class ChatConsumer:
    def __init__(self, queue, schema, team):
        self._queue = queue
        self._avro_reader = avro.io.DatumReader(writer_schema=schema)
        self._consumer = None
        self._event = threading.Event()
        self._thread = None
        self.consumer_group = "nonainterface"
        self.team = team

    def start(self):
        self._thread = threading.Thread(target=self.poll, name="ChatConsumer")
        self._thread.start()

    def stop(self):
        self._event.set()
        self._thread.join(timeout=3.0)

    def poll(self):
        consumer = confluent_kafka.Consumer({'bootstrap.servers': 'localhost:9092',
                                             'group.id': self.consumer_group,
                                             'default.topic.config': {'auto.offset.reset': 'smallest'}})
        consumer.subscribe(["nona_{team}_Chat".format(team=self.team)])
        while not self._event.is_set():
            msg = consumer.poll(1)
            if msg is None:
                continue
            if not msg.error():
                binary_message = base64.b64decode(msg.value())
                chat_message = self._decode_chat_event(binary_message)
                self._queue.put(chat_message)
            elif msg.error().code() != confluent_kafka.KafkaError._PARTITION_EOF:
                print(msg.error())
        consumer.close()

    def _decode_chat_event(self, data):
        """
        Decode a chat event to a ChatMessage.

        :param data: a binary blob of Avro encoded data.
        :return: a ChatMessage object
        """
        decoder = avro.io.BinaryDecoder(io.BytesIO(data))
        message = self._avro_reader.read(decoder)
        return ChatMessage(user_id=message['user_id'], team=message['team'], text=message['text'])


class NonaInterface:
    def __init__(self, team=None):
        self.team = team if team else 'konsulatet'
        self.chat_events = queue.Queue(maxsize=1000)
        with open('../schema/Chat.avsc', 'r') as schema_file:
            schema = avro.schema.Parse(schema_file.read())
        with open('../schema/UserRequestsPuzzle.avsc', 'r') as schema_file:
            self.user_req_puzzle_schema = avro.schema.Parse(schema_file.read())
        self.chat_consumer = ChatConsumer(self.chat_events, schema, self.team)
        self.producer = confluent_kafka.Producer({'bootstrap.servers': 'localhost:9092'})
        self.user_req_puzzle_topic = 'nona_{team}_UserRequestsPuzzle'.format(team=self.team)

    def start(self):
        self.chat_consumer.start()

    def stop(self):
        self.chat_consumer.stop()

    def user_requests_puzzle(self, user_id):
        writer = avro.io.DatumWriter(writer_schema=self.user_req_puzzle_schema)
        out = io.BytesIO()
        encoder = avro.io.BinaryEncoder(out)
        data = {'user_id': user_id, 'team': self.team, 'timestamp': 0}
        writer.write(data, encoder)
        self.producer.produce(self.user_req_puzzle_topic, out.getvalue())
