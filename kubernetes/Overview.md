Overview of Kubernetes system
=============================
An overview of all the parts of the Nona micro service system.

Planning
--------
Like the existing ARM cluster we will have one Kubernetes master and two worker nodes.

We will need to rewrite the Kafka containers to work on x86. Perhaps there are ready made containers
we can use now.
Should we run two brokers? I _think_ that you should preferably have an odd number of brokers,
because of voting issues. Therefore we're stuck with one broker, until we invest in a third worker
node.

There should be two namespaces:

1. default, for production.
2. staging, for acceptance testing.

Acceptance test
---------------
The acceptance tests for the system as a whole are in the `features` directory in the Nona repo.
The will run in a Jobs pod inside Kubernetes. This simulates how the Slack service will actually
integrate with the system.

Current tasks
-------------
- Modify the Kafka deployment for x86
- Modify the Kafka deployment for the new cluster
- Deployment for puzzlestore
- Deployment for slackmessaging
- Job pod for acceptance tests
- Reporting of acceptance test results
- Avro Schema configuration
- Slack secrets for staging
- Slack secrets for default

Completed tasks
---------------
- Create the staging namespace
