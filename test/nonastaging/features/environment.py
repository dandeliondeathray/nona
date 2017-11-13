import time
import threading
import asyncio
import pymetamorph.metamorph as metamorph
import pymetamorph.process
from nonastaging import NonaStagingService
from nonastagingclient import NonaStagingClient
from nonainterface import AvroSchemas


def before_all(context):
    context.metamorph_process = pymetamorph.process.start(go='bin/metamorph')
    time.sleep(1)
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()

    context.metamorph.request_kafka_reset(["nona_UserRequestsPuzzle"])
    context.metamorph.await_reset_complete()
    print("Metamorph reset complete")

    context.schemas = AvroSchemas('../../service/schema')
    context.service = NonaStagingService(brokers='localhost:9092',
                                         schema_path="../../service/schema")

    def start_service_with_loop():
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        context.service.start()

    context.service_thread = threading.Thread(target=start_service_with_loop)
    context.service_thread.start()
    time.sleep(2)


def after_all(context):
    context.service.stop()
    context.service_thread.join()
    context.metamorph_process.stop()


def before_scenario(context, _scenario):
    time.sleep(2)
    context.client = NonaStagingClient('ws://localhost:8765')
    context.client.start()


def after_scenario(context, _scenario):
    context.client.stop()
