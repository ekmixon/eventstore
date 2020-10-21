# EventStore

EventStore is a service that store data from events:

- At 3 different scopes: global, bridge and instance
- Inside a time window defined by TTL

## EventStore Server

The event store server exposes three functions: Save, Load and Delete

- Save, given location (scope data and key) stores values for the time defined at TTL.
- Load, given location (scope data and key) retieves non expired values.
- Delete, given location (scope data and key) deletes stored values.

## Support

We would love your feedback and help on these sources, so don't hesitate to let us know what is wrong and how we could improve them, just file an [issue](https://github.com/triggermesh/eventstore/issues/new) or join those of use who are maintaining them and submit a PR.

## Commercial Support

TriggerMesh Inc supports those sources commercially, email info@triggermesh.com to get more details.

## Code of Conduct

This project is by no means part of [CNCF](https://www.cncf.io/) but we abide
by its
[code of conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md)
