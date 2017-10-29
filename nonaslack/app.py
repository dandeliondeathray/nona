"""The Slack integration for the Nona system."""
import sys
import os
from slackclient import SlackClient
from tornado.ioloop import IOLoop, PeriodicCallback


class SlackException(RuntimeError):
    """A specific exception for when the Slack connection fails for whatever reason."""
    def __init__(self, *args, **kwargs):
        RuntimeError.__init__(*args, **kwargs)


def create_slack_client():
    """Create the Slack client used to connect to Slack."""
    try:
        token = os.environ["SLACK_API_TOKEN"].strip()
    except KeyError:
        print("You must specify a Slack API token in the 'SLACK_API_TOKEN' environment variable")
        return
    return SlackClient(token)


def handle_callback_exception(_callback):
    """Stops processing of all events in case of an unhandled exception. Essentially exits."""
    (ex_type, value, traceback) = sys.exc_info()
    print("EXCEPTION {}: {}", ex_type, value)
    print("Traceback: {}", traceback)
    IOLoop.current().stop()


class NonaSlackApp(object):
    """The Slack integration for the Nona system."""
    def __init__(self, notification_channel):
        self.client = None
        self._read_msg_callback = None
        self._notification_channel = notification_channel

    def enqueue(self, message, channel_id):
        """Send a message to a Slack user or channel, asynchronously."""
        print("Replying to Slack with message '{}' and channel id {}".format(message, channel_id))
        IOLoop.current().add_callback(self.client.rtm_send_message, channel_id, message)

    def connect_to_slack(self):
        """Connect to Slack."""
        print("Connecting to Slack...")
        self.client = create_slack_client()
        if not self.client:
            raise SlackException("Failed to create Slack client!")
        self.client.server.rtm_connect()
        print("Connected to Slack. Bot user name is '{}'".format(self.client.server.username))

    def read_slack_messages(self):
        """Poll for received messages from Slack."""
        msgs = self.client.rtm_read()
        if msgs:
            print("Received messages:", msgs)

    def run_forever(self):
        """Connect to Slack and runs forever (until the connection fails)."""
        print("Starting NonaSlackApp...")
        self.connect_to_slack()
        self._read_msg_callback = PeriodicCallback(self.read_slack_messages, 500)
        self._read_msg_callback.start()
        IOLoop.current().handle_callback_exception = handle_callback_exception
        IOLoop.current().start()
