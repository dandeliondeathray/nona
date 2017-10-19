from nonainterface.nonainterface import NonaInterface, AvroSchemas
from nonainterface.testing import ChatQueue


def before_all(context):
    context.schemas = AvroSchemas("../schema")


def before_scenario(context, scenario):
    context.team = "staging"
    # TODO: Use the real hostname inside the cluster
    bootstrap_servers = 'localhost:9092'
    #schemas = AvroSchemas('../schema')
    context.nonainterface = NonaInterface(team=context.team, bootstrap_servers=bootstrap_servers, schemas=schemas)
    context.nonainterface.start()
    context.chat_queue = ChatQueue(context.nonainterface.chat_events)


def after_scenario(context, scenario):
    context.nonainterface.stop()
