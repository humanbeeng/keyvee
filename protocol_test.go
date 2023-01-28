package main

import (
	"bytes"
	"testing"

	"github.com/humanbeeng/distributed-cache/proto"
	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &proto.CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())
	pCmd, err := proto.ParseCommand(r)
	if err != nil {
		assert.Fail(t, "Failed")
	}
	assert.NotNil(t, pCmd, proto.CommandSet{})
	assert.Equal(t, cmd, pCmd)

}

func TestParseGetCommand(t *testing.T) {
	cmd := &proto.CommandGet{Key: []byte("Foo")}
	r := bytes.NewReader(cmd.Bytes())
	cmdGet, err := proto.ParseCommand(r)
	if err != nil {
		assert.Fail(t, "Failed")
	}
	assert.Equal(t, cmd, cmdGet)

}
