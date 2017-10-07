from behave import *


@given(u'a team konsulatet')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given a team konsulatet')


@when(u'a user requests a puzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: When a user requests a puzzle')


@then(u'a UserRequestsPuzzle is sent to topic nona_konsulatet_UserRequestsPuzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then a UserRequestsPuzzle is sent to topic nona_konsulatet_UserRequestsPuzzle')


@when(u'there is a puzzle response in nona_konsulatet_ChatMessage')
def step_impl(context):
    raise NotImplementedError(u'STEP: When there is a puzzle response in nona_konsulatet_ChatMessage')


@then(u'a chat message is in the event queue')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then a chat message is in the event queue')
