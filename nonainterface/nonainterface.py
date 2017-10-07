import queue

class NonaInterface:
    def __init__(self, team=None):
        self.team = team if team else 'konsulatet'
        self.chat_events = queue.Queue(maxsize=1000)

    def user_requests_puzzle(self, user_id):
        pass

