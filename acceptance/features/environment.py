"""Acceptance test environment."""
import os
from nonastagingclient import NonaStagingClient

def before_all(context):
    context.staging_address = os.environ['NONA_STAGING']


def before_scenario(context, _scenario):
    """Create a fresh NonaStagingClient and connect it."""
    context.team = "staging"
    context.client = NonaStagingClient("ws://mystaging:8765")
    context.client.start()


def after_scenario(context, _scenario):
    """Stop the NonaStagingClient."""
    context.client.stop()
