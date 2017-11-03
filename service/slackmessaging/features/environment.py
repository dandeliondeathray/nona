import time

import pymetamorph.metamorph as metamorph
import pymetamorph.process
import nonaspec.avro


def before_all(context):
    context.schemas = nonaspec.avro.AvroSchemas('../schema')
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    print("Metamorph: stdout {}, stderr {}".format(
          context.metamorph_process.out_name(),
          context.metamorph_process.err_name()))
    time.sleep(1)


def after_all(context):
    context.metamorph_process.stop()


def before_scenario(context, scenario):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.metamorph.request_kafka_reset(['nona_staging_Chat'])
    context.metamorph.await_reset_complete()

    env = {'SCHEMA_PATH': '../schema', 'KAFKA_BROKERS': 'localhost:9092', 'TEAMS': 'staging'}
    context.slackmessaging_process = pymetamorph.process.start(go='bin/slackmessaging', env=env)
    print("Slackmessaging output: stdout {}, stderr {}".format(
          context.slackmessaging_process.out_name(),
          context.slackmessaging_process.err_name()))


def after_scenario(context, scenario):
    context.slackmessaging_process.stop()

