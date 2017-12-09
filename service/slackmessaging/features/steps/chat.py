from behave import *
from pymetamorph.metamorph import OnTopic
from hamcrest import *


@given(u'that the team is enabled')
def step_impl(context):
    context.user_id = 'U1'


@when(u'a user is notified of the puzzle PUSSGURKA')
def step_impl(context):
    puzzle_response = context.schemas.encode('PuzzleNotification',
                                             {'user_id': context.user_id, 'team': 'staging', 'puzzle': 'PUSSGURKA'})
    context.metamorph.send_message(topic='nona_PuzzleNotification', value=puzzle_response)


@then(u'a chat message containing "{text}" is sent to the same user')
def step_impl(context, text):
    metamorph_message = context.metamorph.await_message(matcher=OnTopic('nona_staging_Chat'))
    m = context.schemas.decode('Chat', metamorph_message['message'])
    assert_that(m['user_id'], equal_to(context.user_id))
    assert_that(m['text'], contains_string(text))


@when(u'a user solved the puzzle')
def step_impl(context):
    message = context.schemas.encode('CorrectSolution',
                                     {'user_id': context.user_id, 'team': 'staging'})
    context.metamorph.send_message(topic='nona_CorrectSolution', value=message)


@when(u'a user attempts an incorrect word')
def step_impl(context):
    message = context.schemas.encode('IncorrectSolution',
                                     {'user_id': context.user_id, 'team': 'staging'})
    context.metamorph.send_message(topic='nona_IncorrectSolution', value=message)
