import time
import pymetamorph.metamorph as metamorph
import pymetamorph.process
from nonastaging import NonaStagingService
from nonastagingclient import NonaStagingClient


def before_all(context):
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    time.sleep(1)
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()


def after_all(context):
    context.metamorph_process.stop()


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["nona_UserRequestsPuzzle"])
    context.metamorph.await_reset_complete()

    context.service = NonaStagingService(brokers='localhost:9092',
                                         schema_path="../../service/schema")
    context.service.start()
    context.ws = NonaStagingClient('ws://localhost:8765')
    context.ws.start()


def after_scenario(context, _scenario):
    context.ws.stop()
    context.service.stop()
