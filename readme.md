# notesbot

The telegram bot for reading and synchronizing notes in markdown format. The bot consists of the client and server parts. Client part parse notes and send them to the server-side. Server-side handle them as telegram bot. 

Bot commands:
`get list` - shows index of all notes.

Environment variable for server app.
```
s_PORT
TOKEN - telegram bot token
DB_FILE
```