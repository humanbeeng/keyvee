package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/humanbeeng/distributed-cache/cache"
	"github.com/humanbeeng/distributed-cache/client"
	"github.com/humanbeeng/distributed-cache/proto"
)

type ServerOpts struct {
	ListenAddr string
	LeaderAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	follwersMu sync.RWMutex
	followers  map[*client.Client]struct{}
	cache      cache.Cacher
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      cache.New(),
		follwersMu: sync.RWMutex{},
		followers:  make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("Server starting on port %v", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	log.Println("Connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn)

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("unable to parse", err)
			break
		}
		s.handleCommand(conn, cmd)

	}
}

func (s *Server) handleCommand(conn net.Conn, cmd proto.Command) {
	switch cmd {
	case proto.CmdGet:
		{
			cmdGet := proto.ParseGetCommand(conn)
			s.handleGetCommand(conn, *cmdGet)
		}
	case proto.CmdSet:
		{
			cmdSet := proto.ParseSetCommand(conn)
			s.handleSetCommand(conn, *cmdSet)
		}
	case proto.CmdJoin:
		{
			fmt.Println("Follower join request came in", conn.RemoteAddr())
			s.follwersMu.Lock()
			s.followers[&client.Client{Conn: conn}] = struct{}{}
			s.follwersMu.Unlock()
		}
	case proto.CmdDel:
		{
			cmdDel := proto.ParseDelCommand(conn)
			s.handleDelCommand(conn, *cmdDel)
		}
	}

}

func (s *Server) handleSetCommand(conn net.Conn, cmdSet proto.CommandSet) {
	resp := proto.ResponseSet{}

	go func() {
		for follower := range s.followers {
			fmt.Printf("Sending SET %s : %s\n", cmdSet.Key, cmdSet.Value)
			follower.Set(string(cmdSet.Key), string(cmdSet.Value), cmdSet.TTL)
		}
	}()

	err := s.cache.Set(cmdSet.Key, cmdSet.Value, time.Second*time.Duration(cmdSet.TTL))
	if err != nil {
		resp.Status = proto.Error
		conn.Write(resp.Bytes())
		return
	}

	resp.Status = proto.Ok
	conn.Write(resp.Bytes())
}

func (s *Server) handleGetCommand(conn net.Conn, cmdGet proto.CommandGet) {
	resp := proto.ResponseGet{}

	val, err := s.cache.Get(cmdGet.Key)
	if err != nil {
		resp.Status = proto.Error
		conn.Write(resp.Bytes())
		return
	}

	resp.Value = val
	resp.Status = proto.Ok
	conn.Write(resp.Bytes())
}

func (s *Server) handleDelCommand(conn net.Conn, cmdDel proto.CommandDel) {
	resp := proto.ResponseDel{}

	go func() {
		for follower := range s.followers {
			fmt.Printf("Sending DEL %s\n", cmdDel.Key)
			follower.Del(string(cmdDel.Key))
		}
	}()

	err := s.cache.Delete(cmdDel.Key)
	if err != nil {
		resp.Status = proto.Error
		conn.Write(resp.Bytes())
		return
	}

	resp.Status = proto.Ok
	conn.Write(resp.Bytes())
}
