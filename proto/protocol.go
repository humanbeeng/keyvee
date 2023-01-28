package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte
type Status byte

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

const (
	Ok Status = iota
	Error
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}
type CommandGet struct {
	Key []byte
}

type CommandJoin struct{}

type ResponseGet struct {
	Status Status
	Value  []byte
}

type ResponseSet struct {
	Status Status
}

func (c *CommandSet) Bytes() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)

	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	valueLen := int32(len(c.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int32(c.TTL))
	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)

	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)

	valLen := int32(len(r.Value))
	binary.Write(buf, binary.LittleEndian, valLen)
	binary.Write(buf, binary.LittleEndian, r.Value)
	return buf.Bytes()
}

func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (Command, error) {
	var cmd Command

	binary.Read(r, binary.LittleEndian, &cmd)
	fmt.Println("command from r: ", cmd)
	return cmd, nil
}

func ParseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = ttl

	return cmd
}

func ParseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}
	var keyLen int32

	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)

	binary.Read(r, binary.LittleEndian, &cmd.Key)
	return cmd
}

func ParseSetResponse(r io.Reader) *ResponseSet {
	resp := &ResponseSet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp
}

func ParseGetResponse(r io.Reader) *ResponseGet {
	resp := &ResponseGet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)
	var valLen int32
	binary.Read(r, binary.LittleEndian, &valLen)
	resp.Value = make([]byte, valLen)
	binary.Read(r, binary.LittleEndian, &resp.Value)
	return resp
}
