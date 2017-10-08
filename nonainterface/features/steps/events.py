from behave import *
from pymetamorph.metamorph import MatchThese, OnTopic
import base64


@given(u'a team konsulatet')
def step_impl(context):
    context.nonainterface.team = 'konsulatet'


@when(u'a user requests a puzzle')
def step_impl(context):
    context.nonainterface.user_requests_puzzle(user_id='U1')


@then(u'a UserRequestsPuzzle is sent to topic nona_konsulatet_UserRequestsPuzzle')
def step_impl(context):
    # TODO: Match against user id and team
    match_message = MatchThese(OnTopic('nona_konsulatet_UserRequestsPuzzle'))
    context.metamorph.await_message(matcher=match_message)


@when(u'there is a puzzle response in nona_konsulatet_Chat')
def step_impl(context):
    puzzle_response = context.schemas.encode('Chat', {'user_id': 'U1', 'team': 'konsulatet', 'text': 'PUSSGURKA'})
    context.metamorph.send_message(topic='nona_konsulatet_Chat', value=puzzle_response)


@then(u'a chat message is in the event queue')
def step_impl(context):
    context.chat_queue.has({'user_id': 'U1', 'team': 'konsulatet', 'text': 'PUSSGURKA'}, timeout=2.0)
