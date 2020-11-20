# EventStore

EventStore is an interface for storing ephemeral data in an event flow.

## Clients and Servers

EventStore server functionality can be implemented by using the [Protocol Buffers definition](./pkg/protob/eventstore.proto) included in this repository, find the generated go client at the same folder.

A recommended go client with a much simpler interface [is also provided](./pkg/client/eventstore.go).

## EventStore Interface

The EventStore interface stores data at three eventing levels: global, bridge and instance

- `Global`: data stored at the global level is shared among all components that use the same EventStore by just specifying the key where the values are stored.

- `Bridge`: the term bridge is used at Triggermesh as the grouping of eventing components that play a role in an event flow, which are marked with the bridge name they belong to. The bridge level stored values can be used by components at the same bridge only. Note that only the name _bridge_ is borrowed, and no dependency exists with Triggermesh bridges.

- `Instance`: a bridge flow is started by one component that might be able to stamp an instance identifier on it. As long as that identifier exists and is propagate through the event flow, it can be used to store data only available at that single instance inside the bridge.

Depending on the level where data is stored a location needs to be specified:

- Global data location needs: key
- Bridge data location needs: bridge, key
- Instance data location needs: bridge, instance, key

At all three levels the available operations are `Save`, `Load`, `Delete`.

### Stored Data

Data is stored as a byte array. It is up to the client code storing or retrieving the data to perform serialization.

### Data Expiry

Data needs to include at `Save` time a value for `TTL` (Time to Live) parameter that informs the number of seconds (int32) that the data will be retrievable at the store.

## EventStore Client

The Go client provided in this repo exposes the connection management at its base object, and the same `Save`, `Load`, `Delete` functions at each of the `Global`, `Bridge` and `Instance` levels.

### Connection

Client instantiation and connection requires the EventStore server and a timeout in milliseconds.

```go
import "github.com/triggermesh/eventstore/pkg/client"

...

c := client.New("dns:///inmemorystorage-triggermesh.tm-demo:8080", 5000)
err := c.Connect(ctx)

...
```

It is recommended that after succesful connection a defered function is created for disconnection.

```go
defer func() { err = c.Disconnect() }()
```

### Levels

Each of the EventStore levels can be chosen by informing their parameters.

```go
global := c.Global()
myBrige := c.Bridge("my-bridge")
myBrigeInstance := c.Instance("my-bridge","aaee-1122")
```

### Interface Methods

Given the client for one of the levels, we can load, delete and delete values. When saving we need to provide the value and also the time to live in seconds.

```go
err := myBrigeInstance.SaveValue(ctx, "invoice.total", []byte("103¥"), 20)

...

err := myBrigeInstance.LoadValue(ctx, "invoice.total")

...

err := myBrigeInstance.DeleteValue(ctx, "invoice.total")
```

`LoadValue` function will return an error when trying to load a value that doesn't exists or have been expired.

## Example Client

An example client is included at this repository. When running in kubernetes the easiest way to test it is using ko to create a pod where the binary will be present.

```sh
ko apply -n mynamespace -f ./config/
```

Once the pod is running, you can exec bash in it and interact with the backing EventStore:

```sh
$ kubectl exec -ti -n tmsamples eventstores-client -- bash

# eventstore-client  \
    --command load \
    --server dns:///inmemorystorage-triggermesh.tm-demo:8080 \
    --scope Bridge \
    --bridge salesforce-elastic \
    --key /data/ChangeEvents.replayId

```

## Support

We would love your feedback and help on these sources, so don't hesitate to let us know what is wrong and how we could improve them, just file an [issue](https://github.com/triggermesh/eventstore/issues/new) or join those of use who are maintaining them and submit a PR.

## Commercial Support

TriggerMesh Inc supports those sources commercially, email info@triggermesh.com to get more details.

## Code of Conduct

This project is by no means part of [CNCF](https://www.cncf.io/) but we abide
by its
[code of conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md)
