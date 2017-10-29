from behave import *


@when(u'Erik has requested the puzzle')
def step_impl(context):
    context.nonainterface.user_requests_puzzle("Erik")


@then(u'the response is a nine letter word')
def step_impl(context):
	# TODO: Team should be "staging"
	# TODO: Match against word length, not actual word.
    context.chat_queue.has("Erik", "konsulatet", "PUSSGURKA")


@given(u'Erik has gotten five puzzles')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given Erik has gotten five puzzles')


@when(u'Johan has gotten five puzzles')
def step_impl(context):
    raise NotImplementedError(u'STEP: When Johan has gotten five puzzles')


@then(u'Erik and Johan have received the same puzzles')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then Erik and Johan have received the same puzzles')


@given(u'Erik has requested the puzzle')
def step_impl(context):
    context.nonainterface.user_requests_puzzle("Erik")


@given(u'he has solved the puzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given he has solved the puzzle')


@when(u'Erik has requested a new puzzle')
def step_impl(context):
    context.nonainterface.user_requests_puzzle("Erik")


@then(u'the response is a different puzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then the response is a different puzzle')
