package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/humanbeeng/distributed-cache/client"
)

func main() {

	var (
		leader = flag.Bool("leader", false, "true/false: Determines whether server is leader of not")
	)

	flag.Parse()

	go func() {
		time.Sleep(time.Second * 1)
		SendCommand()
	}()

	if *leader {
		startLeader()
	} else {
		startFollower("localhost:3000")
	}

}

func SendCommand() {

	client, err := client.New("localhost:3000")
	if err != nil {
		log.Fatal("Unable to establish connection")
		return
	}

	client.Set("heyy", "there", 1)

	val, _ := client.Get("heyy")
	fmt.Println(val)

}

func startLeader() {
	fmt.Println("Leader UP at :3000")
	serverOpts := ServerOpts{
		ListenAddr: "localhost:3000",
		LeaderAddr: "",
		IsLeader:   true,
	}

	server := NewServer(serverOpts)
	server.Start()
}

func startFollower(leaderAddr string) {

	serverOpts := ServerOpts{
		ListenAddr: "",
		LeaderAddr: leaderAddr,
		IsLeader:   false,
	}
	fmt.Println("Follower UP")

	server := NewServer(serverOpts)
	server.Start()
}
