Specifications
==============
These executable specifications are intended to be run against a staging environment. The
implementation of the specifications make use of the team name "staging".

The specifications will (once I write them) test only the public interface of the system, from the
perspective of the Slack interface. It will not make use of the Slack service, as it's awkward to
test against Slack.

Currently, the idea is to run these tests as a Kubernetes Job.
