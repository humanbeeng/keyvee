# Distributed Cache

A simple KV-Store which runs as a Leader Follower cluster. 

## Features:
- Handwritten binary protocol for client-server communication (No JSON or Protobuf 🙅‍♂️).
- Auto registration of new followers added to cluster.
- TTL for KV entries.

## Workings:

- Leader follower architecture.
- Leader maintains a map of its followers with their addresses.
- Leader runs on `:3000`
- Follower runs on random ports.
- When follower boots up, it establishes a new TCP connection with Leader and Leader registers the follower by storing its details.
- When a SET, DEL command is given, leader ensures that all the followers are in sync with the given command.
- SET command can be paired up with a TTL(in seconds), after which entry will be removed from the store

```sh
$ make leader
$ make follower
```

## Enhancements

- [ ] Have a response timeout from the followers so that leader can know when a follower has went down.
- [ ] GET queries should be forwarded to the followers by the leader. (Optional, bring in a reverse proxy(handwritten) which can distribute the load in a Round-Robin fashion)
- [ ] Implement a consensus algorithm which can elect a leader amongst themselves, when a leader fails to send a heartbeat to its followers.

## Working sample through [CLI](https://github.com/humanbeeng/kv-cli)

![Screenshot from 2023-01-30 10-21-06](https://user-images.githubusercontent.com/37271977/215390645-a0a9debd-3378-4860-827d-2e960ed25f67.png)
