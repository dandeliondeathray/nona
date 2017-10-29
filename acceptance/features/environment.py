"""Acceptance test environment."""
import os
from nonainterface import NonaInterface, AvroSchemas
from nonaspec.interface import ChatQueue


def before_all(context):
    """Set up Kafka brokers URLs and load Avro schemas."""
    context.bootstrap_servers = os.environ['KAFKA_BROKERS']
    context.schemas = AvroSchemas(os.environ['SCHEMA_PATH'])


def before_scenario(context, _scenario):
    """Create a fresh NonaInterface and connect it."""
    # TODO: Team should be "staging"
    context.team = "konsulatet"
    context.nonainterface = NonaInterface(team=context.team,
                                          bootstrap_servers=context.bootstrap_servers,
                                          schemas=context.schemas)
    context.nonainterface.start()
    context.chat_queue = ChatQueue(context.nonainterface.chat_events)


def after_scenario(context, _scenario):
    """Stop the NonaInterface."""
    context.nonainterface.stop()
