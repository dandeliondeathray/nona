from behave import *
from pymetamorph.metamorph import OnTopic
from hamcrest import *


@given(u'that the team is enabled')
def step_impl(context):
    pass


@when(u'a user is notified of the puzzle PUSSGURKA')
def step_impl(context):
    context.user_id = 'U1'
    puzzle_response = context.schemas.encode('PuzzleNotification',
                                             {'user_id': context.user_id, 'team': 'konsulatet', 'puzzle': 'PUSSGURKA'})
    context.metamorph.send_message(topic='nona_PuzzleNotification', value=puzzle_response)


@then(u'a chat message containing PUSSGURKA is sent to the same user')
def step_impl(context):
    metamorph_message = context.metamorph.await_message(matcher=OnTopic('nona_konsulatet_Chat'))
    print(metamorph_message)
    m = context.schemas.decode('Chat', metamorph_message['message'])
    assert_that(m['user_id'], equal_to(context.user_id))
    assert_that(m['text'], equal_to('PUSSGURKA'))
