import pymetamorph.metamorph as metamorph
import pymetamorph.process

import nonaspec.avro


def before_all(context):
    context.schemas = nonaspec.avro.AvroSchemas('../schema')


def before_scenario(context, scenario):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.metamorph.request_kafka_reset(['nona_PuzzleNotification'])
    context.metamorph.await_reset_complete()

    context.puzzlestore_process = pymetamorph.process.start(go='bin/nona_puzzlestore')


def after_scenario(context, scenario):
    context.puzzlestore_process.stop()

