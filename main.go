package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/humanbeeng/distributed-cache/proto"
)

func main() {

	var (
		leader = flag.Bool("leader", false, "true/false: Determines whether server is leader of not")
	)

	flag.Parse()

	if *leader {
		startLeader()
	} else {
		startFollower("localhost:3000")
	}

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

	conn, err := net.Dial("tcp", leaderAddr)
	if err != nil {
		log.Fatalf("Unable to connect to leader: %s", leaderAddr)
	}

	cmdJoin := proto.CmdJoin
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, cmdJoin)
	conn.Write(buf.Bytes())

	server := NewServer(serverOpts)
	server.Start()
}
