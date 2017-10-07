import pymetamorph.metamorph as metamorph


def before_all(context):
    context.metamorph = metamorph.Metamorph()
    context.metamorph.connect()


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["nona_konsulatet_UserRequestsPuzzle"])
    context.metamorph.await_reset_complete()

