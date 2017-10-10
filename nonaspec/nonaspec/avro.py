"""
Helpers for encoding and decoding Avro messages.

Schemas are read from a directory.
"""

import glob
import os.path
import io
import base64
import avro.schema
import avro.io


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

    def decode(self, schema_name, data):
        schema = self._schemas[schema_name]
        reader = avro.io.DatumReader(writer_schema=schema)
        decoder = avro.io.BinaryDecoder(io.BytesIO(data))
        return reader.read(decoder)

    def get(self, schema_name):
        return self._schemas[schema_name]