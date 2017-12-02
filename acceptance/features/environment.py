"""Acceptance test environment."""
import os
from nonacontrolclient import NonaControlClient

def before_all(context):
    context.staging_address = os.environ['NONA_STAGING']


def before_scenario(context, _scenario):
    """Create a fresh NonaControlClient and connect it."""
    context.team = "staging"
    context.client = NonaControlClient(context.staging_address)
    context.client.start()


def after_scenario(context, _scenario):
    """Stop the NonaControlClient."""
    context.client.stop()
