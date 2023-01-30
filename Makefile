build:
	go build -o bin/distributed-cache

leader: build
	./bin/distributed-cache -leader true

follower:build
	./bin/distributed-cache