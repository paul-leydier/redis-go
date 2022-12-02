# ⚡️redis-go ⚡️

Initially inspired by the great work at [CodeCrafters](https://app.codecrafters.io/courses/redis?track=go), I decided
to try and implement Redis from scratch in golang, as a learning exercise.

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
