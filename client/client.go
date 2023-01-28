package client

import (
	"fmt"
	"net"

	"github.com/humanbeeng/distributed-cache/proto"
)

type Client struct {
	Conn net.Conn
}

func New(serveraddr string) (*Client, error) {
	conn, err := net.Dial("tcp", serveraddr)
	if err != nil {
		return nil, fmt.Errorf("server conn err: %v", err.Error())
	}
	return &Client{Conn: conn}, nil
}

func (c *Client) Get(key string) (string, error) {
	cmdGet := &proto.CommandGet{
		Key: []byte(key),
	}

	_, err := c.Conn.Write(cmdGet.Bytes())
	if err != nil {
		return "", fmt.Errorf("error while reading from server %s", err.Error())
	}

	resp := proto.ParseGetResponse(c.Conn)
	if resp.Status == proto.Error {
		return "", fmt.Errorf("no key found for %s", key)
	}
	return string(resp.Value), nil

}

func (c *Client) Set(key string, value string, ttl int32) error {
	cmdset := &proto.CommandSet{Key: []byte(key), Value: []byte(value), TTL: ttl}
	_, err := c.Conn.Write(cmdset.Bytes())
	fmt.Printf("SET %v %v\n", key, value)

	if err != nil {
		return err
	}

	resp := proto.ParseSetResponse(c.Conn)

	if resp.Status == proto.Error {
		return fmt.Errorf("unable to set %s", key)
	}
	return nil
}

func (c *Client) Del(key string) error {
	cmdDel := &proto.CommandDel{Key: []byte(key)}
	_, err := c.Conn.Write(cmdDel.Bytes())
	if err != nil {
		return err
	}

	resp := proto.ParseDelResponse(c.Conn)

	if resp.Status == proto.Error {
		return fmt.Errorf("unable to del %s", key)
	}
	return nil
}
