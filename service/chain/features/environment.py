import time

import pymetamorph.metamorph as metamorph
import pymetamorph.process
import nonaspec.avro

service_port = 8100

def before_all(context):
    context.schemas = nonaspec.avro.AvroSchemas('../schema')
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    time.sleep(1)

def before_scenario(context, scenario):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.metamorph.request_kafka_reset([])
    context.metamorph.await_reset_complete()

    context.port = service_port

    env = {'SCHEMA_PATH': '../schema',
           'KAFKA_BROKERS': 'localhost:9092',
           'DICTIONARY_PATH': '../../test/dictionary.txt',
           'NONA_CHAIN_PORT': str(context.port)}
    context.chain_process = pymetamorph.process.start(go='bin/nona_chain', env=env)
    time.sleep(2)


def after_scenario(context, scenario):
    context.chain_process.stop()
    global service_port
    service_port += 1


def after_all(context):
    context.metamorph_process.stop()