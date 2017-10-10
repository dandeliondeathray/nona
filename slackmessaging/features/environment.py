import pymetamorph.metamorph as metamorph
import pymetamorph.process

import nonaspec.avro


def before_all(context):
    context.schemas = nonaspec.avro.AvroSchemas('../schema')


def before_scenario(context, scenario):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.metamorph.request_kafka_reset(['nona_konsulatet_Chat'])
    context.metamorph.await_reset_complete()

    context.slackmessaging_process = pymetamorph.process.start(go='bin/slackmessaging')


def after_scenario(context, scenario):
    context.slackmessaging_process.stop()

