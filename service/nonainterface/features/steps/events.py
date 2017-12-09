from behave import *
from pymetamorph.metamorph import MatchThese, OnTopic
import base64


@given(u'a team konsulatet')
def step_impl(context):
    context.nonainterface.team = 'konsulatet'


@when(u'a user requests a puzzle')
def step_impl(context):
    context.nonainterface.user_requests_puzzle(user_id='U1')

@when(u'a user sends an attempted solution')
def step_impl(context):
    context.nonainterface.try_word("PUSSGURKA", user_id='U1')

@then(u'a {message_type} is sent to topic {topic}')
def step_impl(context, message_type, topic):
    # TODO: Match against user id and team
    match_message = MatchThese(OnTopic(topic))
    context.metamorph.await_message(matcher=match_message)


@when(u'there is a puzzle response in nona_konsulatet_Chat')
def step_impl(context):
    puzzle_response = context.schemas.encode('Chat', {'user_id': 'U1', 'team': 'konsulatet', 'text': 'PUSSGURKA'})
    context.metamorph.send_message(topic='nona_konsulatet_Chat', value=puzzle_response)


@then(u'a chat message is in the event queue')
def step_impl(context):
    context.chat_queue.has(user_id='U1', team='konsulatet', text='PUSSGURKA', timeout=5.0)
