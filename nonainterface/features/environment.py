import pymetamorph.metamorph as metamorph
from nonainterface import NonaInterface, AvroSchemas
from testing import ChatQueue


def before_all(context):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.schemas = AvroSchemas("../schema")


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["nona_konsulatet_UserRequestsPuzzle"])
    context.metamorph.await_reset_complete()

    schemas = AvroSchemas('../schema')
    context.nonainterface = NonaInterface(team='konsulatet', bootstrap_servers='localhost:9092', schemas=schemas)
    context.nonainterface.start()
    context.chat_queue = ChatQueue(context.nonainterface.chat_events)


def after_scenario(context, scenario):
    context.nonainterface.stop()
