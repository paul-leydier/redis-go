# ⚡️redis-go ⚡️

Initially inspired by the great work at [CodeCrafters](https://app.codecrafters.io/courses/redis?track=go), I decided
to try and implement Redis from scratch in golang, as a learning exercise.

## Features

This redis implementation includes the following features:
- TCP listening on exposed port `6379`
- Handle concurrent client connections
- Serve [RESP](https://redis.io/docs/reference/protocol-spec/) encoded queries coming from clients
- Return [RESP](https://redis.io/docs/reference/protocol-spec/) encoded responses, including
  error messages
- Serve the following Redis commands:
  - [PING](https://redis.io/commands/ping/)
  - [ECHO](https://redis.io/commands/echo/)
  - [GET](https://redis.io/commands/get/)
  - [SET](https://redis.io/commands/set/)

## Repository structure

The source code can be found in the `src` folder.

```
- src
    |- client
    |- core
    |- server
```

There, you'll find:
- The `client` folder (`redis-go/redis`) - an implementation of a Redis golang client
- The `server` folder (`redis-go/server`) - an implementation of a Redis server
- The `core` folder (`redis-go/core`) - code which can be used by both the client and server packages.
