Feature: Probes

  Scenario: Readiness probe
    Given that the service is ready
     When an HTTP ready probe checks the service
     Then an OK response is returned

  Scenario: Liveness probe
    Given that the service is up and running
     When an HTTP liveness probe checks the service
     Then an OK response is returned
