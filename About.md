# Distributed Cache

Workings:

- Leader follower architecture.
- Leader maintains a map of its followers with their addresses.
- Leader runs on :3000.
- Follower runs on random ports.
- When follower boots up, it establishes a new TCP connection with Leader and Leader registers the follower by storing its details.
- When a SET, DELETE command is given, leader ensures that all the followers are in sync with the given command.
- (Optional) Have a response timeout from the followers so that leader can know when a follower has went down.
- (Optional)GET queries should be forwarded to the followers by the leader. (Optional, bring in a reverse proxy(handwritten) which can distribute the load in a Round-Robin fashion)
