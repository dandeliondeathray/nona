import os
import tornado.web
from app import NonaSlackApp


class HealthHandler(tornado.web.RequestHandler):
    def get(self):
        self.finish()


class ReadinessHandler(tornado.web.RequestHandler):
    def get(self):
        self.finish()


def read_environment_var(name):
    try:
        return os.environ[name].strip()
    except KeyError:
        raise OSError("Missing required environment variable {}".format(name))


if __name__ == "__main__":
    health_app = tornado.web.Application([
        (r"/health", HealthHandler),
        (r"/readiness", ReadinessHandler)
    ])
    health_app.listen(24689)

    notification_channel = read_environment_var("NOTIFICATION_CHANNEL")
    app = NonaSlackApp(notification_channel)
    app.run_forever()
