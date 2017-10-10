from behave import *
import nonaspec.probe as probe
import time
from hamcrest import *


@given(u'that the service is ready')
def step_impl(context):
    # TODO: Put the service in a ready state, or wait a few seconds?
    time.sleep(3)


@when(u'an HTTP ready probe checks the service')
def step_impl(context):
    context.probe_response = probe.http_probe('http://localhost:24689/readiness')


@then(u'an OK response is returned')
def step_impl(context):
    assert_that(context.probe_response, "Probe did not respond with OK")


@given(u'that the service is up and running')
def step_impl(context):
    # TODO: Put the service in a live state?
    time.sleep(3)


@when(u'an HTTP liveness probe checks the service')
def step_impl(context):
    context.probe_response = probe.http_probe('http://localhost:24689/liveness')
