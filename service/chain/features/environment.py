import time

import pymetamorph.metamorph as metamorph
import pymetamorph.process
import nonaspec.avro


def before_all(context):
    context.schemas = nonaspec.avro.AvroSchemas('../schema')
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    time.sleep(1)

    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.metamorph.request_kafka_reset([])
    context.metamorph.await_reset_complete()

    env = {'SCHEMA_PATH': '../schema', 'KAFKA_BROKERS': 'localhost:9092'}
    context.chain_process = pymetamorph.process.start(go='bin/nona_chain', env=env)


def after_all(context):
    context.chain_process.stop()
    context.metamorph_process.stop()