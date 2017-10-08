import glob
import os.path
import io
import base64

import pymetamorph.metamorph as metamorph
from nonainterface import NonaInterface
from testing import ChatQueue
import avro.io
import avro.schema


class AvroSchemas:
    def __init__(self, schema_path):
        schema_files = glob.glob(os.path.join(schema_path, '*.avsc'))
        self._schemas = {}
        for schema_file in schema_files:
            with open(schema_file, 'r') as schema_handle:
                schema = avro.schema.Parse(schema_handle.read())
                self._schemas[schema.name] = schema

    def encode(self, schema_name, fields):
        schema = self._schemas[schema_name]
        out = io.BytesIO()
        writer = avro.io.DatumWriter(schema)
        encoder = avro.io.BinaryEncoder(out)
        writer.write(fields, encoder)
        out_bytes = out.getvalue()
        return base64.b64encode(out_bytes).decode('UTF-8')

    def get(self, schema_name):
        return self._schemas[schema_name]


def before_all(context):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.schemas = AvroSchemas("../schema")


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["nona_konsulatet_UserRequestsPuzzle"])
    context.metamorph.await_reset_complete()

    context.nonainterface = NonaInterface(team='konsulatet', bootstrap_servers='localhost:9092')
    context.nonainterface.start()
    context.chat_queue = ChatQueue(context.nonainterface.chat_events)


def after_scenario(context, scenario):
    context.nonainterface.stop()
