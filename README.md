# learn-pub-sub-starter (Peril)

This is the starter code used in Boot.dev's [Learn Pub/Sub](https://learn.boot.dev/learn-pub-sub) course.

## Publisher & Queue

### Server Commands

The server runs an interactive REPL. Available commands:

- `pause` — publishes a pause message to all connected clients
- `resume` — publishes a resume message to all connected clients
- `quit` — stops the server

### Architecture Notes

This project uses a Pub/Sub pattern via RabbitMQ instead of direct
point-to-point messaging. The server publishes game state changes
(pause/resume) without knowing which or how many clients are
listening. Each client gets its own durable queue
(e.g. `pause.<username>`), so messages aren't lost even if a
client is offline when the message is published.

### Client CLI Commands

Run the client to enter the interactive REPL:

* `spawn <location> <unit_type>` - Spawns a unit (`infantry`, `cavalry`, `artillery`) in a region.
* `move <location> <unit_id>` - Moves a specified unit to a destination.
* `status` - Displays the current state of your units and locations.
* `help` - Lists available commands.
* `quit` - Exits the game client.

## Game Logs Queue

The server declares a durable `game_logs` queue and binds it to the Peril topic exchange. Durable queues survive RabbitMQ restarts and can be shared by multiple consumers.

## Consumers

The client now subscribes to the `pause` routing key on the direct exchange as soon as it starts, using `pubsub.SubscribeJSON`. This runs in a background goroutine, so the client can simultaneously listen for pause/resume messages from the server *and* accept player commands (`spawn`, `move`, `status`, etc.) in its main loop.

Each client binds to its own transient queue named `pause.<username>`, ensuring that pause state is tracked independently per player and cleaned up automatically when the client disconnects.

When a pause message is received, `handlerPause` unmarshals the payload into a `routing.PlayingState` struct and calls `gameState.HandlePause` to update the local game state, then acknowledges the message so it's removed from the queue.


## Army Moves (Topic Routing)

Each client subscribes to `army_moves.*` on the `peril_topic` exchange using
a transient queue named `army_moves.<username>`. This lets every connected
client receive moves broadcast by any player.

When a player runs the `move` command, the client publishes the resulting
`ArmyMove` to the `peril_topic` exchange with the routing key
`army_moves.<username>`, where `<username>` is the player who issued the move.

This demonstrates RabbitMQ's topic exchange wildcard matching:
- `*` matches exactly one word in the routing key
- `#` matches zero or more words


## Message Acknowledgements

Consumers now acknowledge messages after processing:

- `Ack`: message processed successfully
- `NackRequeue`: processing failed; retry the message
- `NackDiscard`: processing failed; discard the message

Move events are acknowledged only for safe moves or war outcomes. Invalid or self-originated moves are discarded.


## Dead Letter Queue (DLQ) Integration

To prevent unprocessable or discarded messages from being permanently lost, all client and server queues are configured with a Dead Letter Exchange (`x-dead-letter-exchange`).

- **Handling Rejected Messages**: When a consumer rejects a message without requeuing (`NackDiscard`), RabbitMQ forwards the message to the DLX rather than dropping it.
- **Inspection & Debugging**: Failed or unhandled messages (e.g., self-originating moves) land in `peril_dlq` where they can be inspected, analyzed, or replayed without disrupting active game queues.

## NackRequeue

Publish the return outcome of move, in HandlerMove which will be used to display in both of the client,

- first we did NackRequeue then after messing around and getting the endless loop of NackRequeue, we fix it with AckType (Ack) for case when fails
