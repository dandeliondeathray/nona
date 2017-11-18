from behave import *
import pymetamorph.metamorph as metamorph
from pymetamorph.metamorph import MatchThese, OnTopic
from hamcrest import *


@when(u'a chat message is sent to topic nona_staging_Chat')
def step_impl(context):
    chat_message = context.schemas.encode('Chat', {'user_id': 'U1',
                                                   'team': 'staging',
                                                   'text': 'Some chat text'})
    context.metamorph.send_message(topic='nona_staging_Chat', value=chat_message)


@then(u'that chat message is received on the WebSocket')
def step_impl(context):
    message = context.client.await_chat(user_id='U1', team='staging', text='Some chat text')
    assert_that(message['user_id'], equal_to('U1'))
    assert_that(message['team'], equal_to('staging'))
    assert_that(message['text'], equal_to('Some chat text'))


@when(u'a user requests a puzzle')
def step_impl(context):
    context.client.user_requests_puzzle(user_id='U1')


@then(u'a request is sent to nona_UserRequestsPuzzle')
def step_impl(context):
    match_message = MatchThese(OnTopic('nona_UserRequestsPuzzle'))
    m = context.metamorph.await_message(matcher=match_message)