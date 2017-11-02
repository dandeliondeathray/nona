"""The Slack integration for the Nona system."""
import sys
import os
import threading
from slackclient import SlackClient
from tornado.ioloop import IOLoop, PeriodicCallback
import nonainterface


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
        self._nona = None
        self._user_channels = {}
        self._notification_channel = notification_channel
        self._chat_events_thread = threading.Thread(target=self._read_chat_events)

    def _read_chat_events(self):
        """Read events from Nona and send them to Slack."""
        print("Chat event thread started.")
        while True:
            chat_event = self._nona.chat_events.get()
            try:
                response_channel = self._user_channels[chat_event.user_id]
                self.enqueue(chat_event.text, response_channel)
            except KeyError:
                print("Did not find a response channel for user", chat_event.user_id)

    def enqueue(self, message, channel_id):
        """Send a message to a Slack user or channel, asynchronously."""
        print("Replying to Slack with message '{}' and channel id {}".format(message, channel_id))
        IOLoop.current().add_callback(self.client.rtm_send_message, channel_id, message)

    def _load_schemas(self):
        try:
            return nonainterface.AvroSchemas(os.environ['SCHEMA_PATH'])
        except KeyError:
            print("You must specify SCHEMA_PATH as the path to the schema directory.")
            sys.exit(1)

    def connect_to_nona(self):
        """Create and start the NonaInterface object, to connect to the Nona system."""
        schemas = self._load_schemas()
        try:
            kafka_brokers = os.environ["KAFKA_BROKERS"]
        except KeyError:
            print("You must specify KAFKA_BROKERS as a comma separated list of Kafka broker URLs")
            sys.exit(1)
        self._nona = nonainterface.NonaInterface("konsulatet", kafka_brokers, schemas)
        self._nona.start()

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
            for m in msgs:
                if 'type' in m and m['type'] == 'message' and 'channel' in m and 'user' in m:
                    self.handle_message(m)

    def handle_message(self, message):
        """Handle a Slack event of type 'message'."""
        print("MESSAGE:", message)
        user_id = message['user']
        self._user_channels[user_id] = message['channel']
        if message['text'].strip() == "!gemig":
            self._nona.user_requests_puzzle(user_id)
        else:
            print("Not a !gemig command.")

    def run_forever(self):
        """Connect to Slack and runs forever (until the connection fails)."""
        print("Starting NonaSlackApp...")
        self.connect_to_nona()
        self.connect_to_slack()
        self._read_msg_callback = PeriodicCallback(self.read_slack_messages, 500)
        self._read_msg_callback.start()
        self._chat_events_thread.start()
        IOLoop.current().handle_callback_exception = handle_callback_exception
        IOLoop.current().start()
