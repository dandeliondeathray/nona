from behave import *


@given(u'that the service is ready')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given that the service is ready')


@when(u'an HTTP ready probe checks the service')
def step_impl(context):
    raise NotImplementedError(u'STEP: When an HTTP ready probe checks the service')


@then(u'an OK response is returned')
def step_impl(context):
    raise NotImplementedError(u'STEP: Then an OK response is returned')


@given(u'that the service is up and running')
def step_impl(context):
    raise NotImplementedError(u'STEP: Given that the service is up and running')


@when(u'an HTTP liveness probe checks the service')
def step_impl(context):
    raise NotImplementedError(u'STEP: When an HTTP liveness probe checks the service')
