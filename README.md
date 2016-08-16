# drs

DRS is a library that implements:
* a standard way for services to communicate, using pluggable transport methods, message protocols. Each transport method implements both a server and a client.
* It supports one way (`Fire(...)`) and two-way communication (`Call(...)`)
* The standard unit of communication is a Command that looks like this
```
{
  "action": "myAction",
  "body": {
    "nested": {
      "map": 1
    }
  },
  "key": "i-am-unique"
}
```
* Handlers. Whenever a message is received, depending on its `action`, functions to be triggered can be specified. It is fault tolerant so it will not crash, but report the errors in the logs.

Communication between two services can be seen as a `Connection` between a `Server` and one or many `Clients`. Each `Connection` consists of an open `Stream` on which each side can either read or write, this stream is created with the `Transport` method chosen. Each message that is written or read is called a `Command`, and follows a specific format. A `Command` is one unit of communication. This way both sides know what to expect. This command can then be serialized to a string using different `Protocols`.

# Concepts

## Command
`Command` is a structure that has two required fields, an `Action` of type `string` and a `Body`, which is an `interface{}` but should really be a `map[string]interface{}`. It can also have a unique `Key` that will help the processors keep track of who sent what.

## Handlers - subscribing to messages
On both server and client side, doing `Conn.On(Action, doThis...)` with `Action` being a string to listen to, and `doThis...` functions to be called whenever a message is received.

## Sending a message
`Conn.Fire(cmd)` where cmd is a `Command` type, with an `Action` and a `Body`.

## Broadcast a message across all the connections (from the server side)
`Server.Broadcast(cmd)` basically doing `Conn.Fire(cmd)` for every `Conn` registered by the server.

## Reserved messages

With actions `drs.error`, `drs.exception`, `drs.response`. Whenever a message is fired, the other part will process it depending on its handlers and responds to the first part with either of these three messages, to report the status of its processing.

It's only really read by the first part if the message was sent using a `Call`. Then the message key will be stored in `pending` until getting the response. And deleted once got.

# How to use

## Creating a server

1. Needs to create the connection by doing `drs.New` with appropriate protocol and transport.
2. Then define the handlers.
3. Then `conn.Listen` on the appropriate port. It will essentially just run the boilerplate for the specific transport, as well as `connect` and `disconnect` functions to be run on start and end. Then, it will do `conn.Read`, which waits for messages in the stream.

```
server := drs.New(transport, prot)

// Add a handler
server.On("thisAction", func(msg *drs.Message) (interface{}, error) {
  // do Stuff
})

// Send something
server.Fire(&Command{
  Action: "yo",
  Body: map[string]string{}{
    "this": "that"
  }
})

// At the very end, listen for incoming connections, that really starts the server
server.Listen(host)
```

## Client connecting to a server
1. Create the client doing `drs.Dial` with the appropriate host, transport and protocol.
2. Start a goroutine that will `conn.Read()`
3. Declare handlers if needed
4. `Fire` or `Call` when wanted

```
// Create the connection
client, dialErr := drs.Dial(prot, transport, clientHost)

// Listen to messages
go func() {
  client.Read()
  os.Exit(0)
}()

// Add a handler
client.On("thisAction", func(msg *drs.Message) (interface{}, error) {
  // do sth
})

// Send something
client.Fire(&Command{
  Action: "yo",
  Body: map[string]string{}{
    "this": "that"
  }
})

```

where `transport` can be `ipc.New()`, `tcp.New()`, `ws.New(dynamic.Empty())` and the protocol `prot` can be `protocol.JSON`.
