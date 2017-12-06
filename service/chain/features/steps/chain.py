from behave import *
from hamcrest import *
import requests
import time


def request_new_round(context, team, seed):
    new_round = context.schemas.encode('NewRound', {'team': team, 'seed': seed})
    context.metamorph.send_message(topic='nona_NewRound', value=new_round)


@given(u'a new round')
def step_impl(context):
    context.team = "myteam"
    request_new_round(context, context.team, 0)
    time.sleep(2)


@when(u'a request is made for a puzzle at index {index}')
def step_impl(context, index):
    response = requests.get('http://localhost:8080/puzzle/{}/{}'.format(context.team, index))
    context.status_code = response.status_code
    context.body = response.text


@then(u'a puzzle is returned')
def step_impl(context):
    assert_that(context.status_code, equal_to(200))
    assert_that(context.body, is_not(equal_to("")))


@then(u'it is a different puzzle than the one before')
def step_impl(context):
    response = requests.get('http://localhost:8080/puzzle/{}/0'.format(context.team))
    context.prev_status_code = response.status_code
    context.prev_body = response.text
    assert_that(context.status_code, equal_to(200))
    assert_that(context.prev_status_code, equal_to(200))
    assert_that(context.body, is_not(equal_to(context.prev_body)))


@then(u'the solution is a different word')
def step_impl(context):
    letters_in_puzzle = ''.join(sorted(context.body))
    letters_in_prev_puzzle = ''.join(sorted(context.prev_body))
    assert_that(letters_in_puzzle, is_not(equal_to(letters_in_prev_puzzle)))


@given(u'no new round for a team')
def step_impl(context):
    context.team = "anotherteam"


@then(u'no puzzle was found')
def step_impl(context):
    assert_that(context.status_code, equal_to(404))
