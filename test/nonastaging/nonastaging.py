"""
Nonastaging is a WebSocket based interface to the Nona system. Its purpose is to simplify
acceptance testing of the Nona system deployed in a staging environment.
"""

import asyncio
import os
from nonainterface import NonaInterface, AvroSchemas


class Nonastaging:
    """Nonastaging implements the entire micro service."""
    def __init__(self, bootstrap_servers, schema_path):
        schemas = AvroSchemas(schema_path)
        self._nona = NonaInterface("staging", bootstrap_servers, schemas)

    def run_forever(self):
        """Start the event loop and only return when it is stopped."""
        loop = asyncio.get_event_loop()
        loop.run_forever()


if __name__ == "__main__":
    nonastaging = Nonastaging(os.environ['KAFKA_BROKERS'], os.environ['SCHEMA_PATH'])
    nonastaging.run_forever()
