build:
	go build -o bin/distributed-cache

run: build
	./bin/distributed-cache -leader true

follower:build
	./bin/distributed-cache