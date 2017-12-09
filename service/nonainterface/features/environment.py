import pymetamorph.metamorph as metamorph
import pymetamorph.process
from nonainterface import NonaInterface, AvroSchemas
from nonaspec.interface import ChatQueue
import time


def before_all(context):
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    time.sleep(1)
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()
    context.schemas = AvroSchemas("../schema")


def after_all(context):
    context.metamorph_process.stop()

def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["nona_UserRequestsPuzzle",
                                           "nona_SolutionRequest"])
    context.metamorph.await_reset_complete()

    schemas = AvroSchemas('../schema')
    context.nonainterface = NonaInterface(team='konsulatet', bootstrap_servers='localhost:9092', schemas=schemas)
    context.nonainterface.start()
    context.chat_queue = ChatQueue(context.nonainterface.chat_events)


def after_scenario(context, scenario):
    context.nonainterface.stop()
