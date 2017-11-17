from behave import *
from nonastagingclient import Matcher
import re


class PuzzleMatcher(Matcher):
    def matches(self, value):
        """Match that the value matches what a puzzle is displayed as."""
        # TODO: Verify I'm using this right when I have an internet connection.
        return re.match("[A-ZÅÄÖ]{3} [A-ZÅÄÖ]{3} [A-ZÅÄÖ]{3}", value) is not None

    def __str__(self):
        return "<a puzzle>"


@when(u'Erik has requested the puzzle')
def step_impl(context):
    context.client.user_requests_puzzle("Erik")


@then(u'the response is a nine letter word')
def step_impl(context):
    context.client.await_chat("Erik", "konsulatet", PuzzleMatcher())


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
    context.client.user_requests_puzzle("Erik")


@given(u'he has solved the puzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given he has solved the puzzle')


@when(u'Erik has requested a new puzzle')
def step_impl(context):
    context.client.user_requests_puzzle("Erik")


@then(u'the response is a different puzzle')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then the response is a different puzzle')
