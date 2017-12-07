from behave import *
from hamcrest import *
import requests
import time
import pymetamorph.process


def request_new_round(context, team, seed):
    new_round = context.schemas.encode('NewRound', {'team': team, 'seed': seed})
    context.metamorph.send_message(topic='nona_NewRound', value=new_round)


def start_chain_process(context):
    env = {'SCHEMA_PATH': '../schema',
           'KAFKA_BROKERS': 'localhost:9092',
           'DICTIONARY_PATH': '../../test/dictionary.txt',
           'NONA_CHAIN_PORT': str(context.port)}
    context.chain_process = pymetamorph.process.start(go='bin/nona_chain', env=env)
    time.sleep(2)


def stop_chain_process(context):
    context.chain_process.stop()


def wait_for_ready_fail(port):
    for i in range(5):
        time.sleep(2)
        try:
            requests.get('http://localhost:{}/readiness'.format(port))
        except:
            return
    raise RuntimeError("Service did not go down in time at port {}...".format(port))


def wait_for_ready(port):
    for i in range(5):
        time.sleep(2)
        try:
            response = requests.get('http://localhost:{}/readiness'.format(port))
            if response.status_code == 200:
                return
        except:
            pass
    raise RuntimeError("Service did not become ready in time at port {}...".format(port))


@given(u'a new round')
def step_impl(context):
    context.team = "myteam"
    request_new_round(context, context.team, 0)
    time.sleep(2)


@when(u'a request is made for a puzzle at index {index}')
def step_impl(context, index):
    response = requests.get('http://localhost:{}/puzzle/{}/{}'.format(context.port, context.team, index))
    context.status_code = response.status_code
    context.body = response.text


@then(u'a puzzle is returned')
def step_impl(context):
    assert_that(context.status_code, equal_to(200))
    assert_that(context.body, is_not(equal_to("")))


@then(u'it is a different puzzle than the one before')
def step_impl(context):
    response = requests.get('http://localhost:{}/puzzle/{}/0'.format(context.port, context.team))
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

@when(u'the service goes down and comes back up again')
def step_impl(context):
    stop_chain_process(context)
    wait_for_ready_fail(context.port)
    start_chain_process(context)
    wait_for_ready(context.port)


@then(u'a request for the puzzle at index 5 gives the same answer as before')
def step_impl(context):
    response = requests.get('http://localhost:{}/puzzle/{}/5'.format(context.port, context.team))
    previous_puzzle = context.body
    assert_that(response.status_code, equal_to(200))
    assert_that(response.text, equal_to(previous_puzzle))

