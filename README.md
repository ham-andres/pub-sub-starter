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
