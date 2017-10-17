from behave import *
from pymetamorph.metamorph import OnTopic
from hamcrest import *



#@when(u'a user is notified of the puzzle PUSSGURKA')
#def step_impl(context):
#    context.user_id = 'U1'
#    puzzle_response = context.schemas.encode('PuzzleNotification',
#                                             {'user_id': context.user_id, 'team': 'konsulatet', 'puzzle': 'PUSSGURKA'})
#    context.metamorph.send_message(topic='nona_PuzzleNotification', value=puzzle_response)


#@then(u'a chat message containing PUSSGURKA is sent to the same user')
#def step_impl(context):
#    metamorph_message = context.metamorph.await_message(matcher=OnTopic('nona_konsulatet_Chat'))
#    print(metamorph_message)
#    m = context.schemas.decode('Chat', metamorph_message['message'])
#    assert_that(m['user_id'], equal_to(context.user_id))
#    assert_that(m['text'], equal_to('PUSSGURKA'))
@given(u'that a users current puzzle is PUSSGURKA')
def step_impl(context):
    #context.puzzle = "PUSSGURKA"
    pass

@when(u'a request is sent for the current puzzle')
def step_impl(context):
    context.user_id = 'U1'
    puzzle_request = context.schemas.encode('UserRequestsPuzzle',
                                            {'user_id': context.user_id, 'team': 'konsulatet', 'timestamp': 0})
    context.metamorph.send_message(topic='nona_UserRequestsPuzzle', value=puzzle_request)

@then(u'a puzzle notification is sent for puzzle PUSSGURKA')
def step_impl(context):
    metamorph_message = context.metamorph.await_message(matcher=OnTopic('nona_PuzzleNotification'))
    print(metamorph_message)
    m = context.schemas.decode('PuzzleNotification', metamorph_message['message'])
    assert_that(m['user_id'], equal_to(context.user_id))
    assert_that(m['puzzle'], equal_to('PUSSGURKA'))
    assert_that(m['team'], equal_to('konsulatet'))